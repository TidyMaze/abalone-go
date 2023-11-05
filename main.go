package main

import (
	. "abalone-go/engine"
	"fmt"
	"log"
	"math/rand"
)

func main() {
	// Create a new game
	currentGame := NewGame()
	fillRandomly(currentGame)

	log.Println("Random game:")
	log.Println(currentGame.Show())

	for !currentGame.IsOver() {
		log.Println(fmt.Sprintf("=========== Turn %d ===========", currentGame.Turn))

		validMoves := currentGame.GetValidMoves()

		if len(validMoves) == 0 {
			log.Println("No more valid moves")
			break
		}

		log.Println(fmt.Sprintf("Valid moves size: %d", len(validMoves)))

		firstValidMove := randIn(validMoves)

		switch t := firstValidMove.(type) {
		case PushLine:
			pushLine := firstValidMove.(PushLine)

			log.Println(fmt.Sprintf("Pushing line: %v", t))
			err := currentGame.Push(pushLine.From, pushLine.Direction)
			if err != nil {
				panic(err)
			}

			log.Println(fmt.Sprintf("New game state: %s", currentGame.Show()))
		default:
			panic("Invalid move type" + fmt.Sprintf("%T", t))
		}
	}

	log.Println("Game over")

	if currentGame.Winner == 0 {
		log.Println("Draw")
	} else {
		log.Println(fmt.Sprintf("Winner: %d", currentGame.Winner))
	}
}

func randIn[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}

func fillRandomly(game *Game) {
	for y := -4; y <= 4; y++ {
		for x := -4; x <= 4; x++ {
			coord := Coord2D{X: x, Y: y}.To3D()

			if IsValidCoord(coord) {
				game.SetGrid(coord, rand.Intn(3))
			}
		}
	}
}
