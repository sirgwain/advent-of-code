package day10

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code/advent"
	"github.com/sirgwain/advent-of-code/advent/color"
)

var correctStyle = lipgloss.NewStyle().Foreground(color.BrightGreen82)
var incorrectStyle = lipgloss.NewStyle().Foreground(color.BrightRed124)
var startStyle = lipgloss.NewStyle().Foreground(color.BlazeOrange202)
var light = lipgloss.NewStyle().Foreground(color.BrightGreen82).Render("●")

type dir uint

const (
	dirN       dir = 1 << 1  // 0b000_000_010
	dirW       dir = 1 << 3  // 0b000_001_000
	dirCenter  dir = 1 << 4  // 0b000_010_000
	dirE       dir = 1 << 5  // 0b000_100_000
	dirS       dir = 1 << 7  // 0b010_000_000
	dirLoop    dir = 1 << 9  // flag, outside the 3x3 grid
	dirOutLoop dir = 1 << 10 // flag, outside the 3x3 grid
	dirInLoop  dir = 1 << 11 // flag, outside the 3x3 grid
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
			switch board[y][x] {
			case bNS:
				dirBoard[y][x] = dirNS
			case bEW:
				dirBoard[y][x] = dirEW
			case bNE:
				dirBoard[y][x] = dirNE
			case bNW:
				dirBoard[y][x] = dirNW
			case bSW:
				dirBoard[y][x] = dirSW
			case bSE:
				dirBoard[y][x] = dirSE
			case bStart:
				dirBoard[y][x] = dirAll
			default:
				dirBoard[y][x] = 0
			}
		}
	}
	sx, sy := advent.FindValue(dirBoard, dirAll)

	dirs := [4]dir{dirN, dirE, dirS, dirW}
	for _, d := range dirs {
		steps := followPipe(d.opposite(), 1, sx+d.xoff(), sy+d.yoff(), dirBoard)
		if steps > 0 {
			s1 = steps / 2
			break
		}
	}

	if advent.Viz.Enabled {
		fmt.Println(view(dirBoard))
	}

	expandedBoard := expandBoard(dirBoard)
	floodFill(advent.Point{X: 0, Y: 0}, expandedBoard)

	for y := range dirBoard {
		for x := 0; x < len(dirBoard[y]); x++ {
			if (dirBoard[y][x] & dirLoop) > 0 {
				continue
			}
			allZero := true
		zeroCheck:
			for y1 := range 3 {
				for x1 := range 3 {
					if expandedBoard[y*3+y1][x*3+x1] != 0 {
						allZero = false
						break zeroCheck
					}
				}
			}
			if allZero {
				dirBoard[y][x] |= dirInLoop
				s2++
			} else {
				dirBoard[y][x] |= dirOutLoop
			}
		}
	}

	if advent.Viz.Enabled {
		fmt.Println(view(dirBoard))

		// drawBoardCell := func(c byte) string {
		// 	if c == 0 {
		// 		return "."
		// 	}
		// 	if c == '#' {
		// 		return light
		// 	}
		// 	return string(c)
		// }
		// fmt.Println(advent.BoardToString(expandedBoard, drawBoardCell))
	}

	return s1, s2, nil
}

// followPipe recursively follows a pipe from a direction through its end
// if the pipe has no input from the from direction, return 0
func followPipe(from dir, step, sx, sy int, board [][]dir) int {
	pipe := advent.GetBoardValue(sx, sy, board)

	if pipe == dirAll {
		// we reached the start again, drop
		return step
	}

	// if our current pipe has no incoming from the from dir
	// i.e. we are coming from the S and this pipe has no S part, we are done
	if (pipe & from) == 0 {
		return 0
	}

	// if the pipe is NS and we are coming from S, find the N endpoint
	// AND NOT bitwise op
	// 0b0101 (NS) &^ 0b0100 (S) = 0b0001 (N)
	end := pipe &^ from

	result := followPipe(end.opposite(), step+1, sx+end.xoff(), sy+end.yoff(), board)
	if result > 0 {
		board[sy][sx] |= dirLoop
	}

	return result
}

// expandBoard expands a resolved board into a byte board for filling
// ┐ becomes this 3x3 grid
//
//	012
//
// 0 ...
// 1 ##.
// 2 .#.
func expandBoard(board [][]dir) [][]byte {
	out := advent.MakeBoard[byte](len(board[0])*3, len(board)*3)

	for y := range board {
		for x := 0; x < len(board[y]); x++ {
			// clear this grid
			if (board[y][x] & dirLoop) == 0 {
				continue
			}

			// center is always a wall for edges
			out[y*3+1][x*3+1] = '#'

			val := (board[y][x] &^ dirLoop)
			for y1 := range 3 {
				for x1 := range 3 {
					if (val & (1 << (y1*3 + x1))) > 0 {
						out[y*3+y1][x*3+x1] = '#'
					}
				}
			}
		}
	}
	return out
}

// floodFill fills all empty spaces from a start point
func floodFill(start advent.Point, board [][]byte) {
	work := []advent.Point{start}
	w := len(board[0])
	h := len(board)
	for len(work) > 0 {
		p := work[len(work)-1]
		work = work[:len(work)-1]

		if p.X < 0 || p.X >= w || p.Y < 0 || p.Y >= h {
			continue
		}
		if board[p.Y][p.X] != 0 {
			continue
		}

		board[p.Y][p.X] = 'O'
		work = append(work,
			advent.Point{X: p.X, Y: p.Y + 1},
			advent.Point{X: p.X + 1, Y: p.Y},
			advent.Point{X: p.X, Y: p.Y - 1},
			advent.Point{X: p.X - 1, Y: p.Y},
		)
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
				sb.WriteString(light)
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
