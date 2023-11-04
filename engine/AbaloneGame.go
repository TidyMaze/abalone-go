package engine

import (
	"errors"
	"fmt"
	"log"
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

func (g Game) Push(from Coord3D, direction Direction, count int) error {
	if !IsValidCoord(from) {
		return errors.New(fmt.Sprintf("Invalid from coord: %v", from))
	}

	cells := findAllCells(from, direction)
	log.Println(cells)

	// check that there are between 1 and 3 marbles in the first cells,
	// followed by 0 marbles or an inferior number of enemy marbles

	myColor := g.grid[from]

	myFirstCells := []Coord3D{}
	nextEnemyCells := []Coord3D{}

	for _, cell := range cells {
		cellContent, cellExists := g.grid[cell]
		if !cellExists {
			break
		}

		if cellContent == 0 {
			break
		} else if cellContent == myColor && len(nextEnemyCells) == 0 {
			myFirstCells = append(myFirstCells, cell)
		} else if cellContent != myColor && len(myFirstCells) > 0 {
			nextEnemyCells = append(nextEnemyCells, cell)
		}
	}

	if len(myFirstCells) == 0 {
		return errors.New("no marble to push")
	} else if len(myFirstCells) > 3 {
		return errors.New(fmt.Sprintf("too many marbles to push (max 3, got %d)", len(myFirstCells)))
	} else if len(nextEnemyCells) > 0 && len(myFirstCells) <= len(nextEnemyCells) {
		return errors.New(fmt.Sprintf("not enough marbles to push enemy (got %d, need %d)", len(myFirstCells), len(nextEnemyCells)+1))
	}

	log.Println(fmt.Sprintf("Pushing my marbles: %v and enemy marbles: %v", myFirstCells, nextEnemyCells))

	// push enemy marbles in inverse order (from the last to the first)
	for i := len(nextEnemyCells) - 1; i >= 0; i-- {
		err := g.pushSingle(nextEnemyCells[i], direction)
		if err != nil {
			return err
		}
	}

	// push my marbles in reverse order (from the last to the first)
	for i := len(myFirstCells) - 1; i >= 0; i-- {
		err := g.pushSingle(myFirstCells[i], direction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g Game) pushSingle(from Coord3D, direction Direction) error {
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

func findAllCells(from Coord3D, direction Direction) []Coord3D {
	var cells []Coord3D

	currentCell := from

	// find all cells in the direction from the from cell, until we find a non-existent cell (out of the hexagon)
	for {
		cells = append(cells, currentCell)

		destination := currentCell.Add(direction)

		if !IsValidCoord(destination) {
			break
		}

		currentCell = destination
	}

	return cells

}

func (g Game) Copy() Game {
	newGame := Game{}
	newGame.grid = copyGrid(g.grid)
	return newGame
}

func copyGrid(grid map[Coord3D]int) map[Coord3D]int {
	newGrid := make(map[Coord3D]int)

	for k, v := range grid {
		newGrid[k] = v
	}

	return newGrid
}
