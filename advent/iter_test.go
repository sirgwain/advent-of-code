package advent_test

import (
	"slices"
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
)

func TestNumberFieldsSeq(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []int
	}{
		{"empty", "", nil},
		{"single", "42", []int{42}},
		{"spaces", "1 22 333", []int{1, 22, 333}},
		{"tabs", "5806\t6444\t1281", []int{5806, 6444, 1281}},
		{"leading and trailing", "  7  8  ", []int{7, 8}},
		{"trailing newline", "1\t2\n", []int{1, 2}},
		{"crlf", "1 2\r\n", []int{1, 2}},
		{"only separators", " \t \n", nil},
		{"zero is a number", "0 5 0", []int{0, 5, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(advent.NumberFieldsSeq([]byte(tt.line)))
			if !slices.Equal(got, tt.want) {
				t.Errorf("NumberFieldsSeq(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestNumberFieldsSeqBreak(t *testing.T) {
	var got []int
	for n := range advent.NumberFieldsSeq([]byte("1 2 3 4")) {
		got = append(got, n)
		if n == 2 {
			break
		}
	}
	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
