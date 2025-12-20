package day18

import (
	"bytes"
	"testing"

	"github.com/sirgwain/advent-of-code/advent"
)

func Test_toggleBoardLights(t *testing.T) {
	tests := []struct {
		name      string
		board     [][]byte
		times     int
		cornersOn bool
		want      [][]byte
	}{
		{
			name:  "turn off 1x1",
			board: [][]byte{{'#'}},
			times: 1,
			want:  [][]byte{{'.'}},
		},
		{
			name:  "leave off 1x1",
			board: [][]byte{{'.'}},
			times: 1,
			want:  [][]byte{{'.'}},
		},
		{
			name: "turn on 3 neighbors lit",
			board: [][]byte{
				[]byte(".#."),
				[]byte("#.#"),
				[]byte("..."),
			},
			times: 1,
			want: [][]byte{
				[]byte(".#."),
				[]byte(".#."),
				[]byte("..."),
			},
		},
		{
			name: "stable",
			board: [][]byte{
				[]byte(".#."),
				[]byte("#.#"),
				[]byte(".#."),
			},
			times: 10,
			want: [][]byte{
				[]byte(".#."),
				[]byte("#.#"),
				[]byte(".#."),
			},
		},
		{
			name: "test board",
			board: [][]byte{
				[]byte(".#.#.#"),
				[]byte("...##."),
				[]byte("#....#"),
				[]byte("..#..."),
				[]byte("#.#..#"),
				[]byte("####.."),
			},
			times: 4,
			want: [][]byte{
				[]byte("......"),
				[]byte("......"),
				[]byte("..##.."),
				[]byte("..##.."),
				[]byte("......"),
				[]byte("......"),
			},
		},
		{
			name: "test board corner on step 1",
			board: [][]byte{
				[]byte("##.#.#"),
				[]byte("...##."),
				[]byte("#....#"),
				[]byte("..#..."),
				[]byte("#.#..#"),
				[]byte("####.#"),
			},
			times:     1,
			cornersOn: true,
			want: [][]byte{
				[]byte("#.##.#"),
				[]byte("####.#"),
				[]byte("...##."),
				[]byte("......"),
				[]byte("#...#."),
				[]byte("#.####"),
			},
		},
		{
			name: "test board corner on",
			board: [][]byte{
				[]byte("##.#.#"),
				[]byte("...##."),
				[]byte("#....#"),
				[]byte("..#..."),
				[]byte("#.#..#"),
				[]byte("####.#"),
			},
			times:     5,
			cornersOn: true,
			want: [][]byte{
				[]byte("##.###"),
				[]byte(".##..#"),
				[]byte(".##..."),
				[]byte(".##..."),
				[]byte("#.#..."),
				[]byte("##...#"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toggleBoardLights(tt.times, tt.board, tt.cornersOn)
			fail := false
			for y := range tt.board {
				if !bytes.Equal(tt.board[y], tt.want[y]) {
					fail = true
					break
				}
			}

			if fail {
				t.Errorf("toggleBoardLights() = \n%v\n, want \n%v", advent.BoardToString(tt.board, drawBoardCell), advent.BoardToString(tt.want, drawBoardCell))
			}
		})
	}
}
