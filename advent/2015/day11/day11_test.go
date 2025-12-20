package day11

import "testing"

func Test_rotatePassword(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "a", want: "b"},
		{name: "z", want: "a"},
		{name: "az", want: "ba"},
		{name: "ghijklmn", want: "ghijklmo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rotatePassword([]byte(tt.name)); string(got) != tt.want {
				t.Errorf("rotatePassword() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "hijklmmn", want: false},
		{name: "abbceffg", want: false},
		{name: "abbcegjk", want: false},
		{name: "abcdffaa", want: true},
		{name: "ghjaabcc", want: true},
		{name: "ghjaabcd", want: false}, // only one pair
		{name: "ghjaabaa", want: false}, // only one unique pair
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePassword([]byte(tt.name)); got != tt.want {
				t.Errorf("validatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextPassword(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "abcdefgh", want: "abcdffaa"},
		{name: "ghijklmn", want: "ghjaabcc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextPassword([]byte(tt.name)); string(got) != tt.want {
				t.Errorf("nextPassword() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

