// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import "github.com/sylabs/compute-service/internal/pkg/core"

type pageArgs struct {
	After  *string
	Before *string
	First  *int32
	Last   *int32
}

func convertPageArgs(args pageArgs) core.PageArgs {
	pa := core.PageArgs{
		After:  args.After,
		Before: args.Before,
	}
	if args.First != nil {
		first := int(*args.First)
		pa.First = &first
	}
	if args.Last != nil {
		last := int(*args.Last)
		pa.Last = &last
	}
	return pa
}
