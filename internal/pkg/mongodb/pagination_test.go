// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package mongodb

import (
	"reflect"
	"testing"

	"github.com/sylabs/fuzzball-service/internal/pkg/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParsePageOpts(t *testing.T) {
	id := primitive.NewObjectID()

	// We need pointers for PageArgs...
	bad := "bad"
	negativeOne := -1
	five := 5
	eleven := 11
	idHex := id.Hex()

	tests := []struct {
		name        string
		maxPageSize int
		pa          core.PageArgs
		wantFirst   int
		wantLast    int
		wantAfter   primitive.ObjectID
		wantBefore  primitive.ObjectID
		wantErr     bool
	}{
		{"BadFirst", 10, core.PageArgs{First: &negativeOne}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"BadLast", 10, core.PageArgs{Last: &negativeOne}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"BadAfter", 10, core.PageArgs{After: &bad}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"BadBefore", 10, core.PageArgs{Before: &bad}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"ZeroValues", 10, core.PageArgs{}, 10, 0, primitive.NilObjectID, primitive.NilObjectID, false},
		{"MaxPlusFirst", 10, core.PageArgs{First: &eleven}, 10, 0, primitive.NilObjectID, primitive.NilObjectID, false},
		{"MaxPlusLast", 10, core.PageArgs{Last: &eleven}, 0, 10, primitive.NilObjectID, primitive.NilObjectID, false},
		{"GoodFirst", 10, core.PageArgs{First: &five}, 5, 0, primitive.NilObjectID, primitive.NilObjectID, false},
		{"GoodLast", 10, core.PageArgs{Last: &five}, 0, 5, primitive.NilObjectID, primitive.NilObjectID, false},
		{"GoodAfter", 10, core.PageArgs{After: &idHex}, 10, 0, id, primitive.NilObjectID, false},
		{"GoodBefore", 10, core.PageArgs{Before: &idHex}, 10, 0, primitive.NilObjectID, id, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, l, a, b, err := parsePageOpts(tt.maxPageSize, tt.pa)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, wantErr %v", err, tt.wantErr)
			}

			if got, want := f, tt.wantFirst; got != want {
				t.Errorf("got first %v, want %v", got, want)
			}
			if got, want := l, tt.wantLast; got != want {
				t.Errorf("got last %v, want %v", got, want)
			}
			if got, want := a, tt.wantAfter; got != want {
				t.Errorf("got after %v, want %v", got, want)
			}
			if got, want := b, tt.wantBefore; got != want {
				t.Errorf("got before %v, want %v", got, want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	// Generate some test data.
	b, err := bson.Marshal(bson.D{{Key: "values", Value: []string{"one", "two"}}})
	if err != nil {
		t.Fatal(err)
	}
	res := struct {
		Values []bson.RawValue `bson:"values"`
	}{}
	if err := bson.Unmarshal(b, &res); err != nil {
		t.Fatal(err)
	}

	var result string
	var results []string

	tests := []struct {
		name        string
		rvs         []bson.RawValue
		results     interface{}
		wantResults []string
		wantErr     bool
	}{
		{"NilResult", res.Values, nil, nil, true},
		{"NotPointer", res.Values, result, nil, true},
		{"NotSlice", res.Values, &result, nil, true},
		{"NotSlice", res.Values, &results, []string{"one", "two"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unmarshal(tt.rvs, tt.results)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if got, want := results, tt.wantResults; !reflect.DeepEqual(got, want) {
					t.Errorf("got results %v, want %v", got, want)
				}
			}
		})
	}
}
