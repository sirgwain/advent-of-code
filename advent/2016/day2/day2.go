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

// RunX returns both codes as packed hex digits, one nibble per key. Part 1's
// keys are all 1-9, so its hex form reads the same as the decimal code.
func (d *Day) RunX(input []byte) (uint, uint, error) {
	s1, s2 := solve(input)
	return s1, s2, nil
}

func solve(input []byte) (s1, s2 uint) {

	var keypad = [5][5]byte{
		{0, 0, 0, 0, 0},
		{0, 1, 2, 3, 0},
		{0, 4, 5, 6, 0},
		{0, 7, 8, 9, 0},
		{0, 0, 0, 0, 0},
	}

	//     1
	//   2 3 4
	// 5 6 7 8 9
	//   A B C
	//     D
	var keypadDiag = [7][7]byte{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 2, 3, 4, 0, 0},
		{0, 5, 6, 7, 8, 9, 0},
		{0, 0, 0xA, 0xB, 0xC, 0, 0},
		{0, 0, 0, 0xD, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	}

	// the current location on the keypad, starts at 5
	x1, y1 := 2, 2
	x2, y2 := 1, 3
	for i := range input {
		b := input[i]

		if b == '\n' {
			// record number
			s1 = s1<<4 | uint(keypad[y1][x1])
			s2 = s2<<4 | uint(keypadDiag[y2][x2])
			if advent.Viz.Enabled {
				fmt.Printf("Record (%d,%d): %X\n", x1, y1, keypad[y1][x1])
				fmt.Printf("Record (%d,%d): %X\n", x2, y2, keypadDiag[y2][x2])
			}

			continue
		}
		dir := dirs[b]

		if keypad[y1+dir[1]][x1+dir[0]] != 0 {
			x1, y1 = x1+dir[0], y1+dir[1]
		}
		if keypadDiag[y2+dir[1]][x2+dir[0]] != 0 {
			x2, y2 = x2+dir[0], y2+dir[1]
		}
		if advent.Viz.Enabled {
			fmt.Printf("%v dir: %d,%d -> (%d,%d)\n", string(b), dir[0], dir[1], x1, y1)
			fmt.Printf("%v dir: %d,%d -> (%d,%d)\n", string(b), dir[0], dir[1], x2, y2)
		}
	}

	return s1, s2
}
