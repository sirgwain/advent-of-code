package day3_test

import (
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
	"github.com/sirgwain/advent-of-code/advent/2016/day3"
)

func TestDay_Run(t *testing.T) {
	tests := []struct {
		name string 
		input   string
		wants1    int
		wants2   int
	}{
		{name: "5+10 not > 25", input: "  5 10 25\n", wants1: 0},
		{name: "3 4 5 triangle", input: "  3 4 5\n", wants1: 1},
		{name: "three lines", input: "  3 4 5\n  30 40 50\n  32 45 54\n", wants1: 3, wants2: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advent.Viz.Enabled = true
			var d day3.Day
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
