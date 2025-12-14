package day6

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

type action int

const (
	turnOn action = iota
	turnOff
	toggle
)

type order struct {
	action action
	start  advent.Point
	end    advent.Point
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int
	board1 := advent.MakeBoard[bool](1000, 1000)
	board2 := advent.MakeBoard[int](1000, 1000)
	for _, line := range strings.Split(string(bytes.TrimSpace(input)), "\n") {
		order, err := parseLine(line)
		if err != nil {
			return 0, 0, err
		}
		runOrder1(board1, order)
		runOrder2(board2, order)
	}

	// count lights on
	for y := range len(board1) {
		for x := range len(board1[y]) {
			if board1[y][x] {
				s1++
			}
		}
	}

	// sum lights values
	for y := range len(board2) {
		for x := range len(board2[y]) {
			s2 += board2[y][x]
		}
	}

	return s1, s2, nil
}

func runOrder1(board [][]bool, o order) {
	x1 := min(o.start.X, o.end.X)
	y1 := min(o.start.Y, o.end.Y)
	x2 := max(o.start.X, o.end.X)
	y2 := max(o.start.Y, o.end.Y)
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			switch o.action {
			case turnOn:
				board[y][x] = true
			case turnOff:
				board[y][x] = false
			case toggle:
				board[y][x] = !board[y][x]
			}
		}
	}
}

func runOrder2(board [][]int, o order) {
	x1 := min(o.start.X, o.end.X)
	y1 := min(o.start.Y, o.end.Y)
	x2 := max(o.start.X, o.end.X)
	y2 := max(o.start.Y, o.end.Y)
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			switch o.action {
			case turnOn:
				board[y][x]++
			case turnOff:
				board[y][x] = max(0, board[y][x]-1)
			case toggle:
				board[y][x] += 2
			}
		}
	}
}

var re = regexp.MustCompile(`(turn on|turn off|toggle) (\d+,\d+) through (\d+,\d+)`)

func parseLine(line string) (order, error) {

	m := re.FindStringSubmatch(line)
	if len(m) != 4 {
		return order{}, fmt.Errorf("can't parse line %s", line)
	}

	var a action
	switch m[1] {
	case "turn on":
		a = turnOn
	case "turn off":
		a = turnOff
	case "toggle":
		a = toggle
	}

	start, err := advent.PointFromString(m[2])
	if err != nil {
		return order{}, fmt.Errorf("can't parse point from line %s %v", line, err)
	}
	end, err := advent.PointFromString(m[3])
	if err != nil {
		return order{}, fmt.Errorf("can't parse point from line %s %v", line, err)
	}

	return order{a, start, end}, nil
}
