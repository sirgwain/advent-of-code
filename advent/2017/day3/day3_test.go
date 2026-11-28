package day3

import (
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
)

func TestDay_Run(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wants1  int
		wants2  int
		wantErr bool
	}{
		{name: "1 at start", input: "1", wants1: 0},
		{name: "2 one place right", input: "2", wants1: 1},
		{name: "3 one place right, one up", input: "3", wants1: 2},
		{name: "12 3 steps away", input: "12", wants1: 3},
		{name: "1024 31 steps away", input: "1024", wants1: 31},
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
