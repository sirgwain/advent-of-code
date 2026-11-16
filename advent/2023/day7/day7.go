package day7

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strconv"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

type rank uint

const (
	rankA     rank = 14
	rankK     rank = 13
	rankQ     rank = 12
	rankJ     rank = 11
	rankT     rank = 10
	rank9     rank = 9
	rank8     rank = 8
	rank7     rank = 7
	rank6     rank = 6
	rank5     rank = 5
	rank4     rank = 4
	rank3     rank = 3
	rank2     rank = 2
	rankJoker rank = 1
	rankNone  rank = 0
)

func (r rank) String() string {
	switch r {
	case rankA:
		return "A"
	case rankK:
		return "K"
	case rankQ:
		return "Q"
	case rankJ, rankJoker:
		return "J"
	case rankT:
		return "T"
	default:
		return strconv.Itoa(int(r))
	}
}

type cat int

const (
	catHigh cat = iota
	catPair
	cat2Pair
	cat3Kind
	catFullHouse
	cat4Kind
	cat5Kind
)

func (c cat) String() string {
	switch c {
	case catHigh:
		return "High Card"
	case catPair:
		return "Pair"
	case cat2Pair:
		return "Two Pairs"
	case cat3Kind:
		return "Three of a Kind"
	case catFullHouse:
		return "Full House"
	case cat4Kind:
		return "Four of a Kind"
	case cat5Kind:
		return "Five of a Kind"
	default:
		return "unknown"
	}
}

type hand struct {
	hand      [5]rank
	category  cat // 0 = high, 1 = pair, 2 = 2-pair, 3 = three of a kind, 4 = full house, 5 = four of a kind, 6 = five of a kind
	category2 cat // category for part2
	bid       int
}

func (h hand) String() string {
	return h.hand[0].String() + h.hand[1].String() + h.hand[2].String() + h.hand[3].String() + h.hand[4].String()
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1 int
	var s2 int

	hands := make([]hand, 0, bytes.Count(input, []byte{'\n'})+1)

	// build a mapping of card character to rank
	var cardVal = func() (t [256]rank) {
		for i, c := range []byte("23456789TJQKA") {
			t[c] = rank(i + 2) // 2..14, so J = 11
		}
		return
	}()

	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte{'\n'}) {

		// input is a card hand and a bid
		// 32T3K 765

		var cards [15]int
		var h [5]rank
		for i := range 5 {
			r := cardVal[line[i]]
			cards[r]++ // add this card to the cards count
			h[i] = r   // remember this card for the hand
		}

		// bid is the number after the 5 cards
		bid, err := strconv.Atoi(string(line[6:]))
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse bid %s %v", string(line[6:]), err)
		}

		// start with the highest value cards
		cat1 := catHigh
		cat2 := catHigh
		for i := len(cards) - 1; i >= 0; i-- {
			// part 1, Js are jacks
			cat1 = catForCardCount(cards[i], cat1)

			// part 2, Js are jokers so skip them when
			// checking initial category
			if i == int(rankJ) {
				continue
			}
			cat2 = catForCardCount(cards[i], cat2)
		}

		// category gets better per joker, joker is low rank
		cat2 = catUpgradeForJokers(cards[rankJ], cat2)

		// store hands as a score and a bid
		hands = append(hands, hand{hand: h, category: cat1, category2: cat2, bid: bid})
	}

	// p1: sort hands by category, then high card
	slices.SortFunc(hands, func(a, b hand) int {
		result := cmp.Compare(a.category, b.category)
		if result != 0 {
			return result
		}

		// sort by cards in each hand in order, first higher card wins
		for i := range 5 {
			if a.hand[i] == b.hand[i] {
				continue
			}
			return cmp.Compare(a.hand[i], b.hand[i])
		}
		return 0
	})

	if advent.Viz.Enabled {
		fmt.Println()
		fmt.Println("part1")
	}

	// sum up part 1 bids * hand rank
	for i, h := range hands {
		s1 += h.bid * (i + 1)
		if advent.Viz.Enabled {
			fmt.Printf("%00d: %s - %d (%s)\n", i+1, h.String(), h.bid, h.category)
		}
	}

	// p2: sort hands by category, then high card
	slices.SortFunc(hands, func(a, b hand) int {
		result := cmp.Compare(a.category2, b.category2)
		if result != 0 {
			return result
		}

		// sort by cards in each hand in order, first higher card wins
		for i := range 5 {
			// jokers are lowest rank when sorting
			ah := a.hand[i]
			bh := b.hand[i]

			if ah == bh {
				continue
			}

			if ah == rankJ {
				ah = rankJoker
			}
			if bh == rankJ {
				bh = rankJoker
			}
			return cmp.Compare(ah, bh)
		}
		return 0
	})

	if advent.Viz.Enabled {
		fmt.Println()
		fmt.Println("part2")
	}

	// sum up part 2 bids * hand rank
	for i, h := range hands {
		s2 += h.bid * (i + 1)
		if advent.Viz.Enabled {
			fmt.Printf("%00d: %s - %d (%s)\n", i+1, h.String(), h.bid, h.category2)
		}
	}

	return s1, s2, nil
}

// catForCardCount returns a new category for a number of cards, how many pairs
// we've already counted in this hand, and the current category (for full house comparison)
func catForCardCount(numCards int, currentCat cat) cat {
	switch numCards {
	case 5:
		return cat5Kind
	case 4:
		return cat4Kind
	case 3:
		switch currentCat {
		case catPair:
			// we have a pair, now we have a full house
			return catFullHouse
		}
		return cat3Kind
	case 2:
		switch currentCat {
		case catPair:
			// we have a pair, now we have 2
			return cat2Pair
		case cat3Kind:
			// we have a 3 of a kind, now we have a ful house
			return catFullHouse
		}
		// we formed a pair
		return catPair
	}

	return currentCat
}

// catUpgradeForJokers returns an upgrade category for a number of jokers
// based on the current category
func catUpgradeForJokers(jokers int, currentCat cat) cat {
	switch jokers {
	case 4, 5: // 4 or 5 jacks is always 5 of a kind
		return cat5Kind
	case 3:
		if currentCat == catPair {
			return cat5Kind
		} else {
			return cat4Kind
		}
	case 2:
		switch currentCat {
		case catPair:
			// 2 jokers and a pair gives 4 of a kind
			return cat4Kind
		case cat3Kind:
			// 2 jokers and a 3 of a kind gives 5 of a kind
			return cat5Kind
		default:
			// 2 jokes and a single gives 3 of a kind
			return cat3Kind
		}
	case 1:
		switch currentCat {
		case cat4Kind:
			return cat5Kind
		case cat3Kind:
			return cat4Kind
		case cat2Pair:
			return catFullHouse
		case catPair:
			return cat3Kind
		default:
			return catPair
		}
	}
	return currentCat
}
