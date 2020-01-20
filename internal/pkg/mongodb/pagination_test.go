// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"testing"

	"github.com/sylabs/compute-service/internal/pkg/model"
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
		pa          model.PageArgs
		wantFirst   int
		wantLast    int
		wantAfter   primitive.ObjectID
		wantBefore  primitive.ObjectID
		wantErr     bool
	}{
		{"BadFirst", 10, model.PageArgs{First: &negativeOne}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"BadLast", 10, model.PageArgs{Last: &negativeOne}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"BadAfter", 10, model.PageArgs{After: &bad}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"BadBefore", 10, model.PageArgs{Before: &bad}, 0, 0, primitive.NilObjectID, primitive.NilObjectID, true},
		{"ZeroValues", 10, model.PageArgs{}, 10, 0, primitive.NilObjectID, primitive.NilObjectID, false},
		{"MaxPlusFirst", 10, model.PageArgs{First: &eleven}, 10, 0, primitive.NilObjectID, primitive.NilObjectID, false},
		{"MaxPlusLast", 10, model.PageArgs{Last: &eleven}, 0, 10, primitive.NilObjectID, primitive.NilObjectID, false},
		{"GoodFirst", 10, model.PageArgs{First: &five}, 5, 0, primitive.NilObjectID, primitive.NilObjectID, false},
		{"GoodLast", 10, model.PageArgs{Last: &five}, 0, 5, primitive.NilObjectID, primitive.NilObjectID, false},
		{"GoodAfter", 10, model.PageArgs{After: &idHex}, 10, 0, id, primitive.NilObjectID, false},
		{"GoodBefore", 10, model.PageArgs{Before: &idHex}, 10, 0, primitive.NilObjectID, id, false},
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
