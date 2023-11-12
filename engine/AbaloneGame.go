package engine

import (
	"errors"
	"fmt"
	"sort"
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
 * The hexagon is 5 cells wide in the outer ring (side size) and 9 cells wide in the inner ring (diameter).
 * There are 6 directions expressed in (x,y,z) coordinates:
 * - top right (1, 0, -1)
 * - right (1, -1, 0)
 * - bottom right (0, -1, 1)
 * - bottom left (-1, 0, 1)
 * - left (-1, 1, 0)
 * - top left (0, 1, -1)
 */

type Game struct {
	grid          map[Coord3D]int
	score         map[int]int
	currentPlayer int
	Turn          int
	Winner        int // 0 if no Winner, 1 or 2 if there is a Winner
}

var emptyGrid = buildEmptyGrid()
var startingGrid = buildStartingGrid()

func NewGame(grid *map[Coord3D]int) *Game {
	game := &Game{
		score:         make(map[int]int),
		currentPlayer: 1,
	}

	game.score[1] = 0
	game.score[2] = 0

	game.grid = copyGrid(grid)
	return game
}

func (g *Game) Show() string {
	res := ""

	res += fmt.Sprintf("Current player: %d\n", g.currentPlayer)
	res += fmt.Sprintf("Score: %v\n", g.score)
	res += fmt.Sprintf("Grid:\n%s", showGrid(g.grid))

	return res
}

func (g *Game) IsOver() bool {
	return g.Winner != 0
}

func (g *Game) SetGrid(c Coord3D, v int) {
	g.grid[c] = v
}

func (g *Game) GetGrid(c Coord3D) int {
	return g.grid[c]
}

func (g *Game) Push(from Coord3D, direction Direction) error {
	myFirstCells, nextEnemyCells, err := g.checkCanPush(from, direction)
	if err != nil {
		return err
	}

	//log.Println(fmt.Sprintf("Pushing my marbles: %v and enemy marbles: %v", myFirstCells, nextEnemyCells))

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
		return err
	}

	if capturedMarble {
		g.score[g.currentPlayer] += 1
	}

	//log.Println(fmt.Sprintf("Switching player from %d to %d", g.currentPlayer, 3-g.currentPlayer))
	g.currentPlayer = 3 - g.currentPlayer
	g.Turn += 1

	if g.score[1] == 6 {
		g.Winner = 1
	} else if g.score[2] == 6 {
		g.Winner = 2
	}

	return nil
}

func (g *Game) checkCanPush(from Coord3D, direction Direction) ([]Coord3D, []Coord3D, error) {
	if !IsValidCoord(from) {
		return nil, nil, errors.New(fmt.Sprintf("Invalid from coord: %v", from))
	}

	if g.currentPlayer != g.grid[from] {
		return nil, nil, errors.New(fmt.Sprintf("Cannot push marble from %v: it is not the current player's marble (current player: %d, marble color: %d)", from, g.currentPlayer, g.grid[from]))
	}

	// check that there are between 1 and 3 marbles in the first cells,
	// followed by 0 marbles or an inferior number of enemy marbles

	var myFirstCells []Coord3D
	var nextEnemyCells []Coord3D

	currentCell := from

	// find all cells in the direction from the from cell, until we find a non-existent cell (out of the hexagon)
	for {
		cellContent := g.grid[currentCell]

		if cellContent == 0 {
			break
		} else if cellContent == g.currentPlayer && len(nextEnemyCells) == 0 {
			myFirstCells = append(myFirstCells, currentCell)
		} else if cellContent != g.currentPlayer && len(myFirstCells) > 0 {
			nextEnemyCells = append(nextEnemyCells, currentCell)
		} else if cellContent == g.currentPlayer && len(nextEnemyCells) > 0 {
			return nil, nil, errors.New("my marbles are sandwiching enemy marbles")
		}

		destination := currentCell.Add(direction)

		if !IsValidCoord(destination) {
			break
		}

		currentCell = destination
	}

	nextCellAfterMyMarbles := myFirstCells[len(myFirstCells)-1].Add(direction)

	if !IsValidCoord(nextCellAfterMyMarbles) && len(nextEnemyCells) == 0 {
		return nil, nil, errors.New("cannot push its own marbles out of the hexagon")
	} else if len(myFirstCells) == 0 {
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

func (g *Game) pushSingle(from Coord3D, direction Direction) (bool, error) {
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
		//log.Println(fmt.Sprintf("Player %d has captured a marble in %v", g.currentPlayer, destination))
		capturedMarble = true
	}

	g.grid[from] = 0
	return capturedMarble, nil
}

func (g *Game) Copy() *Game {
	newGame := NewGame(&g.grid)
	newGame.score = copyScore(g.score)
	newGame.currentPlayer = g.currentPlayer
	newGame.Turn = g.Turn
	newGame.Winner = g.Winner
	return newGame
}

func (g *Game) GetValidMoves() []Move {
	moves := make([]Move, 0)

	pushLines := make([]PushLine, 0)

	for coord := range g.grid {
		if g.grid[coord] == g.currentPlayer {
			for _, direction := range Directions {
				move := PushLine{From: coord, Direction: direction}
				_, _, err := g.checkCanPush(coord, direction)
				if err == nil {
					pushLines = append(pushLines, move)
				}
			}
		}
	}

	// sort moves by 2D coordinate (by y then by x) then by direction (top right, right, bottom right, bottom left, left, top left)
	// this is to make sure that the moves are always displayed in the same order
	// this is important for the AI to be deterministic

	sort.Slice(pushLines, func(i, j int) bool {
		iCellIn2D := pushLines[i].From.To2D()
		jCellIn2D := pushLines[j].From.To2D()

		//log.Println(fmt.Sprintf("iCell %v in 2D: %v", pushLines[i].From, iCellIn2D))
		//log.Println(fmt.Sprintf("jCell %v in 2D: %v", pushLines[j].From, jCellIn2D))

		if iCellIn2D.Y == jCellIn2D.Y {
			if iCellIn2D.X == jCellIn2D.X {
				return pushLines[i].Direction < pushLines[j].Direction
			} else {
				return iCellIn2D.X < jCellIn2D.X
			}
		} else {
			return iCellIn2D.Y < jCellIn2D.Y
		}
	})

	for _, pushLine := range pushLines {
		moves = append(moves, pushLine)
	}

	return moves
}

func (g *Game) Move(move Move) error {
	switch move.(type) {
	case PushLine:
		return g.Push(move.(PushLine).From, move.(PushLine).Direction)
	default:
		return errors.New(fmt.Sprintf("Unknown move type: %v", move))
	}
}

func copyScore(score map[int]int) map[int]int {
	newScore := make(map[int]int)

	for k, v := range score {
		newScore[k] = v
	}

	return newScore
}

func copyGrid(grid *map[Coord3D]int) map[Coord3D]int {
	newGrid := make(map[Coord3D]int)

	for k, v := range *grid {
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
