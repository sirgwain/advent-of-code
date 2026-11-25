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
		{ name: "123", input: "123\n", want: 123/3-2, want2: 51},
		{ name: "6", input: "6\n", want: 0, want2: 0},
		{ name: "9", input: "9\n", want: 1, want2: 1},
		{ name: "12", input: "12\n", want: 2, want2: 2},
		{ name: "33", input: "33\n", want: 9, want2: 10},
		{ name: "1969", input: "1969\n", want: 654, want2: 966},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := solve([]byte(tt.input))
			if got != tt.want {
				t.Errorf("solve() = %v, want1 %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("solve() = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}
