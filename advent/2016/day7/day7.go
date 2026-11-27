package day7

import (
	"bytes"
	"iter"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct{}

const abaInSupernet = 1
const abaInHypernet = 2

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	// encoded ABA pairs using 5bit letter packing
	var abaPairs [1024]byte

	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte("\n")) {
		clear(abaPairs[:])
		supernetAbba := false
		hypernetAbba := false
		supernetAba := false

		for seg, inhypernet := range segments(line) {
			// part1, search for ABBAs
			if !hypernetAbba && findAbba(seg) {
				if inhypernet {
					hypernetAbba = true
				} else {
					supernetAbba = true
				}
			}
			// part2, search for ABAs with matching BABs in hypernets
			if !supernetAba {
				for i := 0; i < len(seg)-2; i++ {

					// check for ABA
					a1, b1, a2 := seg[i], seg[i+1], seg[i+2]
					if a1 == a2 && a1 != b1 {
						var index uint16
						var val byte
						if inhypernet {
							// reverse the check in the hypernet
							index = advent.PackLetters(b1, a1)
							val = abaInHypernet
						} else {
							index = advent.PackLetters(a1, b1)
							val = abaInSupernet
						}
						if inhypernet && abaPairs[index] == abaInSupernet {
							// we have a pair from a segment and we found a match in the hypernet
							supernetAba = true
							break
						} else if !inhypernet && abaPairs[index] == abaInHypernet {
							// we have a pair from a hypernet and we found a match in the segment
							supernetAba = true
							break
						}
						abaPairs[index] = val
					}
				}
			}
		}

		if !hypernetAbba && supernetAbba {
			s1++
		}
		if supernetAba {
			s2++
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
