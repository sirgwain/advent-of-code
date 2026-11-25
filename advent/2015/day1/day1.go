package day1

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	s1, s2 = solve(input)

	return s1, s2, nil
}
