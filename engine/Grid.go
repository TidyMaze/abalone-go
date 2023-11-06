package engine

import (
	. "abalone-go/helpers"
	"fmt"
)

func showGrid(grid map[Coord3D]int) string {
	res := ""

	// top left cell (-4, -4) is in (0,4,-4)
	// bottom right cell (4, 4) is in (0,-4,4)

	// for each horizontal line, print the line
	// x and y are in 2D coordinates
	for y := -4; y <= 4; y++ {
		line := ""

		// print the line
		for x := -4; x <= 4; x++ {
			coord := Coord2D{x, y}.To3D()

			if v, ok := grid[coord]; ok {
				line += fmt.Sprintf("%d ", v)
			} else {
				line += "  "
			}
		}

		switch y {
		case -4:
			line = line[4:]
		case -3:
			line = line[3:]
		case -2:
			line = line[2:]
		case -1:
			line = line[1:]
		case 0:
		case 1:
			line = " " + line
		case 2:
			line = "  " + line
		case 3:
			line = "   " + line
		case 4:
			line = "    " + line
		}

		res += fmt.Sprintf("%s\n", line)
	}

	return res
}

func buildEmptyGrid() map[Coord3D]int {
	grid := make(map[Coord3D]int)

	queue := []Coord3D{{0, 0, 0}}

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]

		if _, ok := grid[cell]; !ok {
			grid[cell] = 0

			// Add the 6 neighbors to the queue
			neighbors := []Coord3D{
				// top right
				{cell.X + 1, cell.Y, cell.Z - 1},
				// right
				{cell.X + 1, cell.Y - 1, cell.Z},
				// bottom right
				{cell.X, cell.Y - 1, cell.Z + 1},
				// bottom left
				{cell.X - 1, cell.Y, cell.Z + 1},
				// left
				{cell.X - 1, cell.Y + 1, cell.Z},
				// top left
				{cell.X, cell.Y + 1, cell.Z - 1},
			}

			for _, neighbor := range neighbors {
				if _, ok := grid[neighbor]; !ok && IsValidCoord(neighbor) {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	AssertEqual(61, len(grid))

	return grid
}

func buildStartingGrid() map[Coord3D]int {
	grid := buildEmptyGrid()

	// TODO

	return grid
}
