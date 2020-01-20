// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type pageInfo struct {
	totalSize       int
	hasNextPage     bool
	hasPreviousPage bool
	startCursor     string
	endCursor       string
}

// parsePageOpts parses the page options, validating and converting fields as needed.
func parsePageOpts(maxPageSize int, pa model.PageArgs) (first, last int, after, before primitive.ObjectID, err error) {
	// Validate first.
	if pa.First != nil {
		if *pa.First < 0 {
			return 0, 0, primitive.NilObjectID, primitive.NilObjectID, fmt.Errorf("invalid 'first' field value: %v", pa.First)
		}
		if maxPageSize <= *pa.First {
			first = maxPageSize
		} else {
			first = *pa.First
		}
	}

	// Validate last.
	if pa.Last != nil {
		if *pa.Last < 0 {
			return 0, 0, primitive.NilObjectID, primitive.NilObjectID, fmt.Errorf("invalid 'last' field value: %v", pa.Last)
		}
		if maxPageSize <= *pa.Last {
			last = maxPageSize
		} else {
			last = *pa.Last
		}
	}

	// If neither first nor last were supplied, return maxPageSize elements.
	if first == 0 && last == 0 {
		first = maxPageSize
	}

	// Validate after.
	if pa.After != nil {
		after, err = primitive.ObjectIDFromHex(*pa.After)
		if err != nil {
			return 0, 0, primitive.NilObjectID, primitive.NilObjectID, fmt.Errorf("invalid 'after' field value: %v", err)
		}
	}

	// Validate before.
	if pa.Before != nil {
		before, err = primitive.ObjectIDFromHex(*pa.Before)
		if err != nil {
			return 0, 0, primitive.NilObjectID, primitive.NilObjectID, fmt.Errorf("invalid 'before' field value: %v", err)
		}
	}

	return first, last, after, before, nil
}
