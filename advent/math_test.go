package advent_test

import(
	"github.com/sirgwain/advent-of-code/advent"
	"testing"
)

func TestClamp(t *testing.T) {
	tests := []struct {
		name string
		n    int
		mn   int
		mx   int
		want int
	}{
		{"in range", 5, 0, 10, 5},
		{"low", -1, 0, 10, 0},
		{"high", 11, 0, 10, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := advent.Clamp(tt.n, tt.mn, tt.mx)
			if got != tt.want {
				t.Errorf("Clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
