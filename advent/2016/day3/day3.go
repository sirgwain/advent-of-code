package day3

import (
	"bytes"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct{}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte("\n")) {
		fields := bytes.Fields(line)
		num1, _ := advent.Atoi(fields[0])
		num2, _ := advent.Atoi(fields[1])
		num3, _ := advent.Atoi(fields[2])
		if num1+num2 > num3 &&
			num1+num3 > num2 &&
			num2+num3 > num1 {
			s1++
		}
	}

	return s1, s2, nil
}
