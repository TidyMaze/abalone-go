package engine

import (
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

type Coord struct {
	x int
	y int
	z int
}

type AbaloneGame struct {
	grid map[Coord]int
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

	// for each horizontal line, print the line
	for z := -3; z <= 3; z++ {
		// print the line
		for hCoord := -3; hCoord <= 3; hCoord++ {
			x := hCoord
			y := -hCoord - z

			if between(x, -3, 3) && between(y, -3, 3) && between(z, -3, 3) {
				res += fmt.Sprintf("%d ", g.grid[Coord{x, y, z}])
			} else {
				res += "_ "
			}
		}

		res += "\n"
	}

	return res
}

func buildEmptyGrid() map[Coord]int {
	grid := make(map[Coord]int)

	queue := []Coord{Coord{0, 0, 0}}

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]

		if _, ok := grid[cell]; !ok {
			grid[cell] = 0

			// Add the 6 neighbors to the queue
			neighbors := []Coord{
				// top right
				Coord{cell.x + 1, cell.y, cell.z - 1},
				// right
				Coord{cell.x + 1, cell.y - 1, cell.z},
				// bottom right
				Coord{cell.x, cell.y - 1, cell.z + 1},
				// bottom left
				Coord{cell.x - 1, cell.y, cell.z + 1},
				// left
				Coord{cell.x - 1, cell.y + 1, cell.z},
				// top left
				Coord{cell.x, cell.y + 1, cell.z - 1},
			}

			for _, neighbor := range neighbors {
				if _, ok := grid[neighbor]; !ok && isValidCoord(neighbor) {
					queue = append(queue, neighbor)
				}
			}
		}

		println(fmt.Sprintf("Queue size %d, grid size %d", len(queue), len(grid)))
	}

	return grid
}

func isValidCoord(c Coord) bool {
	return between(c.x, -3, 3) &&
		between(c.y, -3, 3) &&
		between(c.z, -3, 3) &&
		c.x+c.y+c.z == 0
}

func between(v int, min int, max int) bool {
	return v >= min && v <= max
}
