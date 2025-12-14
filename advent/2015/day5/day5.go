package day5

import (
	"bytes"
	"slices"
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

	if slices.Contains(vowels, b[len(b)-1]) {
		numVowels++
	}

	hasRepeat := false
	for i := 0; i < len(s)-1; i++ {
		twoChars := string(b[i : i+2])
		switch twoChars {
		case "ab", "cd", "pq", "xy":
			return false
		}
		if slices.Contains(vowels, b[i]) {
			numVowels++
		}
		hasRepeat = hasRepeat || b[i] == b[i+1]
		if numVowels >= 3 && hasRepeat {
			return true
		}
	}

	return false
}

func isNice2(s string) bool {
	b := []byte(s)

	pairIndexes := make(map[string]int)
	hasSurroundingPair := false
	hasDoublePairs := false
	for i := 0; i < len(s)-1; i++ {
		twoChars := string(b[i : i+2])
		if j, ok := pairIndexes[twoChars]; ok {
			if i-j >= 2 {
				hasDoublePairs = true
			}
		} else {
			pairIndexes[twoChars] = i
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
