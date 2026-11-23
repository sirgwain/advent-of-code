package day2

import (
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

var dirs = [256][2]int{
	'U': {0, -1},
	'D': {0, 1},
	'L': {-1, 0},
	'R': {1, 0},
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	s1 = part1(input)
	s2 = part2(input)

	return s1, s2, nil
}

func part1(input []byte) (s1 int) {
	var keypad = [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	// the current location on the keypad, 5
	// 123
	// 456
	// 789
	x, y := 1, 1
	for i := range input {
		b := input[i]

		if b == '\n' {
			// record number
			s1 = 10*s1 + keypad[y][x]
			if advent.Viz.Enabled {
				fmt.Printf("Record (%d,%d): %d\n", x, y, keypad[y][x])
			}

			continue
		}
		dir := dirs[b]
		x = advent.Clamp(x+dir[0], 0, 2)
		y = advent.Clamp(y+dir[1], 0, 2)
		if advent.Viz.Enabled {
			fmt.Printf("%v dir: %d,%d -> (%d,%d)\n", string(b), dir[0], dir[1], x, y)
		}
	}

	return s1
}

func part2(input []byte) (s2 int) {

	//     1
	//   2 3 4
	// 5 6 7 8 9
	//   A B C
	//     D
	var keypad = [][]byte{
		{0, 0, 1, 0, 0},
		{0, 2, 3, 4, 0},
		{5, 6, 7, 8, 9},
		{0, 0xA, 0xB, 0xC, 0},
		{0, 0, 0xD, 0, 0},
	}

	// the current location on the keypad, starts at 5
	x, y := 0, 2
	for i := range input {
		b := input[i]

		if b == '\n' {
			// record number
			s2 = s2<<4 + int(keypad[y][x])
			if advent.Viz.Enabled {
				fmt.Printf("Record (%d,%d): %X\n", x, y, keypad[y][x])
			}

			continue
		}
		dir := dirs[b]

		x1 := advent.Clamp(x+dir[0], 0, len(keypad[y])-1)
		y1 := advent.Clamp(y+dir[1], 0, len(keypad)-1)
		if keypad[y1][x1] != 0 {
			x, y = x1, y1
		}
		if advent.Viz.Enabled {
			fmt.Printf("%v dir: %d,%d -> (%d,%d)\n", string(b), dir[0], dir[1], x, y)
		}
	}

	return s2
}
