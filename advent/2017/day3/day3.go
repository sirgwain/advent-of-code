package day3

import (
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	num, _ := advent.Atoi(input)

	// build a spiral grid like this
	// 17  16  15  14  13
	// 18   5   4   3  12
	// 19   6   1   2  11
	// 20   7   8   9  10
	// 21  22  23---> ...
	//
	// from 1 we go over 1, up 1, over 2, up 2, over 3, up 3

	x, y := 0, 0
	dist := 1
	hx, hy := 1, 0
	steps := 1
	i := 0
	for steps < num {
		if advent.Viz.Enabled {
			fmt.Printf("%d: (%d, %d) dist: %d, h: (%d, %d)\n", steps, x, y, dist, hx, hy)
		}

		// travel the dist
		x += hx * dist
		y += hy * dist
		hx, hy = -hy, hx
		steps += dist
		if i%2 == 1 {
			// increase dist every other turn
			dist++
		}
		i++

		if steps+dist >= num {
			dist = num - steps
		}
	}

	s1 = advent.Abs(x) + advent.Abs(y)

	// part 2, generate the board
	dirs := [][2]int{
		{0, 1}, // up
		{1, 1},
		{1, 0},  // right
		{1, -1}, // down right
		{0, -1},
		{-1, -1},
		{-1, 0},
		{-1, 1}, // up left
	}
	board := advent.MakeBoard[int](20, 20)
	x, y = 10, 10
	board[x][y] = 1
	dist = 1
	hx, hy = 1, 0
	steps = 1
	i = 0

	// build the spiral (or just look it up here: https://oeis.org/A141481)
findsums:
	for {
		for range dist {
			// take a step
			x += hx
			y += hy

			// add all adjacent tiles
			for _, d := range dirs {
				board[y][x] += board[y+d[1]][x+d[0]]
			}

			steps++

			if advent.Viz.Enabled {
				fmt.Println(advent.BoardToString(board, func(val int) string {
					return fmt.Sprintf("%7d", val)
				}))
				fmt.Println()
			}

			if board[y][x] > num {
				s2 = board[y][x]
				break findsums
			}

		}

		// turn left
		hx, hy = -hy, hx
		if i%2 == 1 {
			dist++
		}
		i++
	}

	return s1, s2, nil
}
