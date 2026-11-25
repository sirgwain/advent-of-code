//go:build !arm64 || purego

package day1

import (
	"bytes"

	"github.com/sirgwain/advent-of-code/advent"
)

func solve(input []byte) (s1, s2 int) {
	dial := 50
	start := dial
	for _, line := range bytes.Split(bytes.TrimSpace(input), []byte("\n")) {
		num, _ := advent.Atoi(line[1:])

		switch line[0] {
		case 'L':
			num = -num
		}

		start = dial
		dial += num

		// lop while dial < 0, adding 100 and checking how many times it passes 0
		for dial < 0 {
			dial += 100
			if start != 0 {
				s2++
			}
			start = dial
		}
		for dial >= 100 {
			dial = dial - 100
			if dial != 0 {
				s2++
			}
		}

		if dial == 0 {
			s1++
			s2++
		}

	}
	return s1, s2
}
