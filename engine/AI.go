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
	"math/rand"
	"os"
)

const CountGames = 10

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
		//log.Println(fmt.Sprintf("[Gen %d] Evaluating organism: %d", epoch.Id, org.Genotype.Id))

		_, err := e.orgEvaluate(org, epoch)

		//log.Println(fmt.Sprintf("[Gen %d] Evaluated organism: %d, fitness: %f", epoch.Id, org.Genotype.Id, org.Fitness))

		if err != nil {
			panic(err)
		}

		if epoch.Champion == nil || org.Fitness > epoch.Champion.Fitness {
			epoch.WinnerNodes = len(org.Genotype.Nodes)
			epoch.WinnerGenes = org.Genotype.Extrons()
			epoch.WinnerEvals = options.PopSize*epoch.Id + org.Genotype.Id
			epoch.Champion = org
		}

		totalFitness = totalFitness + org.Fitness
	}

	log.Println(fmt.Sprintf("[Gen %d] Found new champion with fitness: %f", epoch.Id, epoch.Champion.Fitness))

	if optPath, err := utils.WriteGenomePlain("abalone_champion", e.OutputPath, epoch.Champion, epoch); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to dump champion genome, reason: %s\n", err))
	} else {
		neat.InfoLog(fmt.Sprintf("Dumped champion genome to: %s\n", optPath))
	}

	averageFitness := totalFitness / float64(len(pop.Organisms))

	log.Println(fmt.Sprintf("[Gen %d] Average fitness: %f for total fitness: %f and population size: %d",
		epoch.Id, averageFitness, totalFitness, len(pop.Organisms)))

	for _, specy := range pop.Species {
		max, avg := specy.ComputeMaxAndAvgFitness()
		log.Println(fmt.Sprintf("[Gen %d] Species id %d, organisms: %d, average fitness: %f, max fitness: %f",
			epoch.Id, specy.Id, len(specy.Organisms), avg, max))
	}

	for _, org := range pop.Organisms {
		println(fmt.Sprintf("Org %d, specy %d, fitness: %f, error: %f",
			org.Genotype.Id, org.Species.Id, org.Fitness, org.Error))
	}

	// Fill statistics about current epoch
	epoch.FillPopulationStatistics(pop)

	bestFitnessBySpecy := epoch.Fitness
	log.Println(fmt.Sprintf("[Gen %d] Epoch statistics: %f, fitness: %v", epoch.Id, bestFitnessBySpecy.Mean(), bestFitnessBySpecy))

	pop.MeanFitness = averageFitness

	helpers.AssertEqual(false, pop.MeanFitness == 0.0)
	helpers.AssertEqual(false, bestFitnessBySpecy.Mean() == 0.0)

	// Only print to file every print_every generation
	if epoch.Id%options.PrintEvery == 0 {
		if _, err := utils.WritePopulationPlain(e.OutputPath, pop, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump population, reason: %s\n", err))
			return err
		}
	}

	org := epoch.Champion
	//utils.PrintActivationDepth(org, true)

	genomeFile := "abalone_champion_genome"
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

	// write epoch, average fitness and champion fitness to CSV file
	if err := writeGenerationCSV(e.OutputPath, epoch, averageFitness); err != nil {
		return err
	}

	return nil
}

