package day9

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code/advent"
	"github.com/sirgwain/advent-of-code/advent/color"
)

var nextIncrementStyle = lipgloss.NewStyle().Foreground(color.BrightGreen82)

type Day struct {
	processedNumbers [][]int
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	// input is a series of numbers
	// 0 3 6 9 12 15
	// 1 3 6 10 15 21
	// 10 13 16 21 30 45

	for line := range bytes.SplitSeq(bytes.TrimSpace(input), []byte{'\n'}) {
		n := bytes.Count(line, []byte{' '}) + 1
		numbers := make([]int, 0, n)
		for val := range bytes.FieldsSeq(line) {
			num, err := strconv.Atoi(string(val))
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse num %s %v", val, err)
			}
			numbers = append(numbers, num)
		}

		d.processedNumbers = append([][]int{}, numbers)
		nextIncrement := d.processNumbers(numbers)
		d.processedNumbers[0] = append(d.processedNumbers[0], numbers[len(numbers)-1]+nextIncrement)

		s1 += numbers[len(numbers)-1] + nextIncrement

		if advent.Viz.Enabled {
			d.view()
			fmt.Println()
		}
	}

	return s1, s2, nil
}

func (d *Day) view() {

	if advent.Viz.Enabled {
		for step, processed := range d.processedNumbers {
			// print the zeros
			fmt.Printf("%s", strings.Repeat(" ", step))
			for i := 0; i < len(processed)-1; i++ {
				fmt.Printf("%d ", processed[i])
			}
			fmt.Printf("%s", nextIncrementStyle.Render(fmt.Sprintf("%d", processed[len(processed)-1])))
			fmt.Println()
		}
	}
}

func (d *Day) processNumbers(numbers []int) int {
	processed := make([]int, len(numbers)-1)
	allZero := true
	for i := 1; i < len(numbers); i++ {
		processed[i-1] = numbers[i] - numbers[i-1]
		if processed[i-1] != 0 {
			allZero = false
		}
	}

	// record this for viewing
	step := len(d.processedNumbers)
	d.processedNumbers = append(d.processedNumbers, processed)

	if allZero {
		d.processedNumbers[step] = append(d.processedNumbers[step], 0)
		return 0
	}

	next := processed[len(processed)-1] + d.processNumbers(processed)
	d.processedNumbers[step] = append(d.processedNumbers[step], next)

	return next
}
