package day17

import (
	"math"
	"testing"
)

func Test_fillContainers(t *testing.T) {
	tests := []struct {
		name              string
		amount            int
		containers        []int
		want              int
		wantMinContainers int
	}{
		{name: "test data", amount: 25, containers: []int{20, 15, 10, 5, 5}, want: 4, wantMinContainers: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minContainer{numContainers: math.MaxInt}
			if got := fillContainers(tt.amount, tt.containers, &mc); got != tt.want {
				t.Errorf("fitContainers() = %v, want %v", got, tt.want)
			}
			if mc.count != tt.wantMinContainers {
				t.Errorf("fitContainers() minContainers.count = %v, want %v", mc.count, tt.wantMinContainers)
			}
		})
	}
}
