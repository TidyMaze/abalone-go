package main

import (
	"abalone-go/engine"
	"context"
	"flag"
	"fmt"
	"github.com/yaricom/goNEAT/v4/experiment"
	"github.com/yaricom/goNEAT/v4/neat"
	"github.com/yaricom/goNEAT/v4/neat/genetics"
	"github.com/yaricom/goNEAT/v4/neat/network"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	func main() {
//		// Create a new game
//		currentGame := NewGame()
//		fillRandomly(currentGame)
//
//		log.Println("Random game:")
//		log.Println(currentGame.Show())
//
//		for !currentGame.IsOver() {
//			log.Println(fmt.Sprintf("=========== Turn %d ===========", currentGame.Turn))
//
//			validMoves := currentGame.GetValidMoves()
//
//			if len(validMoves) == 0 {
//				log.Println("No more valid moves")
//				break
//			}
//
//			log.Println(fmt.Sprintf("Valid moves size: %d", len(validMoves)))
//
//			firstValidMove := helpers.RandIn(validMoves)
//
//			switch t := firstValidMove.(type) {
//			case PushLine:
//				pushLine := firstValidMove.(PushLine)
//
//				log.Println(fmt.Sprintf("Pushing line: %v", t))
//				err := currentGame.Push(pushLine.From, pushLine.Direction)
//				if err != nil {
//					panic(err)
//				}
//
//				log.Println(fmt.Sprintf("New game state: %s", currentGame.Show()))
//			default:
//				panic("Invalid move type" + fmt.Sprintf("%T", t))
//			}
//		}
//
//		log.Println("Game over")
//
//		if currentGame.Winner == 0 {
//			log.Println("Draw")
//		} else {
//			log.Println(fmt.Sprintf("Winner: %d", currentGame.Winner))
//		}
//	}
//
//	func fillRandomly(game *Game) {
//		for y := -4; y <= 4; y++ {
//			for x := -4; x <= 4; x++ {
//				coord := Coord2D{X: x, Y: y}.To3D()
//
//				if IsValidCoord(coord) {
//					game.SetGrid(coord, rand.Intn(3))
//				}
//			}
//		}
//	}
func main() {
	var outDirPath = flag.String("out", "./out", "The output directory to store results.")
	var contextPath = flag.String("context", "./data/abalone.neat", "The execution context configuration file.")
	var trialsCount = flag.Int("trials", 0, "The number of trials for experiment. Overrides the one set in configuration.")
	var logLevel = flag.String("log_level", "", "The logger level to be used. Overrides the one set in configuration.")
	var randSeed = flag.Int64("seed", 0, "The seed for random number generator")

	flag.Parse()

	// Seed the random-number generator with current time so that
	// the numbers will be different every time we run.
	seed := time.Now().Unix()
	if randSeed != nil {
		seed = *randSeed
	}
	rand.Seed(seed)

	// Load NEAT options
	neatOptions, err := neat.ReadNeatOptionsFromFile(*contextPath)
	if err != nil {
		log.Fatal("Failed to load NEAT options: ", err)
	}

	// Load Genome
	traits := make([]*neat.Trait, 0)
	traits = append(traits, neat.NewTrait())

	nodes := make([]*network.NNode, 0)
	// Add the 61*2 input nodes
	for i := 0; i < 61*2; i++ {
		nodes = append(nodes, network.NewNNode(i+1, network.InputNeuron))
	}

	// Add the 61 + 6 output nodes
	for i := 0; i < 61+6; i++ {
		nodes = append(nodes, network.NewNNode(i+1+61*2, network.OutputNeuron))
	}

	genes := make([]*genetics.Gene, 0)
	// Link each input to each output
	for i := 0; i < 61*2; i++ {
		for j := 0; j < 61+6; j++ {
			genes = append(genes, genetics.NewGene(0.0, nodes[i], nodes[j+61*2], true, 1.0, 0))
		}
	}

	startGenome := genetics.NewGenome(1, traits, nodes, genes)
	//fmt.Println(startGenome)

	// Check if output dir exists
	outDir := *outDirPath
	if _, err := os.Stat(outDir); err == nil {
		// backup it
		backUpDir := fmt.Sprintf("%s-%s", outDir, time.Now().Format("2006-01-02T15_04_05"))
		// clear it
		err = os.Rename(outDir, backUpDir)
		if err != nil {
			log.Fatal("Failed to do previous results backup: ", err)
		}
	}
	// create output dir
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create output directory: ", err)
	}

	// Override neatOptions configuration parameters with ones set from command line
	if *trialsCount > 0 {
		neatOptions.NumRuns = *trialsCount
	}
	if len(*logLevel) > 0 {
		neat.LogLevel = neat.LoggerLevel(*logLevel)
	}

	// create experiment
	expt := experiment.Experiment{
		Id:       0,
		Trials:   make(experiment.Trials, neatOptions.NumRuns),
		RandSeed: seed,
	}
	var generationEvaluator experiment.GenerationEvaluator
	expt.MaxFitnessScore = 100.0 // as given by fitness function definition
	generationEvaluator = engine.NewAbaloneGenerationEvaluator(outDir)

	// prepare to execute
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	// run experiment in the separate GO routine
	go func() {
		if err = expt.Execute(neat.NewContext(ctx, neatOptions), startGenome, generationEvaluator, nil); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	// register handler to wait for termination signals
	//
	go func(cancel context.CancelFunc) {
		fmt.Println("\nPress Ctrl+C to stop")

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case <-signals:
			// signal to stop test fixture
			cancel()
		case err = <-errChan:
			// stop waiting
		}
	}(cancel)

	// Wait for experiment completion
	//
	err = <-errChan
	if err != nil {
		// error during execution
		log.Fatalf("Experiment execution failed: %s", err)
	}

	// Print experiment results statistics
	//
	expt.PrintStatistics()

	fmt.Printf(">>> Configuration file: %s\n", *contextPath)

	// Save experiment data in native format
	//
	expResPath := fmt.Sprintf("%s/%s.dat", outDir, "abalone")
	if expResFile, err := os.Create(expResPath); err != nil {
		log.Fatal("Failed to create file for experiment results", err)
	} else if err = expt.Write(expResFile); err != nil {
		log.Fatal("Failed to save experiment results", err)
	}

	// Save experiment data in Numpy NPZ format if requested
	//
	npzResPath := fmt.Sprintf("%s/%s.npz", outDir, "abalone")
	if npzResFile, err := os.Create(npzResPath); err != nil {
		log.Fatalf("Failed to create file for experiment results: [%s], reason: %s", npzResPath, err)
	} else if err = expt.WriteNPZ(npzResFile); err != nil {
		log.Fatal("Failed to save experiment results as NPZ file", err)
	}
}
