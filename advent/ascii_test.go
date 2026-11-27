package advent_test

import (
	"strconv"
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
)

func TestIncr(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"zero", "0", "1"},
		{"single digit", "5", "6"},
		{"8 to 9", "8", "9"},
		{"9 to 10", "9", "10"},
		{"no carry", "41", "42"},
		{"single carry", "19", "20"},
		{"multi carry", "1299", "1300"},
		{"99 to 100", "99", "100"},
		{"999 to 1000", "999", "1000"},
		{"leading digit carry", "1999", "2000"},
		{"long all nines", "9999999", "10000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// exact capacity: rollover must grow via append
			if got := advent.Incr([]byte(tt.in)); string(got) != tt.want {
				t.Errorf("Incr(%q) = %q, want %q", tt.in, got, tt.want)
			}

			// spare capacity full of garbage: rollover must not expose it
			buf := make([]byte, len(tt.in), len(tt.in)+4)
			copy(buf, tt.in)
			for i := len(buf); i < cap(buf); i++ {
				buf[:cap(buf)][i] = 'x'
			}
			if got := advent.Incr(buf); string(got) != tt.want {
				t.Errorf("Incr(%q) with spare cap = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestIncrSequence(t *testing.T) {
	d := make([]byte, 1, 8)
	d[0] = '0'
	for n := 1; n <= 100000; n++ {
		d = advent.Incr(d)
		if string(d) != strconv.Itoa(n) {
			t.Fatalf("after %d increments got %q, want %q", n, d, strconv.Itoa(n))
		}
	}
}

func TestPackLetters(t *testing.T) {
	tests := []struct {
		name string
		a    byte
		b    byte
		want uint16
	}{
		{name: "ab=0x01", a: 'a', b: 'b', want: 0x01},
		{name: "ac=0x02", a: 'a', b: 'c', want: 0x02},
		{name: "yz=0b11000_11001", a: 'y', b: 'z', want: 0b11000_11001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := advent.PackLetters(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("PackLetters() = %v, want %v", got, tt.want)
			}
		})
	}
}
