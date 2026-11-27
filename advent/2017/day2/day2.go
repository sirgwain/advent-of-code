package day2

import (
	"bytes"
	"math"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	var nums [16]int
	input = bytes.TrimSpace(input)
	for line := range bytes.Lines(input) {
		num := 0
		largest := 0
		smallest := math.MaxInt
		col := 0
		for i := range line {
			digit := line[i] - '0'
			if digit > 9 {
				largest = max(largest, num)
				smallest = min(smallest, num)
				nums[col] = num
				col++
				num = 0
				continue
			}
			num = num*10 + int(digit)
		}
		if num != 0 {
			nums[col] = num
			largest = max(largest, num)
			smallest = min(smallest, num)
		}

		s1 += largest - smallest

		// part2, find the whole number divide in this row
		for i := range col {
			for j := i + 1; j < col; j++ {
				if nums[i]%nums[j] == 0 {
					s2 += nums[i] / nums[j]
					break
				} else if nums[j]%nums[i] == 0 {
					s2 += nums[j] / nums[i]
					break
				}
			}
		}
	}
	return s1, s2, nil
}
