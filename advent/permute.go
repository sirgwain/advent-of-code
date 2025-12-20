package advent

// HeapPermute uses heap's permutation algorithm to generate every combination of a slice
// and call a visit function for each combo
func HeapPermute[T any](slice []T, n int, visit func(perm []T)) {
	if n == 1 {
		visit(slice)
		return
	}
	for i := range n {
		HeapPermute(slice, n-1, visit)
		if n%2 == 1 {
			slice[0], slice[n-1] = slice[n-1], slice[0]
		} else {
			slice[i], slice[n-1] = slice[n-1], slice[i]
		}
	}
}
