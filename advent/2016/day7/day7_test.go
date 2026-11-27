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
		{name: "has abba", input: "abba[mnop]qrst\n", wants1: 1},
		{name: "has 2 abbas", input: "abba[mnop]qrst\nqrst[mnop]abba\n", wants1: 2},
		{name: "abba in hypernet", input: "abcd[bddb]xyyx\n", wants1: 0},
		{name: "aaaa invalid", input: "aaaa[qwer]tyui\n", wants1: 0},
		{name: "oxxo in string is valid", input: "ioxxoj[asdfgh]zxcvbn\n", wants1: 1},
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
