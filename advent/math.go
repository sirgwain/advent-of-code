package advent

import "golang.org/x/exp/constraints"

// Greatest common denominator
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Least common multiple
func LCM(a, b int) int {
	return a / GCD(a, b) * b
}

func Abs[T constraints.Integer | constraints.Float](num T) T {
	if num < 0 {
		return -num
	}
	return num
}

func Clamp(n, mn, mx int) int {
	return max(mn, min(mx, n))
}