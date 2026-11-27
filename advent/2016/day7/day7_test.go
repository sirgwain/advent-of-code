package day7

import (
	"testing"
)

func TestDay_Run(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wants1 int
		wants2 int
	}{
		{name: "has abba", input: "abba[mnop]qrst", wants1: 1},
		{name: "has 2 abbas", input: "abba[mnop]qrst\nqrst[mnop]abba", wants1: 2},
		{name: "multiple hypernets invalid", input: "abba[mnop]qrst[abba]abba", wants1: 0},
		{name: "multiple hypernets valid", input: "qrst[mnop]qrst[mnop]abba", wants1: 1},
		{name: "abba in hypernet", input: "abcd[bddb]xyyx", wants1: 0},
		{name: "aaaa invalid", input: "aaaa[qwer]tyui", wants1: 0},
		{name: "oxxo in string is valid", input: "ioxxoj[asdfgh]zxcvbn", wants1: 1},
		// part2
		{name: "aba + bab in hypernet", input: "aba[bab]xyz", wants2: 1},
		{name: "xyx but no yxy", input: "xyx[xyx]xyx", wants2: 0},
		{name: "eke with kek in hypernet", input: "aaa[kek]eke", wants2: 1},
		{name: "zbz with bzb in hypernet", input: "zazbz[bzb]cdb", wants2: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Day
			gots1, gots2, _ := d.Run([]byte(tt.input))
			if gots1 != tt.wants1 {
				t.Errorf("Run() s1 = %v, want %v", gots1, tt.wants1)
			}
			if gots2 != tt.wants2 {
				t.Errorf("Run() s2 = %v, want %v", gots2, tt.wants2)
			}
		})
	}
}
