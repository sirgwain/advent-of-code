package day11

import (
	"fmt"
	"strconv"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

var boardMap = [256]int{
	'.': 0,
	'#': 1,
}

// additional growth (-1 because we grow once by default)
const galaxyGrowthFactor = 1_000_000 - 1

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	board := advent.BoardFromBytes(input)

	if advent.Viz.Enabled {
		count := 0
		fmt.Println(advent.BoardToString(board, func(b byte) string {
			if b == '#' {
				count++
				return strconv.Itoa(count)
			}
			return "."
		}))
	}

	emptyRowCounts := make([]int, len(board))
	emptyColCounts := make([]int, len(board[0]))
	galaxyCountPerCol := make([]int, len(board[0]))
	galaxyCount := 0

	// pre-compute the empty rows and empty cols
	// the idea is to keep a count of the empty rows/cols at any index, i.e.
	// 011222333 means no empty rows/cols at index 0, 1 at index 1+, 2 at index 3+, etc
	emptyRows := 0
	for y := range board {
		rowEmpty := 0
		for x, cell := range board[y] {
			rowEmpty += boardMap[cell]
			galaxyCountPerCol[x] += boardMap[cell]
			galaxyCount += boardMap[cell]
		}

		if rowEmpty == 0 {
			emptyRows++
		}
		emptyRowCounts[y] = emptyRows
	}
	emptyCols := 0
	for x, count := range galaxyCountPerCol {
		if count == 0 {
			emptyCols++
		}
		emptyColCounts[x] = emptyCols
	}

	galaxies := make([]advent.Point, 0, galaxyCount)

	// go through the board looking for galaxies
	// every time we find one, compute the distance between it and all previous
	// galaxies, adding expanding empty rows/cols between the two galaxies
	for y := range board {
		for x, cell := range board[y] {
			if boardMap[cell] > 0 {
				p := advent.Point{X: x, Y: y}
				galaxies = append(galaxies, p)
				for i := 0; i < len(galaxies)-1; i++ {
					d1 := dist(p, galaxies[i]) +
						advent.Abs(emptyRowCounts[y]-emptyRowCounts[galaxies[i].Y]) +
						advent.Abs(emptyColCounts[x]-emptyColCounts[galaxies[i].X])
					// for solution2, the distance for empty galaxies is multiplied (1 million times)
					d2 := dist(p, galaxies[i]) +
						(advent.Abs(emptyRowCounts[y]-emptyRowCounts[galaxies[i].Y])*galaxyGrowthFactor +
							advent.Abs(emptyColCounts[x]-emptyColCounts[galaxies[i].X])*galaxyGrowthFactor)

					if advent.Viz.Enabled {
						fmt.Printf("dist %d->%d = %d\n", i+1, len(galaxies), d1)
						fmt.Printf("dist2 %d->%d = %d\n", i+1, len(galaxies), d2)
					}
					s1 += d1
					s2 += d2
				}
			}
		}
	}

	return s1, s2, nil
}

func dist(p1, p2 advent.Point) int {
	return advent.Abs(p1.X-p2.X) + advent.Abs(p1.Y-p2.Y)
}
