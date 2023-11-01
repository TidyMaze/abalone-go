package main

import (
	. "abalone-go/engine"
	"abalone-go/helpers"
	"fmt"
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

	game.Push(Coord3D{X: -2, Y: -1, Z: 3}, Right)

	println("Game pushed to the right:")
	println(game.Show())
}
