package day4

import (
	"bytes"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for line := range bytes.SplitSeq(input, []byte{'\n'}) {
		// a line looks like this
		// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
		// winning numbers are listed, then |, then your numbers
		_, line, _ = bytes.Cut(line, []byte{':', ' '})

		var i int
		var winningNumbers [100]bool

		// find winning numbers
		for i = 0; i < len(line); i++ {
			if line[i] == '|' {
				break
			}

			num, numLen := advent.Atoi(line[i:])
			if num == 0 {
				// skip spaces
				continue
			}
			winningNumbers[num] = true
			i += numLen
		}

		matches := 0
		for ; i < len(line); i++ {
			num, numLen := advent.Atoi(line[i:])
			if num == 0 {
				// skip spaces
				continue
			}

			if winningNumbers[num] {
				matches++
			}

			i += numLen
		}

		// double the score per match
		if matches > 0 {
			s1 += 1 << (matches - 1)
		}
	}

	return s1, s2, nil
}
