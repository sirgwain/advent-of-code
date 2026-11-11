package day2

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	// s1: only count games that have 12 red, 13 green, 14 blue or greater
	max_rbg := []int{12, 13, 14}

	// lines are games with ; separated pulls, comma separated colors
	// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
	// Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
	game := 0
	for line := range strings.SplitSeq(string(bytes.TrimSpace(input)), "\n") {
		// skip "Game n: "
		_, line, _ = strings.Cut(line, ": ")
		game++

		// count the min cubes that make this game possible
		minCubes := []int{0, 0, 0}

		valid := true

		// split into ["3", "blue,", "4", "red;", "1", "red,", "2", "green,", "6", "blue;", "2", "green"]
		fields := strings.Fields(line)
		for i := 0; i < len(fields); i += 2 {

			num, err := strconv.Atoi(fields[i])
			if err != nil {
				return 0, 0, fmt.Errorf("failed to convert %s from line %s to int %v", fields[i], line, err)
			}

			// for the color, check the first letter
			colorIndex := 0
			switch fields[i+1][0] {
			case 'r':
				colorIndex = 0
			case 'g':
				colorIndex = 1
			case 'b':
				colorIndex = 2
			}

			// if this number is greater than our max amount (i.e. 13 red when 12 is the max), it's an invalid game
			if num > max_rbg[colorIndex] {
				valid = false
			}

			// find the minimum num ber of cubes that would make this game valid from all the draws
			minCubes[colorIndex] = max(num, minCubes[colorIndex])
		}
		if valid {
			s1 += game
		}

		s2 += minCubes[0] * minCubes[1] * minCubes[2]
	}

	return s1, s2, nil
}
