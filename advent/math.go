package advent

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
