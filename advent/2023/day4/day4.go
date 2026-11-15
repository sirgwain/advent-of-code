package day4

import (
	"bytes"
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

// printCards is the subset of cards to print when debugging
const printCards = 10

func (d *Day) Run(input []byte) (int, int, error) {
	var s1 int
	var s2 int

	// track bonus cards per card count. Card 1 has 4 matches so we end up with
	// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
	// [1 1 1 1 1 0 0 ...] in bonusCards
	var bonusCards [200]int

	// the input is fixed width, so store these ahead of time
	colon := bytes.Index(input, []byte(": "))
	pipe := bytes.Index(input, []byte(" | ")) - colon - 2

	cardNum := 0
	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte{'\n'}) {
		// a line looks like this
		// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
		// winning numbers are listed, then |, then your numbers
		line := line[colon+2:]

		var i int
		var winningNumbers [100]bool

		// this card counts itself plus any copies from before
		bonusCards[cardNum]++
		currentCardCount := bonusCards[cardNum]

		// find winning numbers
		for i = 0; i < pipe; i++ {
			num, numLen := advent.Atoi(line[i:])
			if num == 0 {
				// skip spaces
				continue
			}
			winningNumbers[num] = true
			i += numLen
		}

		matches := 0
		for ; i < len(line); i++ {
			num, numLen := advent.Atoi(line[i:])
			if num == 0 {
				// skip spaces
				continue
			}

			if winningNumbers[num] {
				matches++
			}

			i += numLen
		}

		if advent.Viz.Enabled {
			fmt.Printf("%v", bonusCards[:printCards])
		}

		// double the score per match
		s2 += currentCardCount
		if matches > 0 {
			s1 += 1 << (matches - 1)
			// each match gives us a bonus card in th enext draw
			for i := range matches {
				bonusCards[cardNum+1+i] += currentCardCount
			}
		}

		if advent.Viz.Enabled {
			fmt.Printf(" %d of card %d, %d matches -> %v\n", currentCardCount, cardNum+1, matches, bonusCards[:printCards])
		}
		cardNum++
	}

	return s1, s2, nil
}
