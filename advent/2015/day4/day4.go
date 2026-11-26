package day4

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	input = bytes.TrimSpace(input)

	s1 = findHash(input, 1, 0xFFFFF000)
	s2 = findHash(input, int64(s1), 0xFFFFFF00)

	return s1, s2, nil
}

func findHash(key []byte, start int64, mask uint32) int {
	keylen := len(key)
	buf := make([]byte, keylen, keylen+20)

	// put the key at the start of the buf
	copy(buf, key)

	// render the start to the message
	msg := strconv.AppendInt(buf, start, 10)

	for count := start; count < math.MaxInt64; count++ {
		// compute hash from new key
		hash := md5.Sum(msg)

		// check first 5 or 6 hex digits
		if (binary.BigEndian.Uint32(hash[:4]) & mask) == 0 {
			return int(count)
		}

		msg = msg[:keylen+len(advent.Incr(msg[keylen:]))]
	}
	return -1
}
