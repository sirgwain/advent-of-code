//go:build arm64 && !purego
package day1

//go:noescape
func solve(input []byte) (floor, basement int)