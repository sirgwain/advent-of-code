package day2

import (
	"bytes"
	"math"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	var nums [16]int
	input = bytes.TrimSpace(input)
	for line := range bytes.Lines(input) {
		largest := 0
		smallest := math.MaxInt
		col := 0
		for num := range advent.NumberFieldsSeq(line) {
			largest = max(largest, num)
			smallest = min(smallest, num)
			nums[col] = num
			col++
		}

		s1 += largest - smallest

		// part2, find the whole number divide in this row
	dividecheck:
		for i := range col {
			for j := i + 1; j < col; j++ {
				a := nums[i]
				b := nums[j]
				if a < b {
					// make sure we check large / small
					a, b = b, a
				}
				if a%b == 0 {
					s2 += a / b
					break dividecheck
				}
			}
		}
	}
	return s1, s2, nil
}
