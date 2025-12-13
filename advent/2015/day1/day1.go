package day1

import "fmt"

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {

	floor := 0
	firstBasementIndex := -1
	for i, c := range input {
		switch c {
		case '(':
			floor++
		case ')':
			floor--
		default:
			return 0, 0, fmt.Errorf("unknown input %v", c)
		}
		if floor == -1 && firstBasementIndex == -1 {
			firstBasementIndex = i
		}
	}
	return floor, firstBasementIndex + 1, nil
}
