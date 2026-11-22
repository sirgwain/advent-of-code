//go:build !arm64 || purego

package day1

import (
	"fmt"

	"github.com/sirgwain/advent-of-code/advent"
)

func solve(input []byte) (s1, s2 int) {
	// start facing north
	hx, hy := 0, 1
	pt := advent.Point{}
	ptCounts := make(map[advent.Point]int, 1024)
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case 'R':
			hx, hy = hy, -hx
		case 'L':
			hx, hy = -hy, hx
		case ',':
			fallthrough
		case ' ':
			continue
		case '\n':
			continue
		default:
			dist, count := advent.Atoi(input[i:])
			i += count

			for j := 1; j <= dist; j++ {
				step := advent.Point{X: pt.X + hx*j, Y: pt.Y + hy*j}
				ptCounts[step]++
				if ptCounts[step] == 2 && s2 == 0 {
					s2 = advent.Abs(step.X) + advent.Abs(step.Y)
					if advent.Viz.Enabled {
						fmt.Printf("First duplicate point (%d,%d) %d away\n", step.X, step.Y, s2)
					}
				}
			}

			pt.X += hx * dist
			pt.Y += hy * dist

			if advent.Viz.Enabled {
				fmt.Printf("%s%d -> (%d, %d)\n", string(input[i-count-1:i-count]), dist, pt.X, pt.Y)
			}

		}
	}

	s1 = advent.Abs(pt.X) + advent.Abs(pt.Y)

	return s1, s2
}
