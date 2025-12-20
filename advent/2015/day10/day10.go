package day10

import (
	"bytes"
	"strconv"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {

	var s1, s2 int
	in := bytes.TrimSpace(input)

	// Two reusable buffers
	buf1 := make([]byte, 0, 5_000_000)
	buf2 := make([]byte, 0, 5_000_000)

	// call it 40 times for part 1
	cur := lookAndSayRepeat(in, 40, buf1, buf2)
	s1 = len(cur)

	// call it 10 more times (50 total) for part 2
	cur = lookAndSayRepeat(cur, 10, buf1, buf2)
	s2 = len(cur)

	return s1, s2, nil
}

// call lookAndSay count times returning the result
// takes two swappable buffers to prevent allocs
// NOTE, assumes even count as it always uses buf1
func lookAndSayRepeat(in []byte, count int, buf1 []byte, buf2 []byte) []byte {
	cur := in

	for range count {
		buf1 = lookAndSay(cur, buf1[:0])

		// swap: next round reads from a, writes into b
		cur = buf1
		buf1, buf2 = buf2, buf1

	}
	return cur
}

func lookAndSay(in, out []byte) []byte {

	last := in[0]
	var count byte = 1
	for i := 1; i < len(in); i++ {
		b := in[i]
		// keep counting the same number
		if b == last {
			count++
			continue
		}
		// record 111, as three one's = 31
		out = appendCount(out, count)
		out = append(out, last)
		last = b
		count = 1
	}
	out = appendCount(out, count)
	out = append(out, in[len(in)-1])

	return out
}

// fast path appendCount to handle majority of
// cases (sped up benchmarks from 24 runs to 42, not bad!)
func appendCount(out []byte, c byte) []byte {
	switch c {
	case 1:
		return append(out, '1')
	case 2:
		return append(out, '2')
	case 3:
		return append(out, '3')
	default:
		return strconv.AppendInt(out, int64(c), 10)
	}
}
