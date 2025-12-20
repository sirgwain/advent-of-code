package day14

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Day struct {
}

type reindeer struct {
	name   string
	v      int
	tFly   int
	tRest  int
	dist   int
	points int
}

func (r *reindeer) move(second int) {
	if second%(r.tFly+r.tRest) < r.tFly {
		// we're in the first n seconds of a block, move this deer
		r.dist += r.v
	}
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	t := 2503
	deer := make([]*reindeer, 0, 10)
	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input)))
	for scanner.Scan() {
		line := scanner.Text()
		// line is like
		// Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds.
		//   0    1   2   3  4    5   6    7      8    9   10   11   12  13   14
		fields := strings.Fields(line)

		name := fields[0]
		vStr := fields[3]
		v, err := strconv.Atoi(vStr)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}

		tFlyStr := fields[6]
		tFly, err := strconv.Atoi(tFlyStr)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}

		tRestStr := fields[13]
		tRest, err := strconv.Atoi(tRestStr)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}

		s1 = max(s1, calcDistance(t, v, tFly, tRest))
		deer = append(deer, &reindeer{name: name, v: v, tFly: tFly, tRest: tRest})
	}

	// progress in seconds
	for second := range t {
		// find the deer that is the farthest along
		best := 0
		for _, r := range deer {
			// move each deer
			r.move(second)
			// find the best distance
			best = max(best, r.dist)
		}

		// the leading deer(s) get a point
		for _, r := range deer {
			if r.dist == best {
				r.points++
			}
		}
	}

	for _, r := range deer {
		s2 = max(s2, r.points)
	}

	return s1, s2, nil
}

// calcDistance returns the total distance traveled for time t based on velocity v and time flying/time resting
func calcDistance(t, v, tFly, tRest int) int {
	// dist is time flying "full time blocks" (fly time + rest time) * velocity + time flying for the remainder * velocity
	return (t/(tFly+tRest))*(v*tFly) + min(tFly, (t%(tFly+tRest)))*v
}
