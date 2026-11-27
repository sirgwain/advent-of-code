package advent

import "iter"

// NumberFieldsSeq yields each run of ASCII digits in line as an int.
// Any non-digit byte (space, tab, newline, comma, ...) separates numbers.
func NumberFieldsSeq(line []byte) iter.Seq[int] {
	return func(yield func(int) bool) {
		num, inNum := 0, false
		for _, c := range line {
			if d := c - '0'; d <= 9 {
				num = num*10 + int(d)
				inNum = true
				continue
			}
			if inNum {
				if !yield(num) {
					return
				}
				num, inNum = 0, false
			}
		}
		if inNum {
			yield(num)
		}
	}
}
