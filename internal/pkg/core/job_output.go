// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

// JobOutputFetcher is the interface to fetch job output.
type JobOutputFetcher interface {
	GetJobOutput(string) (string, error)
}
