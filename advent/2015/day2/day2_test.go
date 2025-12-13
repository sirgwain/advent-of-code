package day2

import "testing"

func Test_wrappingPaperSqFt(t *testing.T) {
	tests := []struct {
		name string
		l    int
		w    int
		h    int
		want int
	}{
		{name: "2x3x4", l: 2, w: 3, h: 4, want: 58},
		{name: "1x1x10", l: 1, w: 1, h: 10, want: 43},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrappingPaperSqFt(tt.l, tt.w, tt.h); got != tt.want {
				t.Errorf("wrappingPaper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ribbonLength(t *testing.T) {
	tests := []struct {
		name string
		l    int
		w    int
		h    int
		want int
	}{
		{name: "2x3x4", l: 2, w: 3, h: 4, want: 34},
		{name: "1x1x10", l: 1, w: 1, h: 10, want: 14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ribbonLength(tt.l, tt.w, tt.h); got != tt.want {
				t.Errorf("wrappingPaper() = %v, want %v", got, tt.want)
			}
		})
	}
}
