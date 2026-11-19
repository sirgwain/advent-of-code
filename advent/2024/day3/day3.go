package day3

import (
	"bytes"
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

var (
	mulTok  = []byte("mul(")
	doTok   = []byte("do()")
	dontTok = []byte("don't()")
)

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	enabled := true
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case 'm':
			x, y, n, ok := parseMul(input[i:])
			if ok {
				if advent.Viz.Enabled {
					fmt.Printf("%s = %d (enabled=%v)\n", input[i:i+n], x*y, enabled)
				}
				s1 += x * y
				if enabled {
					s2 += x * y
				}
				i += n - 1 // i++ will still happen
			}
		case 'd':
			if bytes.HasPrefix(input[i:], doTok) {
				if advent.Viz.Enabled {
					fmt.Println("enabled")
				}
				enabled = true
				i += len(doTok) - 1
			}
			if bytes.HasPrefix(input[i:], dontTok) {
				enabled = false
				if advent.Viz.Enabled {
					fmt.Println("disabled")
				}
				i += len(dontTok) - 1
			}
		}
	}

	return s1, s2, nil
}

// parseMul parsed a mul(\d{1..3},\d{1..3}) byte stream
// returning the numbers and the total length and true if
// it parsed it
func parseMul(b []byte) (x, y, n int, ok bool) {
	if !bytes.HasPrefix(b, mulTok) {
		// not our guy
		return
	}

	n = len(mulTok) // mul(
	x, size := advent.Atoi(b[n:])
	if size < 1 || size > 3 {
		// reject numbers outside the range
		return
	}

	n += size // mul(123
	if n >= len(b) || b[n] != ',' {
		return
	}
	n++ // mul(123,

	y, size = advent.Atoi(b[n:])
	if size < 1 || size > 3 {
		// reject numbers outside the range
		return
	}

	n += size // mul(123,456

	if n >= len(b) || b[n] != ')' {
		return
	}
	n++ // mul(123,456)

	return x, y, n, true
}
