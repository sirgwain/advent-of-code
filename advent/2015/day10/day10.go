package day10

import (
	"bytes"
	"strconv"
	"strings"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {

	var s1, s2 int
	out := bytes.Clone(input)
	for i := range 50 {
		out = lookAndSay(bytes.TrimSpace(out))
		if i == 39 {
			s1 = len(out)
		}
	}

	s2 = len(out)

	return s1, s2, nil
}

func lookAndSay(in []byte) []byte {
	var sb strings.Builder

	var last byte
	count := 0
	for _, b := range in {
		// new number
		if last != b {
			if count > 0 {
				// record 3 1s as 31
				sb.WriteString(strconv.Itoa(count))
				sb.WriteByte(last)
			}
			last = b
			count = 0
		}
		count++
	}
	if count > 0 {
		// record 3 1s as 31
		sb.WriteString(strconv.Itoa(count))
	}
	sb.WriteByte(in[len(in)-1])

	return []byte(sb.String())
}
