// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package schema

import (
	"github.com/graph-gophers/graphql-go"
)

// Get parses the GraphQL schema and attaches the given root resolver. It returns an error if the
// Go type signature of the resolvers does not match the schema. If nil is passed as the resolver,
// then the schema can not be executed, but it may be inspected (e.g. with ToJSON).
func Get(resolver interface{}) (*graphql.Schema, error) {
	s, err := schema()
	if err != nil {
		return nil, err
	}
	return graphql.ParseSchema(s, resolver, graphql.UseStringDescriptions(), graphql.UseFieldResolvers())
}
