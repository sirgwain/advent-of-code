package advent

import (
	"fmt"
	"strconv"
)

func PrintSolution(s1, s2 int) {
	fmt.Printf("solution1: %s solution2: %s",
		solutionStyle.Render(strconv.Itoa(s1)),
		solutionStyle.Render(strconv.Itoa(s2)),
	)
}
