package main

import (
	. "abalone-go/engine"
	"abalone-go/helpers"
	"fmt"
	"math/rand"
)

func main() {
	helpers.AssertEqual(Coord2D{X: 3, Y: -3}.To3D(), Coord3D{X: 3, Z: -3})
	helpers.AssertEqual(Coord2D{X: -3, Y: 3}.To3D(), Coord3D{X: -3, Z: 3})

	fmt.Println("Hello, World!")

	// Create a new game
	game := NewGame()

	println("Empty game:")
	println(game.Show())

	game.SetGrid(Coord2D{X: 3, Y: -3}.To3D(), 1)
	game.SetGrid(Coord2D{X: -3, Y: 3}.To3D(), 2)

	println("Game with 2 marbles:")
	println(game.Show())

	game.Push(Coord3D{X: -3, Y: 0, Z: 3}, Right)

	println("Game pushed to the right:")
	println(game.Show())

	currentMarbleCoord := &Coord3D{X: -2, Y: -1, Z: 3}

	for {
		if currentMarbleCoord == nil {
			println("No more marble to push")
			break
		}

		err := game.Push(*currentMarbleCoord, Right)
		if err != nil {
			panic(err)
		}

		destination := currentMarbleCoord.Add(Right)

		if !IsValidCoord(destination) {
			currentMarbleCoord = nil
		} else {
			currentMarbleCoord = &destination
		}

		println("Game pushed to the right:")
		println(game.Show())
	}

	randomGame := NewGame()
	fillRandomly(randomGame)

	println("Random game:")
	println(randomGame.Show())
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
