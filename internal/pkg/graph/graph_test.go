// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package graph

import (
	"reflect"
	"testing"
)

type testVertex struct {
	id      string
	parents []string
}

func TestAddEdge(t *testing.T) {
	g := New()
	// Add should succeed.
	if err := g.AddVertex("one", nil); err != nil {
		t.Fatalf("failed to add: %s", err)
	}

	// Add should fail.
	if err := g.AddVertex("one", nil); err == nil {
		t.Fatalf("Unexpected success")
	}

	// Add should succeed.
	if err := g.AddVertex("two", nil); err != nil {
		t.Fatalf("failed to add: %s", err)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		vertices []testVertex
		wantErr  bool
	}{
		{
			"SingleVertex",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
			},
			false,
		},
		{
			"ParentChild",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one"},
				},
			},
			false,
		},
		{
			"MultipleChildren",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one"},
				},
				testVertex{
					"three",
					[]string{"one"},
				},
			},
			false,
		},
		{
			"MultipleParents",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					nil,
				},
				testVertex{
					"three",
					[]string{"one", "two"},
				},
			},
			false,
		},
		{
			"Disconnected",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					nil,
				},
				testVertex{
					"three",
					[]string{"one"},
				},
			},
			false,
		},
		{
			"SingleVertexErr",
			[]testVertex{
				testVertex{
					"one",
					[]string{"something"},
				},
			},
			true,
		},
		{
			"MultipleVertexErr",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one", "something"},
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// populate a graph with verticies
			g := New()
			for _, v := range tt.vertices {
				if err := g.AddVertex(v.id, v.parents); err != nil {
					t.Fatal("failed to construct graph: %w", err)
				}
			}

			err := g.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHasCycle(t *testing.T) {

	tests := []struct {
		name        string
		vertices    []testVertex
		expectCycle bool
	}{
		{
			"SingleVertex",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
			},
			false,
		},
		{
			"ParentChild",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one"},
				},
			},
			false,
		},
		{
			"MultipleChildren",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one"},
				},
				testVertex{
					"three",
					[]string{"one"},
				},
			},
			false,
		},
		{
			"MultipleParents",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					nil,
				},
				testVertex{
					"three",
					[]string{"one", "two"},
				},
			},
			false,
		},
		{
			"Disconnected",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					nil,
				},
				testVertex{
					"three",
					[]string{"one"},
				},
			},
			false,
		},
		{
			"Cycle",
			[]testVertex{
				testVertex{
					"one",
					[]string{"two"},
				},
				testVertex{
					"two",
					[]string{"one"},
				},
			},
			true,
		},
		{
			"ComplexCycle",
			[]testVertex{
				testVertex{
					"one",
					[]string{"four"},
				},
				testVertex{
					"two",
					[]string{"one"},
				},
				testVertex{
					"three",
					[]string{"two"},
				},
				testVertex{
					"four",
					[]string{"two"},
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// populate a graph with verticies
			g := New()
			for _, v := range tt.vertices {
				if err := g.AddVertex(v.id, v.parents); err != nil {
					t.Fatal("failed to construct graph: %w", err)
				}
			}

			cycle := g.HasCycle()
			if cycle != tt.expectCycle {
				t.Fatalf("got %t, want %t", cycle, tt.expectCycle)
			}
		})
	}
}

func TestTopoSort(t *testing.T) {
	tests := []struct {
		name           string
		vertices       []testVertex
		expectCycle    bool
		possibleOrders [][]string
	}{
		{
			"SingleVertex",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
			},
			false,
			[][]string{
				[]string{"one"},
			},
		},
		{
			"ParentChild",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one"},
				},
			},
			false,
			[][]string{
				[]string{"one", "two"},
			},
		},
		{
			"MultipleChildren",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one"},
				},
				testVertex{
					"three",
					[]string{"one"},
				},
			},
			false,
			[][]string{
				[]string{"one", "two", "three"},
				[]string{"one", "three", "two"},
			},
		},
		{
			"MultipleParents",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					nil,
				},
				testVertex{
					"three",
					[]string{"one", "two"},
				},
			},
			false,
			[][]string{
				[]string{"one", "two", "three"},
				[]string{"two", "one", "three"},
			},
		},
		{
			"Disconnected",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					nil,
				},
				testVertex{
					"three",
					[]string{"one"},
				},
			},
			false,
			[][]string{
				[]string{"one", "two", "three"},
				[]string{"two", "one", "three"},
			},
		},
		{
			"ComplexGraph",
			[]testVertex{
				testVertex{
					"one",
					nil,
				},
				testVertex{
					"two",
					[]string{"one"},
				},
				testVertex{
					"three",
					[]string{"one"},
				},
				testVertex{
					"four",
					[]string{"one", "two"},
				},
				testVertex{
					"five",
					[]string{"one", "two", "three"},
				},
			},
			false,
			[][]string{
				[]string{"one", "two", "three", "four", "five"},
				[]string{"one", "two", "three", "five", "four"},
				[]string{"one", "two", "four", "three", "five"},
				[]string{"one", "three", "two", "four", "five"},
				[]string{"one", "three", "two", "five", "four"},
			},
		},
		{
			"Cycle",
			[]testVertex{
				testVertex{
					"one",
					[]string{"two"},
				},
				testVertex{
					"two",
					[]string{"one"},
				},
			},
			true,
			nil,
		},
		{
			"ComplexCycle",
			[]testVertex{
				testVertex{
					"one",
					[]string{"four"},
				},
				testVertex{
					"two",
					[]string{"one"},
				},
				testVertex{
					"three",
					[]string{"two"},
				},
				testVertex{
					"four",
					[]string{"two"},
				},
			},
			true,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// populate a graph with verticies
			g := New()
			for _, v := range tt.vertices {
				if err := g.AddVertex(v.id, v.parents); err != nil {
					t.Fatal("failed to construct graph: %w", err)
				}
			}

			s, err := g.TopoSort()
			if err != nil {
				if tt.expectCycle {
					if err == ErrCycle {
						return
					}
				} else {
					t.Fatalf("got err %v", err)
				}
			} else {
				if tt.expectCycle {
					t.Fatalf("got err %v, expect cycle %v", err, tt.expectCycle)
				}
			}

			match := false
			for _, o := range tt.possibleOrders {
				if reflect.DeepEqual(s, o) {
					match = true
					break
				}
			}

			if !match {
				t.Fatalf("invalid ordering: %v", s)
			}
		})
	}
}
