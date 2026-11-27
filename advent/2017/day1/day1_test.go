package day1

import "testing"

func Test_solve(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wants1 int
		wants2 int
	}{
		{name: "1122 = 3", input: "1122", wants1: 3},
		{name: "1111 = 4", input: "1111", wants1: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gots1, gots2 := solve([]byte(tt.input))
			if gots1 != tt.wants1 {
				t.Errorf("solve() s1 = %v, want %v", gots1, tt.wants1)
			}
			if gots2 != tt.wants2 {
				t.Errorf("solve() s2 = %v, want %v", gots2, tt.wants2)
			}
		})
	}
}
