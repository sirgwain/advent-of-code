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
			if got := lookAndSay([]byte(tt.name), make([]byte, 0)); !slices.Equal(got, tt.want) {
				t.Errorf("lookAndSay() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func BenchmarkLookAndSayRepeat(b *testing.B) {

	b.Run("50 loops", func(b *testing.B) {

		// Two reusable buffers
		buf1 := make([]byte, 0, 5_000_000)
		buf2 := make([]byte, 0, 5_000_000)

		start := []byte("3113322113")

		for b.Loop() {
			buf1 := buf1[:0]
			buf2 := buf2[:0]
			lookAndSayRepeat(start, 50, buf1, buf2)
		}
	})

}
