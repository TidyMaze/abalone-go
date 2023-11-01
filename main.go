package main

import (
	"abalone-go/engine"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")

	// Create a new game
	game := engine.NewGame()
	println(game.Show())
}
