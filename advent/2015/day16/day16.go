package day16

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	sue := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	re := regexp.MustCompile(`Sue \d+: (.*)`)
	num := 1

	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input)))
	for scanner.Scan() {
		line := scanner.Text()
		// line is like
		// Sue 1: goldfish: 6, trees: 9, akitas: 0
		m := re.FindStringSubmatch(line)
		if len(m) != 2 {
			return 0, 0, fmt.Errorf("can't parse line %s", string(line))
		}

		a := map[string]int{}
		for category := range strings.SplitSeq(m[1], ",") {
			attrs := strings.Split(category, ": ")
			val, err := strconv.Atoi(attrs[1])
			if err != nil {
				return 0, 0, fmt.Errorf("unable to parse number for %s %s %v", attrs[0], line, err)
			}
			a[strings.TrimSpace(attrs[0])] = val
		}
		if s1 == 0 && isSue(a, sue) {
			s1 = num
		}

		if s2 == 0 && isSueRetroEncabulator(a, sue) {
			s2 = num
		}
		num++

	}

	return s1, s2, nil
}

func isSue(a, sue map[string]int) bool {
	for k, v := range sue {
		if val, ok := a[k]; ok {
			if val != v {
				// the aunt has this key, but the value is wrong
				return false
			}
		}
	}
	return true
}

func isSueRetroEncabulator(a, sue map[string]int) bool {
	for k, suesReading := range sue {
		if val, ok := a[k]; ok {
			switch k {
			case "cats", "trees":
				// for cats and trees, the aunt must have a value reading greater than sue's reading
				if val <= suesReading {
					return false
				}
			case "pomeranians", "goldfish":
				// for pomeranians and golfish, the aunt must have a value reading less than sue's reading
				if val >= suesReading {
					return false
				}
			default:
				if val != suesReading {
					// the aunt has this key, but the value is wrong
					return false
				}
			}
		}
	}
	return true
}
