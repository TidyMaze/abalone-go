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
 * 0: empty
 * 1: player 1
 * 2: player 2
 *
 * The players are represented by their color.
 * 1: white
 * 2: black
 *
 * Grid layout:
 * - the grid is a 3x3 square
 */

type Game struct {
	grid          [3][3]int8 // 0: empty, 1: player 1, 2: player 2
	currentPlayer int8       // 1 or 2
	Turn          int8
	Winner        int8 // 0 if no Winner, 1 for tie, 2 or 3 if there is a Winner
}

var emptyGrid = buildEmptyGrid()
var startingGrid = buildStartingGrid()

func NewGame(grid [3][3]int8) *Game {
	game := &Game{
		currentPlayer: 1,
	}

	game.grid = copyGrid(grid)
	return game
}

func (g *Game) Show() string {
	res := ""

	res += fmt.Sprintf("Current player: %d\n", g.currentPlayer)
	res += fmt.Sprintf("Grid:\n%s", showGrid(g.grid))

	return res
}

func (g *Game) IsOver() bool {
	return g.Winner != 0
}

func (g *Game) SetGrid(c Coord2D, v int8) {
	g.grid[c.Y][c.X] = v
}

func (g *Game) GetGrid(c Coord2D) int8 {
	return g.grid[c.Y][c.X]
}

func (g *Game) Put(at Coord2D) error {
	err := g.checkCanPut(at)
	if err != nil {
		log.Println(fmt.Sprintf("Cannot put at %v: %s", at, err.Error()))
		return err
	}

	g.grid[at.Y][at.X] = g.currentPlayer

	g.currentPlayer = 3 - g.currentPlayer
	g.Turn += 1

	winner := g.checkWinner()
	if winner != 0 {
		g.Winner = winner + 1
		//log.Println(fmt.Sprintf("Winner: %d", g.Winner))
	}

	return nil
}

func (g *Game) checkWinner() int8 {
	// check if one line is full
	for i := 0; i < 3; i++ {
		if g.grid[i][0] != 0 && g.grid[i][0] == g.grid[i][1] && g.grid[i][1] == g.grid[i][2] {
			return g.grid[i][0]
		}

		if g.grid[0][i] != 0 && g.grid[0][i] == g.grid[1][i] && g.grid[1][i] == g.grid[2][i] {
			return g.grid[0][i]
		}

		if g.grid[0][0] != 0 && g.grid[0][0] == g.grid[1][1] && g.grid[1][1] == g.grid[2][2] {
			return g.grid[0][0]
		}

		if g.grid[0][2] != 0 && g.grid[0][2] == g.grid[1][1] && g.grid[1][1] == g.grid[2][0] {
			return g.grid[0][2]
		}
	}

	return 0
}

func (g *Game) checkCanPut(at Coord2D) error {
	if !IsValidCoord(at) {
		return errors.New(fmt.Sprintf("Invalid coord: %v", at))
	}

	if g.grid[at.Y][at.X] != 0 {
		return errors.New(fmt.Sprintf("Cannot play at %v: cell is not empty", at))
	}

	return nil
}

func (g *Game) Copy() *Game {
	newGame := NewGame(g.grid)
	newGame.currentPlayer = g.currentPlayer
	newGame.Turn = g.Turn
	newGame.Winner = g.Winner
	return newGame
}

func (g *Game) GetValidMoves() []Move {
	moves := make([]Move, 0)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			at := Coord2D{X: int8(j), Y: int8(i)}

			if g.grid[i][j] == 0 {
				moves = append(moves, Move{At: at})
			}
		}
	}

	return moves
}

func (g *Game) Move(move Move) error {
	return g.Put(move.At)
}

func copyGrid(grid [3][3]int8) [3][3]int8 {
	var newGrid [3][3]int8

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			newGrid[i][j] = grid[i][j]
		}
	}

	return newGrid
}
