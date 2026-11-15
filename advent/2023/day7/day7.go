package day7

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

type rank uint

const (
	rankA    rank = 0xE
	rankK    rank = 0xD
	rankQ    rank = 0xC
	rankJ    rank = 0xB
	rankT    rank = 0xA
	rank9    rank = 9
	rank8    rank = 8
	rank7    rank = 7
	rank6    rank = 6
	rank5    rank = 5
	rank4    rank = 4
	rank3    rank = 3
	rank2    rank = 2
	rankNone rank = 0
)

const (
	catHigh int = iota
	catPair
	cat2Pair
	cat3Kind
	catFullHouse
	cat4Kind
	cat5Kind
)

func (r rank) String() string {
	switch r {
	case rankA:
		return "A"
	case rankK:
		return "K"
	case rankQ:
		return "Q"
	case rankJ:
		return "J"
	case rankT:
		return "T"
	default:
		return strconv.Itoa(int(r))
	}
}

type hand struct {
	hand     [5]rank
	category int  // 0 = high, 1 = pair, 2 = 2-pair, 3 = three of a kind, 4 = full house, 5 = four of a kind, 6 = five of a kind
	score    uint // golang uses 64 bit ints on every system I care about
	bid      int
}

func (h hand) String() string {
	return h.hand[0].String() + h.hand[1].String() + h.hand[2].String() + h.hand[3].String() + h.hand[4].String()
}

func (h hand) scoreString() string {
	var sb strings.Builder

	r := rank((h.score >> (10 * 4)) & 0xf)
	if r > 0 {
		sb.WriteString(strings.Repeat(r.String(), 5))
	}
	r = rank((h.score >> (9 * 4)) & 0xf)
	if r > 0 {
		sb.WriteString(strings.Repeat(r.String(), 4))
	}
	r = rank((h.score >> (8 * 4)) & 0xf)
	if r > 0 {
		sb.WriteString(strings.Repeat(r.String(), 3))
	}
	r = rank((h.score >> (7 * 4)) & 0xf)
	if r > 0 {
		sb.WriteString(strings.Repeat(r.String(), 2))
	}
	r = rank((h.score >> (6 * 4)) & 0xf)
	if r > 0 {
		sb.WriteString(strings.Repeat(r.String(), 2))
	}
	for i := 5; i >= 0; i-- {
		r = rank((h.score >> (i * 4)) & 0xf)
		if r > 0 {
			sb.WriteString(r.String())
		}
	}

	return sb.String()
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1 int
	var s2 int

	// input is a card hand and a bid
	// 32T3K 765

	hands := make([]hand, 0)

	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte{'\n'}) {
		var cards [15]int
		var h [5]rank
		for i := range 5 {
			switch line[i] {
			case 'A':
				cards[rankA]++
				h[i] = rankA
			case 'K':
				cards[rankK]++
				h[i] = rankK
			case 'Q':
				cards[rankQ]++
				h[i] = rankQ
			case 'J':
				cards[rankJ]++
				h[i] = rankJ
			case 'T':
				cards[rankT]++
				h[i] = rankT
			case '9':
				cards[rank9]++
				h[i] = rank9
			case '8':
				cards[rank8]++
				h[i] = rank8
			case '7':
				cards[rank7]++
				h[i] = rank7
			case '6':
				cards[rank6]++
				h[i] = rank6
			case '5':
				cards[rank5]++
				h[i] = rank5
			case '4':
				cards[rank4]++
				h[i] = rank4
			case '3':
				cards[rank3]++
				h[i] = rank3
			case '2':
				cards[rank2]++
				h[i] = rank2
			}
		}

		bid, err := strconv.Atoi(string(line[6:]))
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse bid %s %v", string(line[6:]), err)
		}

		// score is a hex number based on the score of a hand, sorted from best to worst
		// 5 of a kind
		// 4 of a kind + high
		// 3 of a kind + 2 of a kind (full house)
		// 2 of a kind + 2 of a kind + high
		// 2 of a kind + 3 highs
		// 5 highs
		//
		// ranks are hex
		//   AKQJT98765432
		// 0xEDCBA98765432 // rank in hex
		//   543pphhhhh
		//
		// some card combos
		// 0xE000000000 = AAAAA (five aces)
		// 0x0E000D0000 = AAAAK (four aces, high king)
		// 0x00ED000000 = AAAKK (three aces, pair of kings)
		// 0x000ED00000 = AAKKQ (pair aces, pair of kings, high queen)
		// 0x0000EDCB00 = AAKQJ (pair aces, high king, queen, jack)
		// 0x00000EDCBA = AKQJT (high ace, high king, queen, jack, ten)

		// start with the highest value cards
		highs := 0
		pairs := 0
		score := uint(0x0)
		cat := catHigh
		for i := len(cards) - 1; i >= 0; i-- {
			switch cards[i] {
			case 5:
				score |= uint(i) << (10 * 4) // five of a kind is 10 hex digit
				cat = cat5Kind
			case 4:
				score |= uint(i) << (9 * 4) // four of a kind is 9th hex digit
				cat = cat4Kind
			case 3:
				score |= uint(i) << (8 * 4)
				if cat == catPair {
					cat = catFullHouse
				} else {
					cat = cat3Kind
				}

			case 2:
				if pairs == 0 {
					score |= uint(i) << (7 * 4)
					pairs++
					if cat == cat3Kind {
						cat = catFullHouse
					} else {
						cat = catPair
					}
				} else {
					score |= uint(i) << (6 * 4)
					cat = cat2Pair
				}
			case 1:
				score |= uint(i) << ((5 - highs) * 4)
				highs++
			}
		}

		// store hands as a score and a bid
		hands = append(hands, hand{hand: h, category: cat, score: score, bid: bid})
	}

	// p1: sort hands by category, then high card
	slices.SortFunc(hands, func(a, b hand) int {
		result := cmp.Compare(a.category, b.category)
		if result != 0 {
			return result
		}
		// sort by first card in each hand
		for i := range 5 {
			if a.hand[i] == b.hand[i] {
				continue
			}
			return cmp.Compare(a.hand[i], b.hand[i])
		}
		return 0
	})
	for i, h := range hands {
		s1 += h.bid * (i + 1)
		if advent.Viz.Enabled {
			fmt.Printf("%00d: %s - %d\n", i+1, h.String(), h.bid)
		}
	}

	return s1, s2, nil
}
