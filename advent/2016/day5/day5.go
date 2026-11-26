package day5

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct{}

const hexDigits = "0123456789abcdef"

func (d *Day) RunS(input []byte) (string, string, error) {
	var s1, s2 string

	input = bytes.TrimSpace(input)

	start := int64(0)
	var hash [16]byte
	password1 := make([]byte, 0, 8)
	password2 := make([]byte, 8)
	count := 0
	for {
		start, hash = findHash(input, start+1, 0xFFFFF000)
		// hash is 16 bytes, each byte is 2 hex digits, so get the 3rd byte, low digit (for digit 6)
		sixthVal := hash[2] & 0x0F
		sixthDigit := hexDigits[sixthVal]
		seventhDigit := hexDigits[hash[3]&0xF0>>4]
		if len(password1) < 8 {
			// solution one is just the six digit in hex appendec
			password1 = append(password1, sixthDigit)
		}
		
		if sixthVal < 8 && password2[sixthVal] == 0 {
			// solution 2, the sixth digit is the index (if valid) and the seventh digit is the
			// password
			password2[sixthVal] = seventhDigit
			count++
			if count == 8 {
				break
			}
		}
	}

	s1 = string(password1)
	s2 = string(password2)

	return s1, s2, nil
}

func findHash(key []byte, start int64, mask uint32) (int64, [16]byte) {
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
			return count, hash
		}

		msg = msg[:keylen+len(advent.Incr(msg[keylen:]))]
	}
	return -1, [16]byte{}
}
