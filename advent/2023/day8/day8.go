package day8

import (
	"bytes"
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

// a key is 3 bytes packed
// 0x414243 = ABC
type key uint

// a val is 6 bytes packed
// 0x414243_444547 = ABC_DEF
type val uint

func keyFromBytes(b []byte) key {
	return key((uint32(b[2]) << 16) | (uint32(b[1]) << 8) | uint32(b[0]))
}

func (k key) String() string {
	return fmt.Sprintf("%c%c%c", (k>>16)&0xFF, (k>>8)&0xFF, k&0xFF)
}

func (v val) String() string {
	return fmt.Sprintf("(%c%c%c, %c%c%c)", (v>>40)&0xFF, (v>>32)&0xFF, (v>>24)&0xFF, (v>>16)&0xFF, (v>>8)&0xFF, v&0xFF)
}

// left returns the key that makes up the top 3 bytes of the val
func (v val) left() key {
	return key((v >> 24) & 0xFFFFFF)
}

// right returns the key that makes up the bottom 3 bytes of the val
func (v val) right() key {
	return key(v & 0xFFFFFF)
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int
	input = bytes.TrimSpace(input)

	// input is directions followed by nodes
	// RL
	//
	// AAA = (BBB, CCC)
	// BBB = (DDD, EEE)
	// CCC = (ZZZ, GGG)

	directions, nodeBytes, _ := bytes.Cut(input, []byte("\n\n"))
	numNodes := bytes.Count(nodeBytes, []byte{'\n'}) - 2

	// a node is just a map of 3 bytes with two endpoints, each 3 bytes
	nodes := make(map[key]val, numNodes)

	for line := range bytes.SplitSeq(nodeBytes, []byte{'\n'}) {
		// 0123456789ABCDEF
		// AAA = (BBB, CCC)
		key := key(uint32(line[0])<<16 | uint32(line[1])<<8 | uint32(line[2]))
		left := uint64(line[7])<<16 | uint64(line[8])<<8 | uint64(line[9])
		right := uint64(line[0xC])<<16 | uint64(line[0xD])<<8 | uint64(line[0xE])

		nodes[key] = val((left << 24) | right)
	}

	s1 = part1(directions, nodes)
	s2 = part2(directions, nodes)

	return s1, s2, nil
}

func part1(directions []byte, nodes map[key]val) int {
	steps := 0
	end := keyFromBytes([]byte("ZZZ"))
	key := keyFromBytes([]byte("AAA"))
	node, ok := nodes[key]
	if !ok {
		// no start node (only true in tests)
		return 0
	}

	i := 0
	for {
		// traverse each direction
		switch directions[i] {
		case 'L':
			key = node.left()
		case 'R':
			key = node.right()
		}

		node = nodes[key]
		steps++ // part1 solution is the number of steps
		if key == end {
			break
		}

		// directions loop
		i++
		if i >= len(directions) {
			i = 0
		}
	}
	return steps
}

func part2(directions []byte, nodes map[key]val) int {

	lcm := 0
	for k := range nodes {
		// for each node that ends in A, that's a ghost path, follow it
		if (k & 0xFF) == 'A' {
			steps := stepsToZ(directions, nodes, k)
			if lcm == 0 {
				lcm = steps
			} else {
				lcm = advent.LCM(lcm, steps)
			}
		}
	}

	return lcm
}

// / stepstoZ finds the steps one ghost takes to find a node ending in Z
func stepsToZ(directions []byte, nodes map[key]val, k key) int {
	steps := 0

	i := 0
	for {
		// traverse each direction
		switch directions[i] {
		case 'L':
			k = nodes[k].left()
		case 'R':
			k = nodes[k].right()
		}

		steps++

		if (k & 0xFF) == 'Z' {
			break
		}

		// directions loop
		i++
		if i >= len(directions) {
			i = 0
		}
	}

	return steps
}
