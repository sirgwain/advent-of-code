package day5

import "testing"

func Test_isNice1(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ugknbfddgicrmopn", want: true},
		{name: "aaa", want: true},
		{name: "jchzalrnumimnmhp", want: false},
		{name: "haegwjzuvuyypxyu", want: false},
		{name: "dvszwmarrgswjxmb", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNice1(tt.name); got != tt.want {
				t.Errorf("isNice1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNice2(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "qjhvhtzxzqqjkmpb", want: true},
		{name: "xxyxx", want: true},
		{name: "uurcxstgmygtbstg", want: false},
		{name: "ieodomkazucvgmuy", want: false},
		{name: "aaa", want: false},
		{name: "qpnxkuldeiituggg", want: false},
		{name: "aaaa", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNice2(tt.name); got != tt.want {
				t.Errorf("isNice2() = %v, want %v", got, tt.want)
			}
		})
	}
}
