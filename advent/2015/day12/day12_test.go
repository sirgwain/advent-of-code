package day12

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed testdata/input.txt
var input []byte

func Test_sumNumbers(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{`[1,2,3]`, 6},
		{`{"a":2,"b":4}`, 6},
		{`[[[3]]]`, 3},
		{`{"a":{"b":4},"c":-1}`, 3},
		{`{"a":[-1,1]}`, 0},
		{`{"a":[-120,131]}`, 11},
		{`{"a":[-131,120]}`, -11},
		{`[-1,{"a":1}]`, 0},
		{`{}`, 0},
		{`[]`, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumNumbers([]byte(tt.name)); got != tt.want {
				t.Errorf("countNums() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sumNumbersJsonDecoder(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{`{"a":[1,"red",3]}`, 4},
		{`[1,{"c":"red","b":2},3]`, 4},
		{`[1,{"c":"red","b":2,"d":{"e":[1,2,3],"f":"red"}},3]`, 4},
		{`{"d":"red","e":[1,2,3,4],"f":5}`, 0},
		{`{"d":"red","e":[1,2,3,4,"red"],"f":5}`, 0},
		{`[1,"red",5]`, 6},
		{`{"a":{"b":4,"c":"red","e":{"f": 10}},"c":-1}`, -1},
		{`[1,2,3]`, 6},
		{`[1,2,3,"red"]`, 6},
		{`{"a":2,"b":4}`, 6},
		{`[[[3]]]`, 3},
		{`{"a":{"b":4},"c":-1}`, 3},
		{`{"a":[-1,1]}`, 0},
		{`{"a":[-120,131]}`, 11},
		{`{"a":[-131,120]}`, -11},
		{`[-1,{"a":1}]`, 0},
		{`{}`, 0},
		{`[]`, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sumNumbersJsonDecoder([]byte(tt.name))
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("countNums() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSumNumbers(b *testing.B) {

	in := bytes.TrimSpace(input)
	b.Run("sumNumbers", func(b *testing.B) {
		for b.Loop() {
			sumNumbers(in)
		}
	})

	b.Run("sumNumbersJsonDecoder", func(b *testing.B) {
		for b.Loop() {
			_, err := sumNumbersJsonDecoder(in)
			if err != nil {
				b.Fatalf("failed to run %v", err)
			}
		}
	})

}
