package engine

import (
	. "abalone-go/helpers"
	"fmt"
)

func showGrid(grid map[Coord3D]int) string {
	res := ""

	// top left cell (-3, -3) is in (-3,3,6)
	// bottom right cell (3, 3) is in (3,-6,3)

	// for each horizontal line, print the line
	// x and y are in 2D coordinates
	for y := -3; y <= 3; y++ {
		// print the line
		for x := -3; x <= 3; x++ {
			coord := Coord2D{x, y}.To3D()

			if v, ok := grid[coord]; ok {
				res += fmt.Sprintf("%d ", v)
			} else {
				res += "  "
			}
		}

		res += "\n"
	}

	return res
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
				if _, ok := grid[neighbor]; !ok && IsValidCoord(neighbor) {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	AssertEqual(37, len(grid))

	return grid
}
