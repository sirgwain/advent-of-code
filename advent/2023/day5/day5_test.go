package day5

import (
	"fmt"
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
)

func Test_mapSeedRangeToRange(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		seed         rng
		r            rng
		wantMapping  rng
		wantLeftover []rng
		wantMapped   bool
	}{
		// seed  |--------|    0 - 9
		// src      ||         3 - 4
		// ->
		// left  |-|           0 - 2
		// map      ||         100 - 101
		// right      |---|    5 - 9
		{
			name:         "contains",
			seed:         rng{src: 0, size: 10},                       // 0 - 9
			r:            rng{src: 3, dst: 100, size: 2},              // 3 - 4 -> 100 - 101
			wantMapping:  rng{src: 100, size: 2},                      // [100-101]
			wantLeftover: []rng{{src: 0, size: 3}, {src: 5, size: 5}}, // [0-2],  [5-9]
			wantMapped:   true,
		},
		// seed  ||              0 - 1
		// src             ||   10 - 11
		// ->
		// seed  ||              0 - 1 (unchanged)
		{
			name:        "out",
			seed:        rng{src: 0, size: 2},
			r:           rng{src: 10, dst: 20, size: 2},
			wantMapping: rng{src: 0, size: 2},
			wantMapped:  false,
		},
		// seed  ||    0 - 1
		// src   ||    0 - 1
		// ->
		// map   ||    10 - 11
		{
			name:        "in",
			seed:        rng{src: 0, size: 2},
			r:           rng{src: 0, dst: 10, size: 2},
			wantMapping: rng{src: 10, size: 2},
			wantMapped:  true,
		},
		// seed   ||     1 - 2
		// src   |---|   0 - 4
		// ->
		// map    ||     11 - 12
		{
			name:        "in-middle",
			seed:        rng{src: 1, size: 2},          // 1 - 2 (2)
			r:           rng{src: 0, dst: 10, size: 5}, // 0 - 4 -> 10 - 14 (5)
			wantMapping: rng{src: 11, size: 2},         // maps to 11 - 12 (2)
			wantMapped:  true,
		},
		// seed   |---|   1 - 5
		// src    ||      1 - 2
		// ->
		// map    ||      10 - 11
		// right    |-|   3 - 5
		{
			name:         "over-right",
			seed:         rng{src: 1, size: 5},          // 1 - 5
			r:            rng{src: 1, dst: 10, size: 2}, // 1 - 2 -> 10 - 11
			wantMapping:  rng{src: 10, size: 2},         // [10 - 11]
			wantLeftover: []rng{{src: 3, size: 3}},      // [3 - 5]
			wantMapped:   true,
		},
		// seed          |---|   8 - 12
		// src        |----|     5 - 10
		// ->
		// map           |-|     3 - 5
		// right            ||   11 - 12
		{
			name:         "over-right mid",
			seed:         rng{src: 8, size: 5},         // 8 - 12
			r:            rng{src: 5, dst: 0, size: 6}, // 5 - 10 -> 0 - 5
			wantMapping:  rng{src: 3, size: 3},         // [3 - 5]
			wantLeftover: []rng{{src: 11, size: 2}},    // [11 - 12]
			wantMapped:   true,
		},
		// seed   |---|          1 - 5
		// src     |--------|    2 - 11
		// ->
		// left   |              1 - 1
		// map     |--|          10 - 13
		{
			name:         "over-left",
			seed:         rng{src: 1, size: 5},           // 1 - 5
			r:            rng{src: 2, dst: 10, size: 10}, // 2 - 11 -> 10 - 19
			wantMapping:  rng{src: 10, size: 4},          // [10 - 13]
			wantLeftover: []rng{{src: 1, size: 1}},       // [1 - 1]
			wantMapped:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMapping, gotLeftover, mapped := mapSeedRangeToRange(tt.seed, tt.r)

			if mapped != tt.wantMapped {
				t.Errorf("wanted mapped: %v, got %v", tt.wantMapped, mapped)
				return
			}

			if tt.wantMapping != gotMapping {
				t.Errorf("wanted mapping: %v, got: %v", tt.wantMapping, gotMapping)
			}

			if len(tt.wantLeftover) != len(gotLeftover) {
				t.Errorf("wanted %d leftover seed ranges, got %d", len(tt.wantLeftover), len(gotLeftover))
				return
			}

			for i := 0; i < len(tt.wantLeftover); i++ {
				if tt.wantLeftover[i] != gotLeftover[i] {
					t.Errorf("%d wanted: %v, got: %v", i, tt.wantLeftover[i], gotLeftover[i])
				}
			}

		})
	}
}

func Test_processSeedRange(t *testing.T) {
	tests := []struct {
		name   string
		s      rng
		ranges []rng
		want   []rng
		wantOk bool
	}{
		{
			name:   "no mapping",
			s:      rng{src: 0, size: 2},
			ranges: []rng{{src: 10, dst: 20, size: 10}},
			want:   []rng{{src: 0, size: 2}},
		},
		{
			name:   "single mapping",
			s:      rng{src: 0, size: 2},
			ranges: []rng{{src: 0, dst: 10, size: 10}},
			want:   []rng{{src: 10, size: 2}},
			wantOk: true,
		},
		{
			name:   "split mapping",
			s:      rng{src: 2, size: 10},
			ranges: []rng{{src: 0, dst: 10, size: 10}},
			want:   []rng{{src: 12, size: 8}, {src: 10, size: 2}},
			wantOk: true,
		},
		{
			name: "split mapping 2 ranges",
			s:    rng{src: 2, size: 10},
			// 2 ranges, one applies to the first part, the other applies to the leftover
			// map 	0 - 9 (10)
			// 		10 - 19 (10)
			ranges: []rng{
				{src: 0, dst: 10, size: 10},
				{src: 10, dst: 20, size: 10},
			},
			want:   []rng{{src: 12, size: 8}, {src: 20, size: 2}},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advent.Viz.Enabled = true
			fmt.Println(seedRangeStyle.Render(seedRangeString(tt.s.src, tt.s.size)))

			got, gotOk := processSeedRange(tt.s, tt.ranges)
			if gotOk != tt.wantOk {
				t.Errorf("processSeedRange() = %v, wantOk %v", gotOk, tt.wantOk)
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("processSeedRange() = %v, want %v", got, tt.want)
				return
			}
			for i, g := range got {
				if tt.want[i] != got[i] {
					t.Errorf("processSeedRange()  %d = %v, want %v", i, g, tt.want[i])
				}
			}
		})
	}
}
