package day5

import (
	"bytes"
	"strings"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for _, line := range strings.Split(string(bytes.TrimSpace(input)), "\n") {
		if isNice1(line) {
			s1++
		}
		if isNice2(line) {
			s2++
		}
	}

	return s1, s2, nil
}

func isNice1(s string) bool {
	b := []byte(s)
	vowels := []byte("aeiou")
	numVowels := 0

	if bytes.IndexByte(vowels, b[len(b)-1]) != -1 {
		numVowels++
	}

	hasRepeat := false
	for i := 0; i < len(s)-1; i++ {
		twoChars := string(b[i : i+2])
		switch twoChars {
		case "ab", "cd", "pq", "xy":
			return false
		}
		if bytes.IndexByte(vowels, b[i]) != -1 {
			numVowels++
		}
		hasRepeat = hasRepeat || b[i] == b[i+1]
	}

	if numVowels >= 3 && hasRepeat {
		return true
	}

	return false
}

func isNice2(s string) bool {
	b := []byte(s)

	pairIndexes := make(map[uint16]int)
	hasSurroundingPair := false
	hasDoublePairs := false
	for i := 0; i < len(s)-1; i++ {
		// build a two byte number from the current character and the next one
		pair := uint16(b[i])<<8 | uint16(b[i+1])
		if j, ok := pairIndexes[pair]; ok {
			if i-j >= 2 {
				hasDoublePairs = true
			}
		} else {
			pairIndexes[pair] = i
		}
		if i+2 < len(s) {
			hasSurroundingPair = hasSurroundingPair || b[i] == b[i+2]
		}
		if hasDoublePairs && hasSurroundingPair {
			return true
		}
	}

	return false
}
