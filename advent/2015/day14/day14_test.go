package day14

import "testing"

func Test_calcDistance(t *testing.T) {
	tests := []struct {
		name  string
		t     int
		v     int
		tFly  int
		tRest int
		want  int
	}{
		{
			name:  "Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds.",
			t:     1000,
			v:     14,
			tFly:  10,
			tRest: 127,
			want:  1120,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcDistance(tt.t, tt.v, tt.tFly, tt.tRest); got != tt.want {
				t.Errorf("calcDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}
