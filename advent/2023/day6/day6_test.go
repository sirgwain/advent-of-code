package day6

import "testing"

func Test_solveRace(t *testing.T) {
	tests := []struct {
		name string
		t    int
		d    int
		want int
	}{
		// {name: "7ms 9mm", t: 7, d: 9, want: 4}, // 4 solutions, 2ms - 6ms
		// {name: "15ms 40mm", t: 15, d: 40, want: 8},
		{name: "30ms 200mm", t: 30, d: 200, want: 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solveRace(tt.t, tt.d)
			if got != tt.want {
				t.Errorf("solveRace() = %v, want %v", got, tt.want)
			}
		})
	}
}
