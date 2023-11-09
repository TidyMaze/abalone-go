package engine

import (
	"abalone-go/helpers"
	"context"
	"fmt"
	"github.com/yaricom/goNEAT/v4/experiment"
	"github.com/yaricom/goNEAT/v4/experiment/utils"
	"github.com/yaricom/goNEAT/v4/neat"
	"github.com/yaricom/goNEAT/v4/neat/genetics"
	"github.com/yaricom/goNEAT/v4/neat/network"
	"log"
	"math"
)

const CountGames = 50

type AbaloneGenerationEvaluator struct {
	OutputPath string
}

func (e *AbaloneGenerationEvaluator) GenerationEvaluate(ctx context.Context, pop *genetics.Population, epoch *experiment.Generation) error {
	options, ok := neat.FromContext(ctx)
	if !ok {
		return neat.ErrNEATOptionsNotFound
	}

	totalFitness := 0.0

	for _, org := range pop.Organisms {
		res, err := e.orgEvaluate(org, epoch)
		if err != nil {
			return err
		}

		if res && (epoch.Champion == nil || org.Fitness > epoch.Champion.Fitness) {
			epoch.Solved = true
			epoch.WinnerNodes = len(org.Genotype.Nodes)
			epoch.WinnerGenes = org.Genotype.Extrons()
			epoch.WinnerEvals = options.PopSize*epoch.Id + org.Genotype.Id
			epoch.Champion = org
			if epoch.WinnerNodes == 5 {
				// You could dump out optimal genomes here if desired
				if optPath, err := utils.WriteGenomePlain("xor_optimal", e.OutputPath, org, epoch); err != nil {
					neat.ErrorLog(fmt.Sprintf("Failed to dump optimal genome, reason: %s\n", err))
				} else {
					neat.InfoLog(fmt.Sprintf("Dumped optimal genome to: %s\n", optPath))
				}
			}
		}

		totalFitness = totalFitness + org.Fitness
	}

	log.Println(fmt.Sprintf("[Gen %d] Average fitness: %f for total fitness: %f and population size: %d",
		epoch.Id, totalFitness/float64(len(pop.Organisms)), totalFitness, len(pop.Organisms)))

	// Fill statistics about current epoch
	epoch.FillPopulationStatistics(pop)

	// Only print to file every print_every generation
	if epoch.Solved || epoch.Id%options.PrintEvery == 0 {
		if _, err := utils.WritePopulationPlain(e.OutputPath, pop, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump population, reason: %s\n", err))
			return err
		}
	}

	if epoch.Solved {
		// print winner organism
		org := epoch.Champion
		utils.PrintActivationDepth(org, true)

		genomeFile := "xor_winner_genome"
		// Prints the winner organism's Genome to the file!
		if orgPath, err := utils.WriteGenomePlain(genomeFile, e.OutputPath, org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's genome, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's genome dumped to: %s\n", epoch.Id, orgPath))
		}

		// Prints the winner organism's phenotype to the DOT file!
		if orgPath, err := utils.WriteGenomeDOT(genomeFile, e.OutputPath, org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's phenome DOT graph, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's phenome DOT graph dumped to: %s\n",
				epoch.Id, orgPath))
		}

		// Prints the winner organism's phenotype to the Cytoscape JSON file!
		if orgPath, err := utils.WriteGenomeCytoscapeJSON(genomeFile, e.OutputPath, org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's phenome Cytoscape JSON graph, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's phenome Cytoscape JSON graph dumped to: %s\n",
				epoch.Id, orgPath))
		}
	}

	return nil
}

func NewAbaloneGenerationEvaluator(outputPath string) experiment.GenerationEvaluator {
	return &AbaloneGenerationEvaluator{OutputPath: outputPath}
}

