//go:build !arm64 || purego

package day1

import "testing"

func Test_solve(t *testing.T) {
	tests := []struct {
		name string
		input string
		want  int
		want2 int
	}{
		{ name: "123", input: "123\n", want: 123/3-2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := solve([]byte(tt.input))
			if got != tt.want {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("solve() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
