package day2

import (
	"bytes"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte{'\n'}) {

		var nums [16]int
		count := 0
		for field := range bytes.FieldsSeq(line) {
			num, _ := advent.Atoi(field)
			nums[count] = num
			count++
		}

		s1 += part1(nums, count)
		s2 += part2(nums, count)
	}

	return s1, s2, nil
}

func check(num, prev, prevDiff int) (int, bool) {
	diff := prev - num
	if advent.Abs(diff) > 3 || diff == 0 {
		// unsafe because diff is > 3 or not changing
		return 0, false
	}

	if prevDiff == 0 {
		prevDiff = diff
	} else if (prevDiff ^ diff) < 0 {
		// different signs, unsafe
		return 0, false
	}

	return diff, true
}

func checkLevels(nums [16]int, count int, skip int) int {
	prevDiff := 0
	prev := -1
	for i := range count {
		if skip == i {
			continue
		}
		if prev == -1 {
			prev = nums[i]
			continue
		}
		num := nums[i]

		var ok bool
		if prevDiff, ok = check(num, prev, prevDiff); !ok {
			return i
		}
		prev = num
	}

	return -1
}

// part1 checks levels in a list of numbers
func part1(nums [16]int, count int) int {
	if checkLevels(nums, count, -1) == -1 {
		return 1
	}
	return 0
}

// part1 checks levels in a range, but drops an unhealthy level
func part2(nums [16]int, count int) int {
	badLevel := checkLevels(nums, count, -1)
	if badLevel == -1 {
		return 1
	}

	// skip the bad level
	// (badLevel-1, badLevel) is bad, so try the pair before that and after it
	if checkLevels(nums, count, badLevel) == -1 {
		return 1
	}
	// skip the level before the bad level
	if badLevel > 0 && checkLevels(nums, count, badLevel-1) == -1 {
		return 1
	}
	// skip the level 2 levels before if it's the start
	if badLevel == 2 && checkLevels(nums, count, badLevel-2) == -1 {
		return 1
	}

	return 0

}
