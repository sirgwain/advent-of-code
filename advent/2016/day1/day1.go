package day1

type Day struct {
}

// keep a grid of 512x512 bits to track directions
var grid [512 * 512 / 64]uint64

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	s1, s2 = solve(input)

	return s1, s2, nil
}
