package day1

import (
	"strings"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	return runPart(input, false), runPart(input, true), nil
}

func runPart(input []byte, countWords bool) int {
	var s2 int

	first := -1
	last := -1

	number_strings := []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	}

	for i, c := range input {
		digit := -1
		if advent.IsDigit(rune(c)) {
			digit = advent.ToDigit(c)
		} else if countWords {
			// word strings count as digits as well
			for j := range len(number_strings) {
				if strings.HasPrefix(string(input[i:]), number_strings[j]) {
					digit = j + 1
				}
			}
		}

		// found a number word or digit char
		if digit != -1 {
			if first == -1 {
				first = digit
			}
			last = digit
		}

		// at end of line or end of input, count
		// the digits as one num, i.e. 1 2 becomes 12
		if c == '\n' || i == len(input) {
			if first == -1 {
				// no digits this line, skip
				continue
			}

			result := first*10 + last
			s2 += result

			// reset digits for new line
			first = -1
			last = -1
		}
	}

	return s2
}
