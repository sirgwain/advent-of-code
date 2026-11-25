package day1

import (
	"fmt"
	"testing"
)

func TestDay_Run(t *testing.T) {
	tests := []struct {
		name   string
		wants1 int
		wants2 int
	}{
		{name: "R1", wants1: 1, wants2: 0},
		{name: "R2, L3", wants1: 5, wants2: 0},
		{name: "R8, R4, R4, R8", wants1: 8, wants2: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// advent.Viz.Enabled = true
			d := &Day{}
			gots1, gots2, _ := d.Run([]byte(tt.name))
			if gots1 != tt.wants1 {
				t.Errorf("Day.Run() gots1 = %v, wants1 %v", gots1, tt.wants1)
			}
			if gots2 != tt.wants2 {
				t.Errorf("Day.Run() gots2 = %v, wants2 %v", gots2, tt.wants2)
			}
			fmt.Println()
		})
	}
}
