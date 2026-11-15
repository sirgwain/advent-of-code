package day6

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1 int
	var s2 int

	// input is time of each race and best distance
	// Time:      7  15   30
	// Distance:  9  40  200

	lines := bytes.Split(bytes.TrimSpace(input), []byte{'\n'})
	i := 0
	times := make([]int, 0)
	p2time := ""
	for t := range bytes.FieldsSeq(lines[0]) {
		if i == 0 {
			i++
			continue
		}
		p2time += string(t)
		num, err := strconv.Atoi(string(t))
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse input %s %v", string(t), err)
		}
		times = append(times, num)
	}

	i = 0
	distances := make([]int, 0)
	p2distance := ""
	for d := range bytes.FieldsSeq(lines[1]) {
		if i == 0 {
			i++
			continue
		}
		p2distance += string(d)
		num, err := strconv.Atoi(string(d))
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse input %s %v", string(d), err)
		}
		distances = append(distances, num)
	}

	// solve each race, final solution is r1 * r2 * r3
	// i.e. all combinations of min time
	for i := range len(times) {
		t := times[i]
		d := distances[i]
		r := solveRace(t, d)
		if s1 == 0 {
			s1 = r
		} else {
			s1 *= r
		}
	}

	// part 2 merges all distances/times together as one big number each
	p2t, err := strconv.Atoi(p2time)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse p2time:%s %v", p2time, err)
	}
	p2d, err := strconv.Atoi(p2distance)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse p2distance:%s %v", p2distance, err)
	}

	// part 2 is just one big race
	s2 = solveRace(p2t, p2d)

	return s1, s2, nil
}

// solveRace determines the number of solutions for a race time t that will exceed distance d
// every ms you hold the boat button increases speed by 1mm/ms
// i.e. for a 7 ms race, if you hold the charge button for 2ms you go 2mm/ms * 5ms remaining = 10mm travelled
// solveRace returns the number of solutions that win the race (distance travelled > d)
//
// this is a quadratic equation. For the above 7ms race with a max distance of 9mm we want
// timeHeld * (7 - timeHeld) > 9
// -timeHeld^2 + 7*timeHeld -9 > 0
// x = (-b ± sqrt(b^2-4ac))/2a
// x = -t ± sqrt(-t^2 -4(-1)(-d))/-2
func solveRace(t, d int) int {

	x1 := (float64(-t) + math.Sqrt(float64(-t*-t-4*-1*-d))) / (-2)
	x2 := (float64(-t) - math.Sqrt(float64(-t*-t-4*-1*-d))) / (-2)

	// the above example solution is x1=1.69ms, x2=5.3ms
	// we solve for == 0, but we need >0 so add a tiny fraction
	// to account for whole edge cases
	solutions := int(math.Ceil(x2-.01)) - int(math.Ceil(x1+.01))

	if advent.Viz.Enabled {
		fmt.Printf("t=%dms, d=%dmm x1=%f, x2=%f, solutions=%d\n", t, d, x1, x2, solutions)
	}

	return solutions
}
