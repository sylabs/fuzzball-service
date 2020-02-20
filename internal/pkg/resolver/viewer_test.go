// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/sylabs/compute-service/internal/pkg/core"
	"github.com/sylabs/compute-service/internal/pkg/token"
)

func TestViewer(t *testing.T) {
	// User token to pass in context.
	tok := token.Token{
		Token: jwt.NewWithClaims(jwt.SigningMethodNone, &token.Claims{
			StandardClaims: jwt.StandardClaims{
				Subject: "jimbob",
			},
			UserID: "507f1f77bcf86cd799439011",
		}),
	}
	ctx := token.NewContext(context.Background(), &tok)

	sc := "startCursor"
	ec := "endCursor"
	wp := core.WorkflowsPage{
		Workflows: []core.Workflow{
			{
				ID:   "id1",
				Name: "name1",
			},
			{
				ID:   "id2",
				Name: "name2",
			},
		},
		PageInfo: core.PageInfo{
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
		wantPA core.PageArgs
	}{
		{"NoArgs", nil, core.PageArgs{}},
		{"After", map[string]interface{}{"after": cursor}, core.PageArgs{After: &cursor}},
		{"Before", map[string]interface{}{"before": cursor}, core.PageArgs{Before: &cursor}},
		{"First", map[string]interface{}{"first": count}, core.PageArgs{First: &count}},
		{"Last", map[string]interface{}{"last": count}, core.PageArgs{Last: &count}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, err := getMockCore(mockCore{
				p: mockPersister{
					wantPA: tt.wantPA,
					wp:     wp,
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			r := Resolver{s: mc}
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

			res := s.Exec(ctx, q, "", tt.args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}
