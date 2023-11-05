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

func (g Game) Push(from Coord3D, direction Direction) (bool, error) {
	myFirstCells, nextEnemyCells, err := g.checkCanPush(from, direction)
	if err != nil {
		return false, err
	}

	log.Println(fmt.Sprintf("Pushing my marbles: %v and enemy marbles: %v", myFirstCells, nextEnemyCells))

	allCellsToPush := concatSlices(myFirstCells, nextEnemyCells)

	capturedMarble := false

	err = reverseForEach(allCellsToPush, func(cellToPush Coord3D) error {
		captured, err := g.pushSingle(cellToPush, direction)
		if err != nil {
			return err
		}

		if captured {
			capturedMarble = true
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return capturedMarble, nil
}

func (g Game) checkCanPush(from Coord3D, direction Direction) ([]Coord3D, []Coord3D, error) {
	if !IsValidCoord(from) {
		return nil, nil, errors.New(fmt.Sprintf("Invalid from coord: %v", from))
	}

	cells := findAllCells(from, direction)
	log.Println(cells)

	// check that there are between 1 and 3 marbles in the first cells,
	// followed by 0 marbles or an inferior number of enemy marbles

	myColor := g.grid[from]

	var myFirstCells []Coord3D
	var nextEnemyCells []Coord3D

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
		} else if cellContent == myColor && len(nextEnemyCells) > 0 {
			return nil, nil, errors.New("my marbles are sandwiching enemy marbles")
		}
	}

	if len(myFirstCells) == 0 {
		return nil, nil, errors.New("no marble to push")
	} else if len(myFirstCells) > 3 {
		return nil, nil, errors.New(fmt.Sprintf("too many marbles to push (max 3, got %d)", len(myFirstCells)))
	} else if len(nextEnemyCells) >= 3 {
		return nil, nil, errors.New(fmt.Sprintf("too many enemy marbles to push (max 2, got %d)", len(nextEnemyCells)))
	} else if len(nextEnemyCells) > 0 && len(myFirstCells) <= len(nextEnemyCells) {
		return nil, nil, errors.New(fmt.Sprintf("not enough marbles to push enemy (got %d, need %d)", len(myFirstCells), len(nextEnemyCells)+1))
	}
	return myFirstCells, nextEnemyCells, nil
}

func (g Game) pushSingle(from Coord3D, direction Direction) (bool, error) {
	cellContent := g.grid[from]
	if cellContent == 0 {
		return false, errors.New("no marble to push")
	}

	destination := from.Add(direction)

	destinationContent, destinationExists := g.grid[destination]

	capturedMarble := false

	if destinationExists {
		if destinationContent != 0 {
			return false, errors.New(fmt.Sprintf("destination is not empty: %v", destination))
		}

		g.grid[destination] = cellContent
	} else {
		println("Captured a marble")
		capturedMarble = true
	}

	g.grid[from] = 0
	return capturedMarble, nil
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

func concatSlices[T any](slices ...[]T) []T {
	var result []T

	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}

func reverseForEach[T any](slice []T, f func(T) error) error {
	for i := len(slice) - 1; i >= 0; i-- {
		err := f(slice[i])
		if err != nil {
			return err
		}
	}

	return nil
}
