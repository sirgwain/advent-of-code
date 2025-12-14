package day8

import (
	"testing"
)

func Test_parseLine(t *testing.T) {
	tests := []struct {
		name         string
		wantUnquoted int
		wantQuoted   int
	}{
		{name: `""`, wantUnquoted: 0, wantQuoted: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUnquoted, gotQuoted, err := unquotedQuoted(tt.name)
			if err != nil {
				t.Fatalf("unquotedQuoted() failed %v", err)
			}
			if gotUnquoted != tt.wantUnquoted || gotQuoted != tt.wantQuoted {
				t.Errorf("unquotedQuoted() = %v, %v, want %v, %v", gotUnquoted, gotQuoted, tt.wantUnquoted, tt.wantQuoted)
			}
		})
	}
}
