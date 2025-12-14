package day2

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for line := range strings.SplitSeq(string(bytes.TrimSpace(input)), "\n") {
		// input is like
		// 2x3x4
		dims := strings.Split(line, "x")
		l, err := strconv.Atoi(dims[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid input: %s %v", line, err)
		}
		w, err := strconv.Atoi(dims[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid input: %s %v", line, err)
		}
		h, err := strconv.Atoi(dims[2])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid input: %s %v", line, err)
		}

		s1 += wrappingPaperSqFt(l, w, h)
		s2 += ribbonLength(l, w, h)
	}

	return s1, s2, nil
}

func wrappingPaperSqFt(l, w, h int) int {
	return 2*l*w + 2*w*h + 2*h*l + min(l*w, w*h, l*h)
}

func ribbonLength(l, w, h int) int {
	vol := l * w * h
	shortestSide := min(l+w, l+h, w+h)
	return shortestSide*2 + vol
}
