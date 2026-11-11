package day3

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code/advent"
	"github.com/sirgwain/advent-of-code/advent/color"
)

var correctStyle = lipgloss.NewStyle().Foreground(color.BrightGreen82)
var incorrectStyle = lipgloss.NewStyle().Foreground(color.BrightRed124)
var symbolStyle = lipgloss.NewStyle().Foreground(color.BlazeOrange202)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	board := advent.BoardFromBytes(input)

	var sb strings.Builder

	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if !advent.IsDigit(board[y][x]) {
				// output the non digit and continue
				if board[y][x] == '.' {
					sb.WriteByte(board[y][x])
				} else {
					sb.WriteString(symbolStyle.Render(string(board[y][x])))
				}
				continue
			}
			num, count := advent.Atoi(board[y][x:])

			// symbols adjacent or diagonal to a number
			// make it count
			// ...&...
			// ..123^.
			// .#...*.
			include := false
		testNumber:
			for y1 := -1; y1 <= 1; y1++ {
				for x1 := -1; x1 <= count; x1++ {
					if y1 == 0 && x1 >= 0 && x1 < count {
						// inside our number
						continue
					}
					val := advent.GetBoardValue(x+x1, y+y1, board)
					if val != 0 && val != '.' && !advent.IsDigit(val) {
						s1 += num
						include = true
						break testNumber
					}
				}
			}

			if include {
				sb.WriteString(correctStyle.Render(strconv.Itoa(num)))
			} else {
				sb.WriteString(incorrectStyle.Render(strconv.Itoa(num)))
			}

			// skip the number
			x += count - 1
		}
		sb.WriteRune('\n')
	}

	if advent.Viz.Enabled {
		fmt.Printf("\n%s\n", sb.String())
	}

	return s1, s2, nil
}
