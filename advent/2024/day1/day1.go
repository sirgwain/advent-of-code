package day1

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	// every line is two same-width numbers separated by 3 spaces, plus '\n':
	// 37033   48086
	// so line i starts at i*stride and we never have to search for newlines.
	numWidth := bytes.IndexByte(input, ' ')
	stride := 2*numWidth + 4
	n := (len(input) + 1) / stride // +1 in case the last line has no '\n'

	// one allocation call shared between two number lists
	buf := make([]int, 2*n)
	nums1, nums2 := buf[:n], buf[n:]
	for i := range n {
		line := input[i*stride:]
		nums1[i] = advent.AtoiFixed(line[:numWidth])
		nums2[i] = advent.AtoiFixed(line[numWidth+3 : 2*numWidth+3])
	}

	// part 1, sort the lists and add up the distance
	slices.Sort(nums1)
	slices.Sort(nums2)
	for i := range len(nums1) {
		s1 += advent.Abs(nums1[i] - nums2[i])
	}

	// part 2: both lists are sorted, so walk them together. When the values
	// match, count the run of that value in each list and add v * left * right.
	for i, j := 0, 0; i < len(nums1) && j < len(nums2); {
		v := nums1[i]
		switch {
		case v < nums2[j]:
			i++
		case v > nums2[j]:
			j++
		default:
			ci, cj := i, j
			for i < len(nums1) && nums1[i] == v {
				i++
			}
			for j < len(nums2) && nums2[j] == v {
				j++
			}
			s2 += v * (i - ci) * (j - cj)
		}
	}

	if advent.Viz.Enabled {
		for i := range len(nums1) {
			fmt.Printf("%d   %d\n", nums1[i], nums2[i])
		}
	}

	return s1, s2, nil
}
