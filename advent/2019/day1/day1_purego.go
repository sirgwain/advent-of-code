//go:build !arm64 || purego

package day1

import (
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

func solve(input []byte) (s1, s2 int) {
	var num int

	for _, b := range input {
		digit := int(b-'0')

		if digit > 9 {
			s1 += (num / 3) - 2
			if advent.Viz.Enabled {
				fmt.Printf("%d / 3 - 2 = %d\n", num, (num/3)-2)
			}
			num = 0
		} else {
			num = 10*num + digit
		}
	}
	return s1, s2
}
