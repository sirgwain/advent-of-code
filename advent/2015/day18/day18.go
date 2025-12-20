package day18

import (
	"bufio"
	"bytes"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code/advent"
	"github.com/sirgwain/advent-of-code/advent/color"
)

type Day struct {
}

const (
	on  byte = '#'
	off byte = '.'
)

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input)))
	board := make([][]byte, 0)
	for scanner.Scan() {
		b := scanner.Bytes()
		row := make([]byte, len(b))
		copy(row, b)
		board = append(board, row)
	}

	board1 := advent.CloneBoard(board)
	toggleBoardLights(100, board1, false)
	s1 = countLights(board1)

	// for part 2, turn on all the corners
	board[0][0] = on
	board[0][len(board[0])-1] = on
	board[len(board)-1][0] = on
	board[len(board)-1][len(board[0])-1] = on
	toggleBoardLights(100, board, true)
	s2 = countLights(board)

	return s1, s2, nil
}

var light = lipgloss.NewStyle().Foreground(color.BrightGreen82).Render("●")

func drawBoardCell(c byte) string {
	if c == '#' {
		return light
	}
	return "."
}

func toggleBoardLights(times int, board [][]byte, cornersOn bool) {
	// fmt.Printf("initial board: \n%s\n\n", advent.BoardToString(board, drawBoardCell))
	emptyBoard := make([][]byte, len(board))
	for y := 0; y < len(board); y++ {
		emptyBoard[y] = make([]byte, len(board[y]))
		for x := 0; x < len(board[y]); x++ {
			emptyBoard[y][x] = off
		}
	}

	// make two boards to stop allocs
	boardA := board
	boardB := advent.MakeBoard[byte](len(board[0]), len(board))

	for range times {

		// make boardB empty
		advent.CopyBoard(boardB, emptyBoard)

		for y := 0; y < len(boardA); y++ {
			for x := 0; x < len(boardA[y]); x++ {
				if cornersOn && (y == 0 || y == len(boardA)-1) && (x == 0 || x == len(boardA[y])-1) {
					boardB[y][x] = on
					continue
				}
				// count how many neighbors are on
				neighborsOn := 0
				if advent.GetBoardValue(x, y-1, boardA) == on {
					neighborsOn++
				}
				if advent.GetBoardValue(x, y+1, boardA) == on {
					neighborsOn++
				}
				if advent.GetBoardValue(x-1, y, boardA) == on {
					neighborsOn++
				}
				if advent.GetBoardValue(x+1, y, boardA) == on {
					neighborsOn++
				}
				if advent.GetBoardValue(x-1, y-1, boardA) == on {
					neighborsOn++
				}
				if advent.GetBoardValue(x+1, y+1, boardA) == on {
					neighborsOn++
				}
				if advent.GetBoardValue(x+1, y-1, boardA) == on {
					neighborsOn++
				}
				if advent.GetBoardValue(x-1, y+1, boardA) == on {
					neighborsOn++
				}

				// if this value is off and 3 neighbors are on, turn it on
				if advent.GetBoardValue(x, y, boardA) == off && neighborsOn == 3 {
					// turn it on
					boardB[y][x] = on
				} else if advent.GetBoardValue(x, y, boardA) == on && (neighborsOn == 3 || neighborsOn == 2) {
					boardB[y][x] = on
				}
			}
		}

		// copy boardB back to boardA
		advent.CopyBoard(boardA, boardB)

		// swap boardA and B
		boardA, boardB = boardB, boardA
		// fmt.Printf("board: \n%s\n\n", advent.BoardToString(boardA, drawBoardCell))
	}

	advent.CopyBoard(board, boardA)
}

func countLights(board [][]byte) (count int) {
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if board[y][x] == on {
				count++
			}
		}
	}
	return count
}
