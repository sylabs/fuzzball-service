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
	wp := model.WorkflowsPage{
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
	}

	cursor := "cursorValue"
	count := 2

	tests := []struct {
		name   string
		args   map[string]interface{}
		wantPA model.PageArgs
	}{
		{"NoArgs", nil, model.PageArgs{}},
		{"After", map[string]interface{}{"after": cursor}, model.PageArgs{After: &cursor}},
		{"Before", map[string]interface{}{"before": cursor}, model.PageArgs{Before: &cursor}},
		{"First", map[string]interface{}{"first": count}, model.PageArgs{First: &count}},
		{"Last", map[string]interface{}{"last": count}, model.PageArgs{Last: &count}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Resolver{
				p: &mockPersister{
					wantPA: tt.wantPA,
					wp:     wp,
				},
			}
			s, err := getSchema(&r)
			if err != nil {
				t.Fatal(err)
			}

			q := `
			query OpName($after: String, $before: String, $first: Int, $last: Int) {
			  viewer {
			    id
			    login
			    workflows(after: $after, before: $before, first: $first, last: $last) {
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

			res := s.Exec(context.Background(), q, "", tt.args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}
