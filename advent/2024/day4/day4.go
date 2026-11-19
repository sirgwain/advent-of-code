package day4

import "github.com/sirgwain/advent-of-code/advent"

type Day struct {
}

var (
	xmasBytes = []byte("XMAS")
)

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	board := advent.BoardFromBytes(input)

	directions := [][2]int{
		{0, -1},  // up
		{1, -1},  // up-right
		{1, 0},   // right
		{1, 1},   // down-right
		{0, 1},   // down
		{-1, 1},  // down-left
		{-1, 0},  // left
		{-1, -1}, // up-left
	}

	for y := range board {
		for x := 0; x < len(board[y]); x++ {
			switch board[y][x] {
			case 'X':
				for _, dir := range directions {
					if findXMAS(x, y, dir, board) {
						s1++
					}
				}
			case 'A':
				// skip edges, we need to look backwards
				if y == 0 || x == 0 || y == len(board)-1 || x == len(board[y])-1 {
					continue
				}
				s2 += findMAS(x, y, board)
			}
		}
	}

	return s1, s2, nil
}

// findXMAS finds the XMAS string on a board in a given direction
func findXMAS(x, y int, dir [2]int, board [][]byte) bool {
	for i := range 3 {
		x, y = x+dir[0], y+dir[1]
		if advent.GetBoardValue(x, y, board) != xmasBytes[i+1] {
			return false
		}
	}
	return true
}

// findMAS finds X shaped MAS text like
// ====
// M.M
// .A.
// S.S
//
// right
// ====
// S.M
// .A.
// S.M
// ====
// M.M
// .A.
// S.S
//
// right
// ====
// S.M
// .A.
// S.M
//
// down
// ====
// S.S
// .A.
// M.M
//
// left
// ====
// M.S
// .A.
// M.S
func findMAS(x, y int, board [][]byte) int {
	matches := 0
	// check top M, bottom S
	if board[y-1][x-1] == 'M' && board[y-1][x+1] == 'M' &&
		board[y+1][x-1] == 'S' && board[y+1][x+1] == 'S' {
		matches++
	}
	// check right M, left S
	if board[y-1][x-1] == 'S' && board[y-1][x+1] == 'M' &&
		board[y+1][x-1] == 'S' && board[y+1][x+1] == 'M' {
		matches++
	}
	// check bottom M, top S
	if board[y-1][x-1] == 'S' && board[y-1][x+1] == 'S' &&
		board[y+1][x-1] == 'M' && board[y+1][x+1] == 'M' {
		matches++
	}
	// check left M, right S
	if board[y-1][x-1] == 'M' && board[y-1][x+1] == 'S' &&
		board[y+1][x-1] == 'M' && board[y+1][x+1] == 'S' {
		matches++
	}

	return matches
}
