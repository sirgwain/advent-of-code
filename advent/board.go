package advent

import (
	"bufio"
	"bytes"
	"strings"
)

func ValidPosition(p Point, width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
}

func MakeBoard[T any](width, height int) [][]T {
	board := make([][]T, height)
	for y := range height {
		board[y] = make([]T, width)
	}
	return board
}

func CloneBoard[T any](source [][]T) [][]T {
	board := make([][]T, len(source))
	for y := range len(source) {
		board[y] = make([]T, len(source[y]))
		copy(board[y], source[y])
	}
	return board
}

func CopyBoard[T any](dest, source [][]T) {
	for y := range len(source) {
		copy(dest[y], source[y])
	}
}

// GetBoardValue returns a rune/int/bool at x,y in the input or the empty value if out of bounds
func GetBoardValue[T int | uint | byte | rune | bool](x, y int, board [][]T) T {
	var zero T
	if y < 0 || y >= len(board) {
		return zero
	}
	if x < 0 || x >= len(board[y]) {
		return zero
	}
	return board[y][x]
}

// find
func FindValue[T int | uint | byte | rune | bool](board [][]T, c T) (x, y int) {
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if board[y][x] == c {
				return x, y
			}
		}
	}
	return 0, 0
}

// read input as a series of rune lines
func BoardFromBytes(input []byte) [][]byte {
	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input)))
	lines := make([][]byte, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Bytes())
	}

	return lines
}

func BoardToString[T any](board [][]T, cell func(T) string) string {
	var sb strings.Builder

	// Optional: pre-grow (rough estimate; assumes rectangular)
	if h := len(board); h > 0 {
		sb.Grow(h * (len(board[0]) + 1))
	}

	for y := range board {
		for x := range board[y] {
			sb.WriteString(cell(board[y][x]))
		}
		if y+1 < len(board) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
