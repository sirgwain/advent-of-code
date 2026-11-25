//go:build !arm64 || purego

package day1

func solve(input []byte) (floor, basement int) {
	basement = -1
	for i, c := range input {
		switch c {
		case '(':
			floor++
		case ')':
			floor--
		}
		if basement == -1 && floor == -1 {
			basement = i
		}
	}
	basement += 1

	return floor, basement
}
