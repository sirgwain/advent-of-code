//go:build !arm64 || purego

package day1

import "bytes"

func solve(input []byte) (s1, s2 int) {
	input = bytes.TrimSpace(input)

	var cur, prev byte

	// part2 iterate half the list first
	mid := len(input) / 2
	for i := range input[mid:] {
		cur = input[i]
		prev = input[mid]

		// crazy. You have to do this tmp assignment to get the go compiler to use CSEL
		tmp := int(cur - '0')
		if cur == prev {
			tmp = 0
		}
		s2 += tmp
		mid++
	}
	s2 *= 2

	// part1, check each number with the next number
	for i := 1; i < len(input); i++ {
		prev = input[i-1]
		cur = input[i]
		tmp := int(cur - '0')
		if prev == cur {
			tmp = 0
		}
		s1 += tmp
	}
	if cur == input[0] {
		s1 += int(cur - '0')
	}

	return s1, s2
}
