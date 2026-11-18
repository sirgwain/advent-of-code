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

	n := bytes.Count(input, []byte{'\n'}) + 1
	nums1, nums2 := make([]int, 0, n), make([]int, 0, n)
	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte{'\n'}) {
		// lines are numbers separated by 3 spaces
		// 37033   48086
		n1, len := advent.Atoi(line)
		n2, _ := advent.Atoi(line[len+3:])
		nums1 = append(nums1, n1)
		nums2 = append(nums2, n2)
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
