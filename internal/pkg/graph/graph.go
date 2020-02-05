// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package graph

import (
	"errors"
	"fmt"
)

var (
	// ErrDuplicateVertex indicates that a vertex with an
	// identical id exists in the DAG.
	ErrDuplicateVertex = errors.New("vertex already exists")
	// ErrVertexNoExist indicates that an id does not correspond
	// to a vertex in the DAG.
	ErrVertexNoExist = errors.New("vertex does not exist")
	// ErrCycle indicates that there is a cycle within the graph
	ErrCycle = errors.New("cycle exists")
)

// DAG represents a Directed Acyclic Graph.
type DAG struct {
	nodes map[string][]string
}

// New returns a new DAG.
func New() DAG {
	return DAG{
		nodes: make(map[string][]string),
	}
}

// AddVertex will add a vertex to the DAG with the adjacent parent verticies.
// Edges are assumed to travel from parents to children.
// If a vertex already exists, an error will be returned
func (g *DAG) AddVertex(id string, parents []string) error {
	if _, ok := g.nodes[id]; ok {
		return fmt.Errorf("%q: %w", id, ErrDuplicateVertex)
	}
	g.nodes[id] = parents
	return nil
}

// Validate checks that all edges are valid.
// If a vertex references a non-existant vertex an error is returned.
// NOTE: this does not check for cycles. HasCycles does.
// this could maybe be a check within HasCycles and not part of the api?
func (g DAG) Validate() error {
	for _, p := range g.nodes {
		for _, n := range p {
			if _, ok := g.nodes[n]; !ok {
				return fmt.Errorf("%q: %w", n, ErrVertexNoExist)
			}
		}
	}
	return nil
}

func (g *DAG) checkSubTree(id string, visited map[string]bool, processing map[string]bool) bool {
	if visited[id] == false {

		visited[id] = true
		processing[id] = true

		for _, parent := range g.nodes[id] {
			if !visited[parent] && g.checkSubTree(parent, visited, processing) {
				return true
			} else if processing[parent] {
				return true
			}
		}

	}
	processing[id] = false
	return false
}

// HasCycle returns true if the graph contains a cycle.
func (g *DAG) HasCycle() bool {
	visited := make(map[string]bool)
	processing := make(map[string]bool)

	for id := range g.nodes {
		visited[id] = false
		processing[id] = false
	}

	for id := range g.nodes {
		if g.checkSubTree(id, visited, processing) {
			return true
		}
	}
	return false
}

// TopoSort performs a topological sort on the DAG and
// returns a slice of vertex ids in a valid topological order.
func (g DAG) TopoSort() ([]string, error) {
	retList := make([]string, 0)
	sources := make([]string, 0)
	inDegrees := make(map[string]int)
	for id, p := range g.nodes {
		if len(p) == 0 {
			sources = append(sources, id)
		}

		inDegrees[id] = len(p)
	}

	// ensure we have at least one source node
	// otherwise we must have a cycle
	if len(sources) == 0 {
		return nil, ErrCycle
	}

	for len(sources) > 0 {
		id := sources[0]
		sources = sources[1:]
		retList = append(retList, id)
		// this search can be sped up by creating a map of
		// parents to children
		for n, ps := range g.nodes {
			if n == id {
				continue
			}

			match := false
			for _, p := range ps {
				if p == id {
					match = true
					break
				}
			}

			// if id is not a parent of this node keep looking
			if !match {
				continue
			}

			inDegrees[n]--
			if inDegrees[n] == 0 {
				sources = append(sources, n)
			}
		}
	}

	// if our list is not the same size as our number
	// of nodes, then we have a cycle within our graph
	// preventing the above loop from successfully removing
	// the edges within the loop.
	if len(retList) != len(g.nodes) {
		return nil, ErrCycle
	}

	return retList, nil
}
