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
var gearPairStyle = lipgloss.NewStyle().Foreground(color.PureMagenta200)
var incorrectStyle = lipgloss.NewStyle().Foreground(color.BrightRed124)
var symbolStyle = lipgloss.NewStyle().Foreground(color.BlazeOrange202)

type Day struct {
}

type gearPair struct {
	pt       advent.Point
	num1     int
	num2     int
	included bool // true if this gear pair should be included
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	board := advent.BoardFromBytes(input)

	gears := map[advent.Point]*gearPair{}

	var sb strings.Builder

	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if !advent.IsDigit(board[y][x]) {
				// output the non digit and continue
				if board[y][x] == '.' {
					sb.WriteByte(board[y][x])
				} else {
					if board[y][x] == '*' {
						sb.WriteString(gearPairStyle.Render(string(board[y][x])))
					} else {
						sb.WriteString(symbolStyle.Render(string(board[y][x])))
					}
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
			for y1 := -1; y1 <= 1; y1++ {
				for x1 := -1; x1 <= count; x1++ {
					if y1 == 0 && x1 >= 0 && x1 < count {
						// inside our number
						continue
					}
					pt := advent.Point{X: x + x1, Y: y + y1}
					val := advent.GetBoardValue(pt.X, pt.Y, board)
					if val != 0 && val != '.' && !advent.IsDigit(val) {
						include = true
					}
					
					// check for gear pairs
					if val == '*' {
						if existingPair, ok := gears[pt]; ok {
							if existingPair.num2 != 0 {
								existingPair.included = false
								continue
							}
							existingPair.num2 = num
							existingPair.included = true
						} else {
							gears[pt] = &gearPair{
								pt:   pt,
								num1: num,
							}
						}
					}
				}
			}

			if include {
				s1 += num
				sb.WriteString(correctStyle.Render(strconv.Itoa(num)))
			} else {
				sb.WriteString(incorrectStyle.Render(strconv.Itoa(num)))
			}

			// skip the number
			x += count - 1
		}
		sb.WriteRune('\n')
	}

	for _, pair := range gears {
		if pair.included {
			s2 += pair.num1 * pair.num2
		}
	}

	if advent.Viz.Enabled {
		fmt.Printf("\n%s\n", sb.String())
	}

	return s1, s2, nil
}
