// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"testing"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

func TestViewer(t *testing.T) {
	sc := "startCursor"
	ec := "endCursor"
	r := Resolver{
		p: &mockPersister{
			wp: model.WorkflowsPage{
				Workflows: []model.Workflow{
					{
						ID:   "id1",
						Name: "name1",
					},
					{
						ID:   "id2",
						Name: "name2",
					},
				},
				PageInfo: model.PageInfo{
					StartCursor:     &sc,
					EndCursor:       &ec,
					HasNextPage:     true,
					HasPreviousPage: false,
				},
				TotalCount: 2,
			},
		},
	}
	s, err := getSchema(&r)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
	}{
		{"OK"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
            query OpName {
              viewer {
                id
                login
                workflows {
                  edges {
                    cursor
                    node {
                      id
                      name
                    }
                  }
                  pageInfo {
                    startCursor
                    endCursor
                    hasNextPage
                    hasPreviousPage
                  }
                  totalCount
                }
              }
            }`

			res := s.Exec(context.Background(), q, "", nil)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}