// orgEvaluate evaluates fitness of the provided organism
func (e *AbaloneGenerationEvaluator) orgEvaluate(organism *genetics.Organism, epoch *experiment.Generation) (bool, error) {
	// evaluate the organism by running 100 games against random opponent
	// fitness is the win rate of the organism

	// INPUT: 61 cells, 2 possible states (1,2) = 122 input nodes
	// OUTPUT: 61 nodes for the push origin, 6 nodes for the push direction = 67 output nodes

	phenotype, err := organism.Phenotype()
	if err != nil {
		return false, err
	}

	netDepth, err := phenotype.MaxActivationDepthWithCap(0) // The max depth of the network to be activated
	if err != nil {
		neat.WarnLog(fmt.Sprintf(
			"Failed to estimate maximal depth of the network with loop:\n%s\nUsing default depth: %d",
			organism.Genotype, netDepth))
	}
	neat.DebugLog(fmt.Sprintf("Network depth: %d for organism: %d\n", netDepth, organism.Genotype.Id))
	if netDepth == 0 {
		neat.DebugLog(fmt.Sprintf("ALERT: Network depth is ZERO for Genome: %s", organism.Genotype))
		return false, nil
	}

	totalCaptured := 0
	totalEnemyCaptured := 0
	invalidMovesCount := 0

	for gameId := 0; gameId < CountGames; gameId++ {
		game := NewGame()
		game.grid = buildStartingGrid()

		for !game.IsOver() && game.Turn < 30 {
			var move Move
			if game.currentPlayer == 1 {
				// player 1 is the organism

				move, err = e.predictSingleMove(phenotype, netDepth, *game)
				if err != nil {
					return false, err
				}

				switch move.(type) {
				case PushLine:
					_, _, pushError := game.checkCanPush(move.(PushLine).From, move.(PushLine).Direction)

					if pushError != nil {
						// invalid move, opponent wins
						game.Winner = 2
						game.score[2] = 6
						invalidMovesCount = invalidMovesCount + 1
					} else {
						//log.Println(fmt.Sprintf("Predicted move: %v", move))

						err := game.Move(move)
						if err != nil {
							return false, err
						}
					}
				}
			} else {
				// player 2 is the random opponent

				// pick a random move
				possibleMoves := game.GetValidMoves()
				move = helpers.RandIn(possibleMoves)

				err := game.Move(move)
				if err != nil {
					return false, err
				}
			}
		}

		totalCaptured += game.score[1]
		totalEnemyCaptured += game.score[2]
	}

	scoreDiff := float64(totalCaptured) - float64(totalEnemyCaptured)
	score := scoreDiff*100 - float64(invalidMovesCount)
	ideal := float64(6 * CountGames * 100) // win every game at 6-0 for player 1

	log.Println(fmt.Sprintf("[Gen %d][Org %d] Finished ranking organism, score diff: %f, invalid moves: %d, score: %f, ideal: %f",
		epoch.Id, organism.Genotype.Id, scoreDiff, invalidMovesCount, score, ideal))

	organism.Fitness = score
	organism.Error = math.Abs(ideal - score)

	organism.IsWinner = false
	return organism.IsWinner, nil
}

func (e *AbaloneGenerationEvaluator) predictSingleMove(phenotype *network.Network, netDepth int, game Game) (Move, error) {
	err := error(nil)

	var in []float64

	// Set the input values
	for y := -4; y <= 4; y++ {
		for x := -4; x <= 4; x++ {
			coord := Coord2D{x, y}.To3D()
			if IsValidCoord(coord) {
				cellOwner := game.GetGrid(coord)
				player1 := 0.0
				player2 := 0.0

				if cellOwner == 1 {
					player1 = 1.0
				} else if cellOwner == 2 {
					player2 = 1.0
				}

				in = append(in, player1, player2)
			}
		}
	}

	if err = phenotype.LoadSensors(in); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to load sensors: %s", err))
		return nil, err
	}

	// Use depth to ensure full relaxation
	if success, err := phenotype.ForwardSteps(netDepth); err != nil || !success {
		neat.ErrorLog(fmt.Sprintf("Failed to activate network: %s", err))
		return nil, err
	}

	// Read output
	// Cell is the index of max output in the first 61 nodes
	// Direction is the index of max output in the last 6 nodes
	cell := -1
	for i := 0; i < 61; i++ {
		if cell == -1 || (phenotype.Outputs[i].Activation > phenotype.Outputs[cell].Activation) {
			cell = i
		}
	}

	direction := -1
	for i := 61; i < 67; i++ {
		if direction == -1 || (phenotype.Outputs[i].Activation > phenotype.Outputs[direction].Activation) {
			direction = i
		}
	}

	// Flush network for subsequent use
	if _, err = phenotype.Flush(); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to flush network: %s", err))
		return nil, err
	}

	indexToCoord3D := [61]Coord3D{
		{0, 4, -4},
		{1, 3, -4},
		{2, 2, -4},
		{3, 1, -4},
		{4, 0, -4},

		{-1, 4, -3},
		{0, 3, -3},
		{1, 2, -3},
		{2, 1, -3},
		{3, 0, -3},
		{4, -1, -3},

		{-2, 4, -2},
		{-1, 3, -2},
		{0, 2, -2},
		{1, 1, -2},
		{2, 0, -2},
		{3, -1, -2},
		{4, -2, -2},

		{-3, 4, -1},
		{-2, 3, -1},
		{-1, 2, -1},
		{0, 1, -1},
		{1, 0, -1},
		{2, -1, -1},
		{3, -2, -1},
		{4, -3, -1},

		{-4, 4, 0},
		{-3, 3, 0},
		{-2, 2, 0},
		{-1, 1, 0},
		{0, 0, 0},
		{1, -1, 0},
		{2, -2, 0},
		{3, -3, 0},
		{4, -4, 0},

		{-4, 3, 1},
		{-3, 2, 1},
		{-2, 1, 1},
		{-1, 0, 1},
		{0, -1, 1},
		{1, -2, 1},
		{2, -3, 1},
		{3, -4, 1},

		{-4, 2, 2},
		{-3, 1, 2},
		{-2, 0, 2},
		{-1, -1, 2},
		{0, -2, 2},
		{1, -3, 2},
		{2, -4, 2},

		{-4, 1, 3},
		{-3, 0, 3},
		{-2, -1, 3},
		{-1, -2, 3},
		{0, -3, 3},
		{1, -4, 3},

		{-4, 0, 4},
		{-3, -1, 4},
		{-2, -2, 4},
		{-1, -3, 4},
		{0, -4, 4},
	}

	fromCoord := indexToCoord3D[cell]

	return PushLine{fromCoord, Direction(direction - 61)}, nil
}
