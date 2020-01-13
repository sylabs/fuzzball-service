// Use `go generate` to pack all *.graphql files under this directory (and sub-directories) into
// a binary format.
//
//go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package schema

import (
	"bytes"
	"fmt"
)

// String reads the .graphql schema files from the generated _bindata.go file, concatenating the
// files together into one string.
//
// If this method complains about not finding functions AssetNames() or MustAsset(),
// run `go generate` against this package to generate the functions.
func String() (string, error) {
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
