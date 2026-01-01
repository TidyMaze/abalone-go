package engine

import "fmt"

func showGrid(grid [3][3]int8) string {
	res := ""

	for i := 0; i < 3; i++ {
		res += "  "
		for j := 0; j < 3; j++ {
			res += fmt.Sprintf("%d ", grid[i][j])
		}
		res += "\n"
	}

	return res
}

func buildEmptyGrid() [3][3]int8 {
	return [3][3]int8{}
}

func buildStartingGrid() [3][3]int8 {
	return copyGrid(emptyGrid)
}
