package day9

import (
	"bytes"
	"fmt"
	"maps"
	"math"
	"strconv"
	"strings"

	"github.com/sirgwain/advent-of-code/advent"
)

type Day struct {
}

func (d *Day) Run(input []byte) (int, int, error) {
	g := advent.NewGraph[string]()
	for line := range strings.SplitSeq(string(bytes.TrimSpace(input)), "\n") {
		// line is like
		// London to Dublin = 464
		fields := strings.Fields(line)
		if len(fields) != 5 {
			return 0, 0, fmt.Errorf("failed to parse line %s", line)
		}
		source := fields[0]
		dest := fields[2]
		distStr := fields[4]

		dist, err := strconv.Atoi(distStr)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse number in line %s %v", line, err)
		}
		g.AddEdge(source, dest, dist)
		g.AddEdge(dest, source, dist)
	}

	var s1, s2 int

	// find the shortest path that traverses all nodes
	s1 = math.MaxInt
	for _, n := range g.Nodes {
		s1 = min(s1, traverseShortest(n, make(map[string]bool)))
		s2 = max(s2, traverseLongest(n, make(map[string]bool)))
	}

	return s1, s2, nil
}

// pick the shortest distance to a non-visited node
// returning the distance
func traverseShortest[T any](n1 *advent.Node[T], visited map[string]bool) int {
	visited[n1.Key] = true

	shortest := math.MaxInt
	var nextNode *advent.Node[T]

	for _, edge := range n1.Edges {
		n := edge.OtherNode(n1)

		if visited[n.Key] {
			// already visited, don't traverse this edge
			continue
		}
		if edge.Weight < shortest {
			shortest = edge.Weight
			nextNode = n
		}
	}

	if nextNode == nil {
		return 0
	}

	return shortest + traverseShortest(nextNode, visited)
}

// for each edge, try it and pick the longest distance
func traverseLongest[T any](n1 *advent.Node[T], visited map[string]bool) int {
	visited[n1.Key] = true

	// try from each edge
	longest := 0
	for _, edge := range n1.Edges {
		n := edge.OtherNode(n1)

		if visited[n.Key] {
			// already visited, don't traverse this edge
			continue
		}

		dist := edge.Weight + traverseLongest(n, maps.Clone(visited))
		if dist > longest {
			longest = dist
		}
	}

	return longest
}
