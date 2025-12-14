package day10

import (
	"slices"
	"testing"
)

func Test_lookAndSay(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{name: "1", want: []byte("11")},
		{name: "11", want: []byte("21")},
		{name: "1211", want: []byte("111221")},
		{name: "111221", want: []byte("312211")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookAndSay([]byte(tt.name)); !slices.Equal(got, tt.want) {
				t.Errorf("lookAndSay() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
