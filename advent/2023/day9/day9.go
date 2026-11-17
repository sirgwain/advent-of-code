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

		if advent.Viz.Enabled {
			d.processedNumbers = append([][]int{}, numbers)
		}
		l, r := d.processNumbers(numbers)
		if advent.Viz.Enabled {
			d.processedNumbers[0] = append(d.processedNumbers[0], numbers[len(numbers)-1]+r)
			d.processedNumbers[0] = append([]int{numbers[0] - l}, d.processedNumbers[0]...)
		}

		s1 += numbers[len(numbers)-1] + r
		s2 += numbers[0] - l

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
			fmt.Printf("%s", nextIncrementStyle.Render(fmt.Sprintf("%d ", processed[0])))
			for i := 1; i < len(processed)-1; i++ {
				fmt.Printf("%d ", processed[i])
			}
			fmt.Printf("%s", nextIncrementStyle.Render(fmt.Sprintf("%d", processed[len(processed)-1])))
			fmt.Println()
		}
	}
}

func (d *Day) processNumbers(numbers []int) (l, r int) {
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
	if advent.Viz.Enabled {
		d.processedNumbers = append(d.processedNumbers, processed)
	}

	if allZero {
		// add 0 ... 0 to the debug lists
		if advent.Viz.Enabled {
			d.processedNumbers[step] = append([]int{0}, append(d.processedNumbers[step], 0)...)
		}
		return 0, 0
	}

	l, r = d.processNumbers(processed)
	r += processed[len(processed)-1]
	l = processed[0] - l
	if advent.Viz.Enabled {
		d.processedNumbers[step] = append([]int{l}, append(d.processedNumbers[step], r)...)
	}
	return l, r
}
