package day6

import (
	"bytes"
	"math"
)

type Day struct{}

const maxCols = 8

func (d *Day) RunS(input []byte) (string, string, error) {
	var s1, s2 string

	// find the length of each line, thats our column count
	n := bytes.IndexByte(input, '\n')
	var counts [maxCols][26]uint16

	// count up the letters in each column
	stride := n + 1
	for row := 0; row+n < len(input); row += stride {
		for col, b := range input[row : row+n] {
			counts[col][b-'a']++
		}
	}

	// part 1, find the max letter count and associated letter for each column
	var maxCounts [maxCols]uint16
	var maxLetters [maxCols]byte

	// part 2, find the min letter in each column
	var minCounts [maxCols]uint16
	var minLetters [maxCols]byte
	for i := range minCounts {
		minCounts[i] = math.MaxUint16
	}

	// go through each column, then each letter count for that column and get the min
	for col, clc := range counts {
		for li, c := range clc {
			if c == 0 {
				// skip empties
				continue
			}

			if c > maxCounts[col] {
				maxCounts[col] = c
				maxLetters[col] = byte(li + 'a')
			}

			if c < minCounts[col] {
				minCounts[col] = c
				minLetters[col] = byte(li + 'a')
			}
		}
	}

	s1 = string(maxLetters[:n])
	s2 = string(minLetters[:n])

	return s1, s2, nil
}
