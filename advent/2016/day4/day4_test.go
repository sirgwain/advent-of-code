package day4

import (
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
)

func TestDay_Run(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wants1 int
		wants2 int
	}{
		{name: "valid simple", input: "aaaaa-bbb-z-y-x-123[abxyz]\n", wants1: 123},
		{name: "invalid dupes", input: "aaaaa-bbb-z-y-x-123[aaxyz]\n", wants1: 0},
		{name: "invalid order", input: "aaaaa-bbb-zz-y-x-123[azbxy]\n", wants1: 0},
		{name: "invalid alphabetical order", input: "aaaaa-bbb-z-y-x-123[abxzy]\n", wants1: 0},
		{name: "b higher than a", input: "aaaaa-bbbbbb-x-y-z-123[baxyz]\n", wants1: 123},
		{name: "abcxy", input: "aaaaa-bbb-c-z-y-x-123[abcxy]\n", wants1: 123},
		{name: "invalid most common", input: "aaaaa-bbb-z-y-x-123[bxyzc]\n", wants1: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advent.Viz.Enabled = true
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

func Test_decrypt(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want string
	}{
		{"qzmt-zixmtkozy-ivhz", 343, "very encrypted name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decrypt(tt.name, tt.id)
			if got != tt.want {
				t.Errorf("decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
