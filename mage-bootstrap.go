// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build bootstrap

package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() {
	os.Exit(mage.Main())
}
