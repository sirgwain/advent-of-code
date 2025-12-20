package day17

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strconv"
)

type Day struct {
}

type minContainer struct {
	numContainers int
	count         int
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	containers := []int{}
	scanner := bufio.NewScanner(bytes.NewReader(bytes.TrimSpace(input)))
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, 0, fmt.Errorf("unable to parse number %v", err)
		}

		containers = append(containers, val)
	}

	// start with 150l of eggnog
	amount := 150

	mc := minContainer{numContainers: math.MaxInt}
	s1 = fillContainers(amount, containers, &mc)
	s2 = mc.count

	return s1, s2, nil
}

// fillContainers recursively fills containers, trying smaller and smaller slices
// as long as they fit
func fillContainers(amount int, containers []int, mc *minContainer) int {
	var count int
	for i, container := range containers {
		count += findAllFitsForContainer(amount, container, containers[i+1:], 0, mc)
	}
	return count
}

// findAllFitsForContainer takes an amount and a container, and tries to pour it into the container, then fill all other containers
// it also tracks how many of the smallest number of containers we've found
func findAllFitsForContainer(amount int, container int, containers []int, numContainers int, mc *minContainer) (count int) {
	if container > amount {
		// too much container for the rest of the amount, we're done
		return 0
	}
	// poor into this container
	leftover := amount - container
	numContainers++
	if leftover == 0 {
		if numContainers < mc.numContainers {
			// reset the counter for mininum containers
			mc.numContainers = numContainers
			mc.count = 0
		}
		if numContainers == mc.numContainers {
			mc.count++
		}
		// all done
		return 1
	}

	// try each other container
	for i, c := range containers {
		if c > leftover {
			// container too big for remainder
			continue
		}

		count += findAllFitsForContainer(leftover, c, containers[i+1:], numContainers, mc)
	}

	return count
}
