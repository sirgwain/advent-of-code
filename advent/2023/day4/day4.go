package day4

import (
	"bytes"
	"strings"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for line := range strings.SplitSeq(string(bytes.TrimSpace(input)), "\n") {
		// a line looks like this
		// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
		// winning numbers are listed, then |, then your numbers
		_, line, _ = strings.Cut(line, ": ")

		var i int
		var winningNumbers [100]bool
		lineBytes := []byte(line)

		// find winning numbers
		for i = 0; i < len(lineBytes); i++ {
			if lineBytes[i] == '|' {
				break
			}

			num, numLen := advent.Atoi(lineBytes[i:])
			if num == 0 {
				// skip spaces
				continue
			}
			winningNumbers[num] = true
			i += numLen
		}

		matches := 0
		for ; i < len(lineBytes); i++ {
			num, numLen := advent.Atoi(lineBytes[i:])
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