func writeGenerationCSV(outputPath string, epoch *experiment.Generation, averageFitness float64) error {
	epochId := epoch.Id
	championFitness := epoch.Champion.Fitness

	// append to file
	filePath := outputPath + "/generation.csv"

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	helpers.AssertEqual(false, averageFitness == 0.0)

	log.Println(fmt.Sprintf("[Gen %d] Writing generation stats (average fitness: %f, champion fitness: %f) to CSV file %s",
		epochId, averageFitness, championFitness, filePath))

	if _, err := f.WriteString(fmt.Sprintf("%d,%f,%f\n", epochId, averageFitness, championFitness)); err != nil {
		return err
	}

	log.Println(fmt.Sprintf("[Gen %d] Wrote generation stats to CSV file %s", epochId, filePath))

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

	for gameId := 0; gameId < CountGames; gameId++ {
		//log.Println(fmt.Sprintf("[Gen %d][Org %d] Starting game %d", epoch.Id, organism.Genotype.Id, gameId))
		game := NewGame(startingGrid)

		for !game.IsOver() && game.Turn < 127 {
			//log.Println(fmt.Sprintf("[Gen %d][Org %d] Game %d, turn %d", epoch.Id, organism.Genotype.Id, gameId, game.Turn))

			var move Move
			if game.currentPlayer == 1 {
				// player 1 is the organism

				movePtr, err := e.predictSingleMove(phenotype, netDepth, *game)

				if err != nil {
					return false, err
				}

				move = *movePtr
				//log.Println(fmt.Sprintf("[Gen %d][Org %d] Predicted move: %v", epoch.Id, organism.Genotype.Id, move))

				switch move.(type) {
				case PushLine:
					err := game.Move(move)

					if err != nil {
						// invalid move, opponent wins
						log.Println(fmt.Sprintf("[Gen %d][Org %d] Invalid move: %v", epoch.Id, organism.Genotype.Id, move))
						game.Winner = 2
						game.score[2] = 6
						panic(fmt.Sprintf("Invalid move: %v", move))
					}
				default:
					panic(fmt.Sprintf("Invalid move type: %T", move))
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

		//log.Println(fmt.Sprintf("[Gen %d][Org %d] Finished game %d, score: %v after %d turns",
		//	epoch.Id, organism.Genotype.Id, gameId, game.score, game.Turn))

		totalCaptured += int(game.score[1])
		totalEnemyCaptured += int(game.score[2])
	}

	avgScoreDiff := (float64(totalCaptured) - float64(totalEnemyCaptured)) / float64(CountGames)
	score := avgScoreDiff
	ideal := float64(6)  // win every game at 6-0 for player 1
	worst := float64(-6) // lose every game at 0-6 for player 1

	// normalized between 0 and 1
	normalized := (score - worst) / (ideal - worst)

	log.Println(fmt.Sprintf("[Gen %d][Org %d] Finished ranking organism, score diff: %f, normalized: %f, ideal: %f",
		epoch.Id, organism.Genotype.Id, avgScoreDiff, normalized, 1.0))

	organism.Fitness = normalized
	organism.Error = math.Abs(1.0 - normalized)

	return false, nil
}

func (e *AbaloneGenerationEvaluator) predictSingleMove(phenotype *network.Network, netDepth int, game Game) (*Move, error) {
	validMoves := game.GetValidMoves()

	if len(validMoves) == 0 {
		return nil, fmt.Errorf("no valid moves")
	}

	bestMoveScore := -1000000.0
	var bestMove *Move

	rand.Shuffle(len(validMoves), func(i, j int) {
		validMoves[i], validMoves[j] = validMoves[j], validMoves[i]
	})

	for _, move := range validMoves {
		nextState := game.Copy()
		err := nextState.Move(move)

		if err != nil {
			return nil, err
		}

		var in []float64

		// Set the input values
		for y := int8(-4); y <= 4; y++ {
			for x := int8(-4); x <= 4; x++ {
				coord := Coord2D{x, y}.To3D()
				if IsValidCoord(coord) {
					cellOwner := nextState.GetGrid(coord)
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
		score := phenotype.Outputs[0].Activation

		//log.Println(fmt.Sprintf("Move: %v, score: %f", move, score))

		if score > bestMoveScore || bestMove == nil {
			bestMoveScore = score
			bestMove = &move
			//log.Println(fmt.Sprintf("New best move: %v, score: %f", move, score))
		}

		// Flush network for subsequent use
		if _, err = phenotype.Flush(); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to flush network: %s", err))
			return nil, err
		}
	}

	//log.Println(fmt.Sprintf("Best move: %v, score: %f among %d valid moves", *bestMove, bestMoveScore, len(validMoves)))

	return bestMove, nil
}
