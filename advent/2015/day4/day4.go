package day3

import (
	"bytes"
	"crypto/md5"
	"math"
	"strconv"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	s1 = findHash(bytes.TrimSpace(input), 1, false)
	s2 = findHash(bytes.TrimSpace(input), s1, true)

	return s1, s2, nil
}

func findHash(key []byte, start int, sixHexDigits bool) int {
	h := md5.New()
	var num [16]byte
	hash := make([]byte, 0, 16)

	for count := start; ; count++ {
		h.Reset()
		h.Write(key)

		b := strconv.AppendInt(num[:0], int64(count), 10)
		h.Write(b)

		// compute hash from new key
		hash = h.Sum(hash[:0]) // no alloc if cap(out) >= 16

		// check first 5 hex digits
		if hash[0] == 0 &&
			hash[1] == 0 &&
			(hash[2]&0xF0) == 0 && (!sixHexDigits || (hash[2]&0x0F) == 0) {
			return count
		}
		if count == math.MaxInt {
			return -1
		}
	}
}
