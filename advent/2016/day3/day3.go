package day3

import (
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct{}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	var nums [3][3]int
	var row, col int
	var num int
	for _, b := range input {
		digit := b - '0'
		if digit <= 9 {
			num = num*10 + int(digit)
			continue
		}
		if b == ' ' || b == '\n' {
			if num > 0 {
				nums[row][col] = num
				if advent.Viz.Enabled {
					fmt.Printf("r: %d, c: %d <- %d\n", row, col, num)
				}
				col++
			}
			num = 0
		}
		if b == '\n' {
			// check cols for solution 1
			num1, num2, num3 := nums[row][0], nums[row][1], nums[row][2]
			if checkTriangleSides(num1, num2, num3) {
				s1++
			}

			// count every three rows
			if row%3 == 2 {
				for r := range 3 {
					num1, num2, num3 := nums[0][r], nums[1][r], nums[2][r]
					if checkTriangleSides(num1, num2, num3) {
						s2++
					}
				}
			}

			row = (row + 1) % 3
			col = 0
		}
	}

	return s1, s2, nil
}

func checkTriangleSides(num1, num2, num3 int) bool {
	return num1+num2 > num3 &&
		num1+num3 > num2 &&
		num2+num3 > num1
}
