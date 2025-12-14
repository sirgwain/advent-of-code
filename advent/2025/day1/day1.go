package day1

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

	dial := 50
	start := dial
	for _, line := range strings.Split(string(bytes.TrimSpace(input)), "\n") {
		num, err := strconv.Atoi(line[1:])
		if err != nil {
			return 0, 0, fmt.Errorf("error parsing number on line: %s", line)
		}

		switch line[0] {
		case 'L':
			num = -num
		}

		start = dial
		dial += num

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
	return s1, s2, nil
}
