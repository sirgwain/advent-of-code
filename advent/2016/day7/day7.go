package day7

import (
	"bytes"
	"iter"
)

type Day struct{}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte("\n")) {
		hasAbba := false
		for seg, inhypernet := range segments(line) {
			ok := findAbba(seg)
			if ok && inhypernet {
				ok = false
				break
			}
			hasAbba = hasAbba || ok
		}
		if hasAbba {
			s1++
		}
	}

	return s1, s2, nil
}

// segments returns each segment of an ip7 address
// and whether it's in a hypernet block
// ip7 addresses are like abba[mnop]qrst
func segments(line []byte) iter.Seq2[[]byte, bool] {
	return func(yield func([]byte, bool) bool) {
		// iterate through the line, returning each segment
		// and toggling hypernet on or off as we hit braces
		inhypernet := false
		for len(line) > 0 {
			sep := byte('[')
			if inhypernet {
				sep = byte(']')
			}
			i := bytes.IndexByte(line, sep)
			if i < 0 {
				// no separator found, return the rest of the line
				yield(line, inhypernet)
				return
			}
			// return up to the separator, and if the caller
			// breaks out, return
			if !yield(line[:i], inhypernet) {
				return
			}

			// keep processing the line
			line = line[i+1:]
			inhypernet = !inhypernet
		}
	}
}

// findAbba returns true if an abba or xyyx style
// input pattern is in the input
func findAbba(input []byte) bool {
	for i := range input {
		if i+3 >= len(input) {
			break
		}
		// look for ABBA patterns
		// where a1 == b1 && a2 == b2
		// a1 b1   b2 a2
		// A  B  = B  A
		a1 := input[i]
		b1 := input[i+1]
		b2 := input[i+2]
		a2 := input[i+3]

		if a1 == a2 && b1 == b2 && a1 != b1 {
			return true
		}
	}
	return false
}
