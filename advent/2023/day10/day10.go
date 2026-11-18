package day10

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code/advent"
	"github.com/sirgwain/advent-of-code/advent/color"
)

var correctStyle = lipgloss.NewStyle().Foreground(color.BrightGreen82)
var startStyle = lipgloss.NewStyle().Foreground(color.BlazeOrange202)
var fill = lipgloss.NewStyle().Foreground(color.BrightGreen82).Render("●")

type dir byte

const (
	dirN       dir = 1 << 1
	dirW       dir = 1 << 2
	dirE       dir = 1 << 3
	dirS       dir = 1 << 4
	dirLoop    dir = 1 << 5
	dirOutLoop dir = 1 << 6
	dirInLoop  dir = 1 << 7
	dirNS          = dirN | dirS
	dirEW          = dirE | dirW
	dirNE          = dirN | dirE
	dirNW          = dirN | dirW
	dirSW          = dirS | dirW
	dirSE          = dirS | dirE
	dirAll         = dirN | dirE | dirS | dirW | dirLoop
)

const (
	bNS    byte = '|' // │
	bEW    byte = '-' // ─
	bNE    byte = 'L' // └
	bNW    byte = 'J' // ┘
	bSW    byte = '7' // ┐
	bSE    byte = 'F' // ┌
	bStart byte = 'S' // ┼
)

// mapping of byte to directions
var pipeDirs = [256]dir{
	bNS:    dirNS,
	bEW:    dirEW,
	bNE:    dirNE,
	bNW:    dirNW,
	bSW:    dirSW,
	bSE:    dirSE,
	bStart: dirAll,
}

func (d dir) yoff() int {
	switch d {
	case dirN:
		return -1
	case dirS:
		return 1
	}
	return 0
}

func (d dir) xoff() int {
	switch d {
	case dirE:
		return 1
	case dirW:
		return -1
	}
	return 0
}

func (d dir) opposite() dir {
	switch d {
	case dirN:
		return dirS
	case dirE:
		return dirW
	case dirS:
		return dirN
	case dirW:
		return dirE
	}
	return d
}

func (d dir) String() string {
	switch d {
	case 0:
		return "."
	// single cardinal directions
	case dirN:
		return "↑"
	case dirE:
		return "→"
	case dirS:
		return "↓"
	case dirW:
		return "←"
	// pipes (two directions)
	case dirNS:
		return "│"
	case dirEW:
		return "─"
	case dirNE:
		return "└"
	case dirNW:
		return "┘"
	case dirSW:
		return "┐"
	case dirSE:
		return "┌"
	// start
	case dirAll:
		return "┼"
	}
	return "?"
}

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	board := advent.BoardFromBytes(input)
	dirBoard := advent.MakeBoard[dir](len(board[0]), len(board))

	// make a copy of the board with directions
	for y := range board {
		for x := 0; x < len(board[y]); x++ {
			dirBoard[y][x] = pipeDirs[board[y][x]]
		}
	}
	sx, sy := advent.FindValue(dirBoard, dirAll)

	// find the real shape of the starting point
	dirs := [4]dir{dirN, dirE, dirS, dirW}
	var startShape dir
	for _, d := range dirs {
		// if the adjacent node has an incoming edge to match, include that direction in the starting shape
		if (advent.GetBoardValue(sx+d.xoff(), sy+d.yoff(), dirBoard) & d.opposite()) > 0 {
			startShape |= d
		}
	}

	// the starting shape is always part of the board
	dirBoard[sy][sx] = startShape | dirLoop

	// halfway point along the loop is solution 1
	s1 = followPipe(sx, sy, dirBoard) / 2

	for y := range dirBoard {
		// imagine a ray shooting just above the middle of the board square
		// every north edge (│└┘) flips whether you are inside or outside the loop polygon
		inside := false
		for x := 0; x < len(dirBoard[y]); x++ {
			if (dirBoard[y][x] & dirLoop) > 0 {
				// loop piece with a north edge, flip inside/outside flag
				if (dirBoard[y][x] & dirN) > 0 {
					inside = !inside
				}
				continue
			}

			if inside {
				// every non loop square inside the polygon is counted for solution 2
				s2++
				dirBoard[y][x] |= dirInLoop
			} else {
				dirBoard[y][x] |= dirOutLoop
			}
		}
	}

	if advent.Viz.Enabled {
		fmt.Println(view(dirBoard))
	}

	return s1, s2, nil
}

// followPipe recursively follows a pipe from a direction through its end
// if the pipe has no input from the from direction, return 0
func followPipe(sx, sy int, board [][]dir) int {
	heading := board[sy][sx] & -board[sy][sx]

	x, y := sx, sy
	steps := 0
	for {
		// check the next square
		x, y = x+heading.xoff(), y+heading.yoff()
		steps++

		if x == sx && y == sy {
			// back to the start, all done!
			return steps
		}

		pipe := board[y][x]
		board[y][x] |= dirLoop

		// if the pipe is NS and we are heading N (coming from S), find the N endpoint as our new heading
		// AND NOT bitwise op
		// 0b0101 (NS) &^ 0b0100 (S) = 0b0001 (N)
		heading = pipe &^ heading.opposite()
	}
}

// view returns the board as a colored pretty string like
// .....
// .┼─┐.
// .│.│.
// .└─┘.
// .....
func view(board [][]dir) string {
	var sb strings.Builder

	for y := range board {
		for x := 0; x < len(board[y]); x++ {
			if board[y][x] == dirAll {
				sb.WriteString(startStyle.Render("┼"))
				continue
			}

			if (board[y][x] & dirLoop) > 0 {
				sb.WriteString(correctStyle.Render((board[y][x] &^ dirLoop).String()))
				continue
			}

			if (board[y][x] & dirInLoop) > 0 {
				sb.WriteString(fill)
				continue
			}

			if (board[y][x] & dirOutLoop) > 0 {
				sb.WriteRune('.')
			} else {
				sb.WriteString(board[y][x].String())
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
