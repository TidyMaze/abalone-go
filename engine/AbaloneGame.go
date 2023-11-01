package engine

import (
	"abalone-go/helpers"
	"fmt"
)

/**
 * AbaloneGame is the main game class.
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

type Coord3D struct {
	X int
	Y int
	Z int
}

type Coord2D struct {
	X int
	Y int
}

func (c Coord2D) To3D() Coord3D {
	return Coord3D{c.X, -c.X - c.Y, c.Y}
}

func (c Coord3D) To2D() Coord2D {
	return Coord2D{c.X, c.Y}
}

type AbaloneGame struct {
	grid map[Coord3D]int
}

func (g AbaloneGame) show() {
}

func NewGame() *AbaloneGame {
	game := &AbaloneGame{}
	game.grid = buildEmptyGrid()
	return game
}

func (g AbaloneGame) Show() string {
	res := ""

	// top left cell (-3, -3) is in (-3,3,6)
	// bottom right cell (3, 3) is in (3,-6,3)

	// for each horizontal line, print the line
	// x and y are in 2D coordinates
	for y := -3; y <= 3; y++ {
		// print the line
		for x := -3; x <= 3; x++ {
			coord := Coord2D{x, y}.To3D()
			if v, ok := g.grid[coord]; ok {
				res += fmt.Sprintf("%d ", v)
			} else {
				res += "  "
			}
		}

		res += "\n"
	}

	return res
}

func (g AbaloneGame) SetGrid(c Coord3D, v int) {
	g.grid[c] = v
}

func buildEmptyGrid() map[Coord3D]int {
	grid := make(map[Coord3D]int)

	queue := []Coord3D{Coord3D{0, 0, 0}}

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]

		if _, ok := grid[cell]; !ok {
			grid[cell] = 0

			// Add the 6 neighbors to the queue
			neighbors := []Coord3D{
				// top right
				Coord3D{cell.X + 1, cell.Y, cell.Z - 1},
				// right
				Coord3D{cell.X + 1, cell.Y - 1, cell.Z},
				// bottom right
				Coord3D{cell.X, cell.Y - 1, cell.Z + 1},
				// bottom left
				Coord3D{cell.X - 1, cell.Y, cell.Z + 1},
				// left
				Coord3D{cell.X - 1, cell.Y + 1, cell.Z},
				// top left
				Coord3D{cell.X, cell.Y + 1, cell.Z - 1},
			}

			for _, neighbor := range neighbors {
				if _, ok := grid[neighbor]; !ok && isValidCoord(neighbor) {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	helpers.AssertEqual(37, len(grid))

	return grid
}

func isValidCoord(c Coord3D) bool {
	return between(c.X, -3, 3) &&
		between(c.Y, -3, 3) &&
		between(c.Z, -3, 3) &&
		c.X+c.Y+c.Z == 0
}

func between(v int, min int, max int) bool {
	return v >= min && v <= max
}
