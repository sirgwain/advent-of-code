package advent

import "fmt"

type Node[T any] struct {
	Key   string
	Edges []*Edge[T]
}

type Edge[T any] struct {
	N1     *Node[T]
	N2     *Node[T]
	Weight int
}

type Graph[T any] struct {
	Nodes []*Node[T]

	nodesByKey map[string]*Node[T]
}

func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		nodesByKey: make(map[string]*Node[T]),
	}
}

// AddEdge adds an edge to the graph, adding the node if
// it doesn't already exist
func (g *Graph[T]) AddEdge(key1, key2 string, weight int) {
	n1, ok := g.nodesByKey[key1]
	if !ok {
		n1 = &Node[T]{Key: key1}
		g.Nodes = append(g.Nodes, n1)
		g.nodesByKey[key1] = n1
	}

	n2, ok := g.nodesByKey[key2]
	if !ok {
		n2 = &Node[T]{Key: key2}
		g.Nodes = append(g.Nodes, n2)
		g.nodesByKey[key2] = n2
	}

	edge := Edge[T]{N1: n1, N2: n2, Weight: weight}
	n1.Edges = append(n1.Edges, &edge)
}

func (n *Node[T]) String() string {
	return n.Key
}

func (n *Node[T]) Edge(other *Node[T]) *Edge[T] {
	for _, e := range n.Edges {
		if e.N2 == other {
			return e
		}
	}
	return nil
}

func (e *Edge[T]) OtherNode(n *Node[T]) *Node[T] {
	if e.N1 == n {
		return e.N2
	}
	return e.N1
}

func (e *Edge[T]) String() string {
	return fmt.Sprintf("%s -(%d)-> %s", e.N1.String(), e.Weight, e.N2.String())
}
