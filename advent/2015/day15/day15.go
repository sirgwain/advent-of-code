package day15

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	re := regexp.MustCompile(`([A-Za-z]+): capacity (-?\d+), durability (-?\d+), flavor (-?\d+), texture (-?\d+), calories (-?\d+)`)

	ingredients := make([][5]int, 0, 4)
	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input)))
	for scanner.Scan() {
		line := scanner.Text()

		m := re.FindStringSubmatch(line)
		if len(m) != 7 {
			return 0, 0, fmt.Errorf("can't parse line %s", string(line))
		}

		capacity, err := strconv.Atoi(m[2])
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}
		durability, err := strconv.Atoi(m[3])
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}
		flavor, err := strconv.Atoi(m[4])
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}
		texture, err := strconv.Atoi(m[5])
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}
		calories, err := strconv.Atoi(m[6])
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}

		ingredients = append(ingredients, [5]int{
			capacity,
			durability,
			flavor,
			texture,
			calories,
		})
	}

	s1 = bestScore(ingredients, -1)
	s2 = bestScore(ingredients, 500)

	return s1, s2, nil
}

// bestScore brute forces the ingredient combos to get
// at the best
func bestScore(ingredients [][5]int, caloriesRequired int) int {
	var best int
	for a := range 100 {
		for b := range 100 {
			for c := range 100 {
				d := 100 - a - b - c
				if d < 0 {
					continue
				}
				count := [4]int{a, b, c, d}

				// if we aren't in the calorie requirement, skip it
				if caloriesRequired != -1 && calculateCalories(ingredients, count) != caloriesRequired {
					continue
				}
				best = max(best, calculateScore(ingredients, count))
			}
		}
	}
	return best
}

func calculateScore(ingredients [][5]int, count [4]int) int {
	var cap, dur, fla, tex int
	for i, ing := range ingredients {
		cap += count[i] * ing[0]
		dur += count[i] * ing[1]
		fla += count[i] * ing[2]
		tex += count[i] * ing[3]
	}

	return max(0, cap) * max(0, dur) * max(0, fla) * max(0, tex)
}

func calculateCalories(ingredients [][5]int, count [4]int) int {
	var cal int
	for i, ing := range ingredients {
		cal += count[i] * ing[4]
	}
	return cal
}
