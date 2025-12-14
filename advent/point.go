package advent

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func PointFromString(s string) (Point, error) {

	xy := strings.Split(s, ",")
	x, err := strconv.Atoi(xy[0])
	if err != nil {
		return Point{}, fmt.Errorf("failed to parse %s %v", s, err)
	}
	y, err := strconv.Atoi(xy[1])
	if err != nil {
		return Point{}, fmt.Errorf("failed to parse %s %v", s, err)
	}

	return Point{x, y}, nil
}
