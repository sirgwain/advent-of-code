package day3

import "testing"

func Test_findHash(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		want int
	}{
		{name: "abcdef", key: []byte("abcdef"), want: 609043},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := findHash(tt.key, 1, false); got != tt.want {
				t.Errorf("findHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
