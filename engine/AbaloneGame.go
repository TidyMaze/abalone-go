package engine

import (
	"errors"
	"fmt"
)

/**
 * Game is the main game class.
 * It contains the grid and the players.
 *
 * The grid is a 2D array of integers.
 * Here is the meaning of the integers:
 * -1: non-existent cell (out of the hexagon)
 * 0: empty
 * 1: player 1
 * 2: player 2
 *
 * The players are represented by their color.
 * 1: white
 * 2: black
 *
 * Grid layout:
 * - center is (0, 0)
 * - x axis is 60 degrees to the top right
 * - y axis is 60 degrees to the top left
 * - z axis is vertical to the bottom
 *
 * The hexagon is 4 cells wide in the outer ring (side size) and 7 cells wide in the inner ring (diameter).
 * There are 6 directions expressed in (x,y,z) coordinates:
 * - top right (1, 0, -1)
 * - right (1, -1, 0)
 * - bottom right (0, -1, 1)
 * - bottom left (-1, 0, 1)
 * - left (-1, 1, 0)
 * - top left (0, 1, -1)
 */

type Game struct {
	grid map[Coord3D]int
}

func (g Game) show() {
}

func NewGame() *Game {
	game := &Game{}
	game.grid = buildEmptyGrid()
	return game
}

func (g Game) Show() string {
	grid := g.grid

	return showGrid(grid)
}

func (g Game) SetGrid(c Coord3D, v int) {
	g.grid[c] = v
}

func (g Game) GetGrid(c Coord3D) int {
	return g.grid[c]
}

func (g Game) Push(from Coord3D, direction Direction) error {
	if !IsValidCoord(from) {
		return errors.New(fmt.Sprintf("Invalid from coord: %v", from))
	}

	cellContent := g.grid[from]
	if cellContent == 0 {
		return errors.New("no marble to push")
	}

	destination := from.Add(direction)

	destinationContent, destinationExists := g.grid[destination]

	if destinationExists {
		if destinationContent != 0 {
			return errors.New("destination is not empty")
		}

		g.grid[destination] = cellContent
	} else {
		println("Captured a marble")
	}

	g.grid[from] = 0

	return nil
}
