package day11

import (
	"bytes"
)

type Day struct {
}

func (d *Day) RunS(input []byte) (string, string, error) {
	pass1 := nextPassword(bytes.TrimSpace(input))
	pass2 := nextPassword(pass1)

	return string(pass1), string(pass2), nil
}

func nextPassword(pass []byte) []byte {
	for {
		pass = rotatePassword(pass)
		if validatePassword(pass) {
			break
		}
	}
	return pass
}

// rotate letters in password counting like
// xx, xy, xz, ya, yb, etc...
func rotatePassword(in []byte) []byte {
	for i := len(in) - 1; i >= 0; i-- {
		in[i] = in[i] + 1
		if in[i] == 'z'+1 {
			in[i] = 'a'
			continue
		}
		break
	}
	return in
}

func validatePassword(pass []byte) bool {

	hasStraight := false
	pairIndexes := make(map[uint16]int)
	pairs := make(map[uint16]bool)
	for i, b := range pass {
		// check invalid chars
		if b == 'i' || b == 'o' || b == 'l' {
			return false
		}

		if i < len(pass)-1 && pass[i] == pass[i+1] {
			// build a two byte number from the current character and the next one
			pair := uint16(pass[i])<<8 | uint16(pass[i+1])
			if j, ok := pairIndexes[pair]; !ok || i-j >= 2 {
				// don't count aaa as two pairs
				pairIndexes[pair] = i
				pairs[pair] = true
			}
		}

		if i < len(pass)-3 {
			hasStraight = hasStraight ||
				// check for abc, xyz, etc
				(pass[i] == pass[i+1]-1 && pass[i] == pass[i+2]-2)
		}
	}

	return hasStraight && len(pairs) > 1
}
