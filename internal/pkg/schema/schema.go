// Use `go generate` to pack all *.graphql files under this directory (and sub-directories) into
// a binary format.
//
//go:generate go-bindata -nometadata -pkg=schema ../../../api/graphql-schema

package schema

import (
	"bytes"
	"fmt"

	"github.com/graph-gophers/graphql-go"
)

// schema reads the .graphql schema files from the generated _bindata.go file, concatenating the
// files together into one string.
//
// If this method complains about not finding functions AssetNames() or MustAsset(),
// run `go generate` against this package to generate the functions.
func schema() (string, error) {
	buf := bytes.Buffer{}
	for _, name := range AssetNames() {
		b, err := Asset(name)
		if err != nil {
			return "", fmt.Errorf("asset: Asset(%q): %w", name, err)
		}

		buf.Write(b)

		// Add a newline if the file does not end in a newline.
		if len(b) > 0 && b[len(b)-1] != '\n' {
			buf.WriteByte('\n')
		}
	}

	return buf.String(), nil
}

// Get parses the GraphQL schema and attaches the given root resolver. It returns an error if the
// Go type signature of the resolvers does not match the schema. If nil is passed as the resolver,
// then the schema can not be executed, but it may be inspected (e.g. with ToJSON).
func Get(resolver interface{}) (*graphql.Schema, error) {
	s, err := schema()
	if err != nil {
		return nil, err
	}
	return graphql.ParseSchema(s, resolver, graphql.UseStringDescriptions())
}
