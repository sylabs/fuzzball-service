// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/sylabs/compute-service/internal/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// getPipeline returns a MongoDB aggregation pipeline that implements filtered pagination.
//
// To obtain the page as well as metadata about it, the aggregation pipeline contains the following
// stages:
//
//  1. Apply filter (if supplied).
//  2. Two sub-pipelines:
//		2a. Count the total number of documents.
//		2b. Accumulate documents that match the parameters specified by opts.
//  3. Coalesce stage 2 sub-pipelines to form a coherent output document.
//
// Here be dragons... 🔥🐉🐉
func getPipeline(filter bson.M, first, last int, after, before primitive.ObjectID) mongo.Pipeline {
	pipeline := mongo.Pipeline{}

	// Pipeline stage 1.
	if filter != nil {
		pipeline = append(pipeline, bson.D{
			// Stage 1: Apply filter criteria.
			{Key: "$match", Value: filter},
		})
	}

	// Sub-pipeline stage 2a.
	p2A := mongo.Pipeline{
		{{Key: "$count", Value: "count"}},
	}

	// Sub-pipeline stage 2b.
	p2B := mongo.Pipeline{}

	// If the "after" cursor is provided, match ID > "after" cursor.
	if after != primitive.NilObjectID {
		p2B = append(p2B, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$gt", Value: after},
				}},
			}},
		})
	}

	// If the "before" cursor is provided, match ID < "before" cursor.
	if before != primitive.NilObjectID {
		p2B = append(p2B, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$lt", Value: before},
				}},
			}},
		})
	}

	// If the "first" argument is provided, sort ascending and limit to (first+1).
	if first > 0 {
		p2B = append(p2B, bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "_id", Value: 1},
			}},
		}, bson.D{
			{Key: "$limit", Value: first + 1},
		})
	}

	// If the "last" argument is provided, sort descending, limit to (last+1), and then sort
	// ascending.
	if last > 0 {
		p2B = append(p2B, bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "_id", Value: -1},
			}},
		}, bson.D{
			{Key: "$limit", Value: last + 1},
		}, bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "_id", Value: 1},
			}},
		})
	}

	// Add stage 2.
	pipeline = append(pipeline, bson.D{
		// Stage 2: Sub-pipelines.
		{Key: "$facet", Value: bson.D{
			// Stage 2a: Count the total number of documents matching the filter.
			{Key: "count", Value: p2A},
			// Stage 2b: Accumulate documents that match the parameters specified by opts.
			{Key: "results", Value: p2B},
		}},
	})

	// Add stage 3.
	pipeline = append(pipeline, bson.D{
		// Stage 3: Coalesce sub-pipelines to produce result.
		{Key: "$project", Value: bson.D{
			{Key: "count", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$count.count", 0}},
			}},
			{Key: "results", Value: "$results"},
		}},
	})

	return pipeline
}

// unmarshal iterates over the supplied array of BSON values, unmarshalling them into results. If
// results does not contain a pointer to a slice, an error is returned.
func unmarshal(rvs []bson.RawValue, results interface{}) error {
	if results == nil {
		return errors.New("results should be a pointer to a slice")
	}
	if k := reflect.TypeOf(results).Kind(); k != reflect.Ptr {
		return errors.New("results should be a pointer to a slice")
	}
	resultType := reflect.TypeOf(results).Elem()
	if k := resultType.Kind(); k != reflect.Slice {
		return errors.New("results should be a pointer to a slice")
	}

	// Unmarshal values.
	resultSlice := reflect.MakeSlice(resultType, len(rvs), len(rvs))
	for i, rv := range rvs {
		result := reflect.New(resultType.Elem())
		if err := rv.Unmarshal(result.Interface()); err != nil {
			return err
		}
		resultSlice.Index(i).Set(result.Elem())
	}
	reflect.ValueOf(results).Elem().Set(resultSlice)
	return nil
}

type pageResult struct {
	Count   int             `bson:"count"`
	Results []bson.RawValue `bson:"results"`
}

// getCursor returns the cursor value associated with the given BSON value.
func getCursor(rv bson.RawValue) (string, error) {
	raw, ok := rv.DocumentOK()
	if !ok {
		return "", errors.New("got non-document raw value")
	}

	rv, err := raw.LookupErr("_id")
	if err != nil {
		return "", err
	}

	id, ok := rv.ObjectIDOK()
	if !ok {
		return "", errors.New("failed to parse object ID")
	}
	return id.Hex(), nil
}

// getPageInfo transforms a given pageResult into a list of raw values and page info.
func getPageInfo(first, last int, pr pageResult) ([]bson.RawValue, model.PageInfo, error) {
	rvs := pr.Results
	var pi model.PageInfo

	// Determine whether there is a next page.
	if first != 0 && len(rvs) > first {
		if last <= 0 {
			pi.HasNextPage = true
		}
		rvs = rvs[:first]
	}

	// Determine whether there is a previous page.
	if last != 0 && len(rvs) > last {
		if first <= 0 {
			pi.HasPreviousPage = true
		}
		rvs = rvs[len(rvs)-last:]
	}

	if len(rvs) > 0 {
		// Get cursor value of first result.
		sc, err := getCursor(rvs[0])
		if err != nil {
			return nil, model.PageInfo{}, err
		}
		pi.StartCursor = &sc

		// Get cursor value of last result.
		ec, err := getCursor(rvs[len(rvs)-1])
		if err != nil {
			return nil, model.PageInfo{}, err
		}
		pi.EndCursor = &ec
	}

	return rvs, pi, nil
}

// findPageEx implements filtered pagination as described in the "Relay Cursor Connections
// Specification" found at https://facebook.github.io/relay/graphql/connections.htm and "Complete
// Connection Model" found at https://graphql.org/learn/pagination/.
func findPageEx(ctx context.Context, col *mongo.Collection, maxPageSize int, filter bson.M, pa model.PageArgs, results interface{}) (model.PageInfo, int, error) {
	// Ensure page options are valid.
	f, l, a, b, err := parsePageOpts(maxPageSize, pa)
	if err != nil {
		return model.PageInfo{}, 0, err
	}

	// Run aggregation pipeline.
	cur, err := col.Aggregate(ctx, getPipeline(filter, f, l, a, b))
	if err != nil {
		return model.PageInfo{}, 0, err
	}
	defer cur.Close(ctx)

	// Advance cursor to first (only) document.
	if ok := cur.Next(ctx); !ok {
		return model.PageInfo{}, 0, cur.Err()
	}

	// Unmarshal document.
	pr := pageResult{}
	if err := cur.Decode(&pr); err != nil {
		return model.PageInfo{}, 0, err
	}

	// Populate page info.
	rvs, pi, err := getPageInfo(f, l, pr)
	if err != nil {
		return model.PageInfo{}, 0, err
	}

	// Unmarshal results.
	if err := unmarshal(rvs, results); err != nil {
		return model.PageInfo{}, 0, err
	}
	return pi, pr.Count, nil
}
