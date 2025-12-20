package day1

import "fmt"

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	s2 = -1
	for i, c := range input {
		switch c {
		case '(':
			s1++
		case ')':
			s1--
		default:
			return 0, 0, fmt.Errorf("unknown input %v", c)
		}
		if s1 == -1 && s2 == -1 {
			s2 = i
		}
	}
	s2 += 1

	return s1, s2, nil
}
