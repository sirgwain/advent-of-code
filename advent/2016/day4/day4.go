package day4

import (
	"fmt"
	"slices"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct{}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	var id int
	var letterCounts [26]uint16
	start := 0
	end := 0
	for i := 0; i < len(input); i++ {
		b := input[i]
		letter := b - 'a'
		if letter < 26 {
			letterCounts[letter]++
		}

		digit := b - '0'
		if digit <= 9 {
			end = i - 1
			for digit <= 9 {
				id = 10*id + int(digit)
				i++
				b = input[i]
				digit = b - '0'
			}
		}

		if b == '[' {
			i++
			checksumStart := i

			// convert letterCounts to be a 2-byte value with
			// the high byte being the count, low byte is the letter
			// |byte count|byte letter|
			for l := range uint16(26) {
				letterCounts[l] = letterCounts[l]<<8 | (25 - l)
			}

			slices.Sort(letterCounts[:])

			// build the checksum by grabbing the last 5 letters in reverse order
			// convert them back to letters
			var checksum [5]byte
			for j := range len(checksum) {
				checksum[j] = 25 - byte(letterCounts[len(letterCounts)-1-j]&0xFF) + 'a'
			}

			// checksum is the last 5 letters in []
			// aaaaa-bbb-z-y-x-123[abxyz]
			valid := true
			for j := range 5 {
				if checksum[j] != input[i+j] {
					valid = false
				}
			}
			i += 6
			if valid {
				s1 += id
			}
			// run to the end of the line
			for input[i] != '\n' {
				i++
			}

			// had to look all through the outputs for this...
			decrypted := decrypt(string(input[start:end]), id)
			if decrypted == "northpole object storage" {
				s2 = id
			}

			// process this line
			if advent.Viz.Enabled {
				checksumEnd := i - 1
				for l := range letterCounts {
					if letterCounts[l] == 0 {
						continue
					}
					fmt.Printf("%c:%d ", rune('a'+l), letterCounts[l])
				}

				validString := "invalid"
				if valid {
					validString = "valid"
				}
				fmt.Printf(" id: %d [%s] -> %s \"%s\" %s\n", id, string(input[checksumStart:checksumEnd]), string(checksum[:]), decrypted, validString)
			}

			// reset for the next run
			clear(letterCounts[:])
			id = 0
			start = i + 1
		}
	}

	return s1, s2, nil
}

// decrypt for part 2 rotates each letter by the id, rolling over
func decrypt(str string, id int) string {
	input := []byte(str)
	for i := range input {
		if input[i] == '-' {
			input[i] = ' '
			continue
		}
		input[i] = 'a' + byte((int(input[i]-'a')+id)%26)
	}
	return string(input)
}
