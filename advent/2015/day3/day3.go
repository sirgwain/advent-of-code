package day3

import "github.com/sirgwain/advent-of-code/advent"

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	s1 = part1(input)
	s2 = part2(input)

	return s1, s2, nil
}

func part1(input []byte) int {
	delivered := map[advent.Point]int{}
	p := advent.Point{}
	delivered[p]++
	for _, dir := range input {
		switch dir {
		case '>':
			p.X++
		case '<':
			p.X--
		case 'v':
			p.Y++
		case '^':
			p.Y--
		}
		delivered[p]++
	}

	return len(delivered)
}

// part2, santa and robo santa alternate following directions from the eggnogg'd elf
func part2(input []byte) int {
	delivered := map[advent.Point]int{}
	pSanta := advent.Point{}
	pRoboSanta := advent.Point{}
	delivered[pSanta]++
	delivered[pRoboSanta]++
	for i, dir := range input {
		var p *advent.Point
		if i%2 == 0 {
			p = &pSanta
		} else {
			p = &pRoboSanta
		}
		switch dir {
		case '>':
			p.X++
		case '<':
			p.X--
		case 'v':
			p.Y++
		case '^':
			p.Y--
		}
		delivered[*p]++
	}

	return len(delivered)
}
