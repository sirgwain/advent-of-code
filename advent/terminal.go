package advent

import (
	"fmt"
)

func PrintSolution(s1, s2 string) {
	fmt.Printf("solution1: %s solution2: %s",
		solutionStyle.Render(s1),
		solutionStyle.Render(s2),
	)
}
