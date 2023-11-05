package main

import (
	. "abalone-go/engine"
	"abalone-go/helpers"
	"fmt"
	"log"
	"math/rand"
)

func main() {
	helpers.AssertEqual(Coord2D{X: 3, Y: -3}.To3D(), Coord3D{X: 3, Z: -3})
	helpers.AssertEqual(Coord2D{X: -3, Y: 3}.To3D(), Coord3D{X: -3, Z: 3})

	fmt.Println("Hello, World!")

	// Create a new game
	currentGame := NewGame()
	fillRandomly(currentGame)

	println("Random game:")
	println(currentGame.Show())

	for {
		validMoves := currentGame.GetValidMoves()

		if len(validMoves) == 0 {
			println("No more valid moves")
			break
		}

		log.Println(fmt.Sprintf("Valid moves size: %d", len(validMoves)))

		firstValidMove := validMoves[0]

		switch t := firstValidMove.(type) {
		case PushLine:
			pushLine := firstValidMove.(PushLine)

			log.Println(fmt.Sprintf("Pushing line: %v", t))
			err := currentGame.Push(pushLine.From, pushLine.Direction)
			if err != nil {
				panic(err)
			}

			println("New game state:")
			println(currentGame.Show())
		default:
			panic("Invalid move type" + fmt.Sprintf("%T", t))
		}
	}
}

func fillRandomly(game *Game) {
	for y := -3; y <= 3; y++ {
		for x := -3; x <= 3; x++ {
			coord := Coord2D{X: x, Y: y}.To3D()

			if IsValidCoord(coord) {
				game.SetGrid(coord, rand.Intn(3))
			}
		}
	}
}
