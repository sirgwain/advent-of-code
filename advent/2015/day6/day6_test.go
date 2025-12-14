package day6

import (
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
)

func Test_parseLine(t *testing.T) {
	tests := []struct {
		name string
		want order
	}{
		{name: "turn on 0,0 through 999,999", want: order{action: turnOn, start: advent.Point{X: 0, Y: 0}, end: advent.Point{X: 999, Y: 999}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLine(tt.name)
			if err != nil {
				t.Fatalf("parseLine() failed %v", err)
			}
			if got != tt.want {
				t.Errorf("parseLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
