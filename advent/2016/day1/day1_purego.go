//go:build !arm64 || purego

package day1

import (
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

func solve(input []byte) (s1, s2 int) {
	clear(grid[:])

	// start facing north
	hx, hy := 0, 1
	dist := 0

	gridLoc := 256*512 + 256
	grid[gridLoc>>6] = 1 // origin is already visited
	found := false

	pt := advent.Point{}
	var dir *byte
	for i := 0; i < len(input); i++ {
		c := byte('\n')
		if i < len(input) {
			c = input[i]
		}
		switch c {
		case 'R':
			hx, hy = hy, -hx
			dir = &input[i]
		case 'L':
			hx, hy = -hy, hx
			dir = &input[i]
		case ' ', '\n':
			if !found {
				// step through each grid point until we find a duplicate
				step := pt
				stride := hx + hy*512

				for j := 0; j < dist; j++ {
					step.X += hx
					step.Y += hy
					gridLoc += stride

					word := gridLoc >> 6
					mask := uint64(1) << uint(gridLoc&63)

					if grid[word]&mask != 0 {
						s2 = advent.Abs(step.X) + advent.Abs(step.Y)
						found = true
						break
					}

					grid[word] |= mask
				}
			}

			pt.X += hx * dist
			pt.Y += hy * dist

			if advent.Viz.Enabled {
				fmt.Printf("%v%d -> (%d, %d)\n", *dir, dist, pt.X, pt.Y)
			}

			dist = 0
		default:
			digit := int(input[i] - '0')
			if digit > 9 {
				continue
			}
			dist = dist*10 + digit
		}
	}

	s1 = advent.Abs(pt.X) + advent.Abs(pt.Y)

	return s1, s2
}
