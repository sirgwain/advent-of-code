package day8

import (
	"bytes"
	"fmt"
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
	return key((uint32(b[0]) << 16) | (uint32(b[1]) << 8) | uint32(b[0]))
}

func (k key) String() string {
	return fmt.Sprintf("%c%c%c", (k>>16)&0xf, (k>>8)&0xf, k&0xf)
}

func (v val) String() string {
	return fmt.Sprintf("(%c%c%c, %c%c%c)", (v>>40)&0xf, (v>>32)&0xf, (v>>24)&0xf, (v>>16)&0xf, (v>>8)&0xf, v&0xf)
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

	end := keyFromBytes([]byte("ZZZ"))
	key := keyFromBytes([]byte("AAA"))
	node := nodes[key]
	i := 0
	for {
		dir := directions[i]

		// traverse each direction
		switch dir {
		case 'L':
			key = node.left()
		case 'R':
			key = node.right()
		}

		node = nodes[key]
		s1++ // part1 solution is the number of steps
		if key == end {
			break
		}

		// directions loop
		i++
		if i >= len(directions) {
			i = 0
		}

	}

	return s1, s2, nil
}
