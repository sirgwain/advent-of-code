//go:build !arm64 || purego

package day1

import "bytes"

func solve(input []byte) (s1, s2 int) {

	start := input[0] - '0'
	var num2 byte
	for i := 1; i < len(bytes.TrimSpace(input)); i++ {
		num1 := input[i-1] - '0'
		num2 = input[i] - '0'
		if num1 == num2 {
			s1 += int(num2)
		}
	}
	if num2 == start {
		s1 += int(num2)
	}

	return s1, s2
}
