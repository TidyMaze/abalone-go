package main

import (
	. "abalone-go/engine"
	"abalone-go/helpers"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")

	// Create a new game
	game := NewGame()
	println(game.Show())

	helpers.AssertEqual(Coord2D{X: 3, Y: -3}.To3D(), Coord3D{X: 3, Z: -3})
	helpers.AssertEqual(Coord2D{X: -3, Y: 3}.To3D(), Coord3D{X: -3, Z: 3})

	game.SetGrid(Coord2D{X: 3, Y: -3}.To3D(), 1)
	game.SetGrid(Coord2D{X: -3, Y: 3}.To3D(), 2)

	println(game.Show())
}
