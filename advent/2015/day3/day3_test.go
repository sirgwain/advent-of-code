package day3

import (
	"testing"
)

func TestDay_part1(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  int
	}{
		{name: ">", input: []byte(">"), want: 2},
		{name: "^>v<", input: []byte("^>v<"), want: 4},
		{name: "^v^v^v^v^v", input: []byte("^v^v^v^v^v"), want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDay_part2(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  int
	}{
		{name: "^v", input: []byte("^v"), want: 3},
		{name: "^>v<", input: []byte("^>v<"), want: 3},
		{name: "^v^v^v^v^v", input: []byte("^v^v^v^v^v"), want: 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
