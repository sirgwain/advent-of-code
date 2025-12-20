package day13

import (
	"bufio"
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/sirgwain/advent-of-code/advent"
)

const me = "sirgwain"

type Day struct {
}

type pair struct {
	a string
	b string
}

func pairKey(a, b string) pair {
	if a < b {
		return pair{a, b}
	}
	return pair{b, a}
}

func (p pair) String() string {
	return fmt.Sprintf("%s - %s", p.a, p.b)
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	happiness := make(map[pair]int)
	people := make(map[string]bool, 0)

	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input)))
	for scanner.Scan() {
		line := scanner.Text()
		// line is like
		// Alice would gain 54 happiness units by sitting next to Bob.
		// Alice would lose 79 happiness units by sitting next to Carol.
		//   0     1     2   3     4       5    6   7      8    9  n-1
		fields := strings.Fields(line)
		source := fields[0]
		dest := strings.TrimSuffix(fields[len(fields)-1], ".")
		gainLose := fields[2]
		weightStr := fields[3]

		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}
		if gainLose == "lose" {
			weight = -weight
		}

		people[source] = true
		people[dest] = true

		key := pairKey(source, dest)
		happiness[key] += weight
	}

	names := slices.Sorted(maps.Keys(people))
	root := names[0]

	// iterate over indices, not strings, so we can quickly compare in our permutation
	// visit function and skip rings that are reverse order of another ring
	// also, skip the first node because we are looking for a ring,
	// so if we have "Alice - Bob - Carol - Alice" we don't also do "Bob - Carol - Alice - Bob" because that's the same ring
	rest := make([]string, len(names)-1)
	for i := range len(names) - 1 {
		rest[i] = names[i+1]
	}

	// find the best happiness for each combination
	advent.HeapPermute(rest, len(rest), func(perm []string) {
		s1 = max(s1, computeHappiness(root, happiness, perm))
	})

	// add me into the mix
	for _, name := range names {
		happiness[pairKey(me, name)] = 0
	}

	// do the permutation again with all the names, not just the names minus the first name
	// I'm the new root
	advent.HeapPermute(names, len(names), func(perm []string) {
		s2 = max(s2, computeHappiness(me, happiness, perm))
	})

	return s1, s2, nil
}

// computeHappiness updates the best int pointer with the happiness from this permutation, if it's better
func computeHappiness(root string, happiness map[pair]int, perm []string) int {
	if perm[0] > perm[len(perm)-1] {
		return 0 // skip the reverse
	}

	// start weight from root -> first name
	weight := happiness[pairKey(root, perm[0])]

	// build weights for 0 -> 1, 1 -> 2, 2 -> n
	for i := 0; i < len(perm)-1; i++ {
		a := perm[i]
		b := perm[i+1]
		weight += happiness[pairKey(a, b)]
	}
	// add in a final weight for n -> root
	// this completes the ring
	weight += happiness[pairKey(perm[len(perm)-1], root)]

	// we are looking to optimize happiness
	return weight
}
