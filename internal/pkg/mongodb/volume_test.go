// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package mongodb

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/sylabs/compute-service/internal/pkg/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// getTestVolume generates a volume for use in testing. The attributes of the volume are
// varied based on the value of i.
func getTestVolume(i int32) core.Volume {
	return core.Volume{
		Name:       fmt.Sprintf("volume-%02d", i),
		WorkflowID: testWorkflowID,
	}
}

// insertTestVolume inserts a volume into the DB.
func insertTestVolume(t *testing.T, db *mongo.Database) core.Volume {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt32)))
	if err != nil {
		t.Fatalf("failed to generate random int: %v", err)
	}
	v := getTestVolume(int32(n.Int64()))
	sr, err := db.Collection(volumeCollectionName).InsertOne(context.Background(), v)
	if err != nil {
		t.Fatalf("failed to insert: %s", err)
	}
	v.ID = sr.InsertedID.(primitive.ObjectID).Hex()
	return v
}

// deleteTestVolume deletes a volume.
func deleteTestVolume(t *testing.T, db *mongo.Database, id string) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Fatalf("failed to parse object ID: %v", err)
	}
	m := bson.M{"_id": oid}
	if err := db.Collection(volumeCollectionName).FindOneAndDelete(context.Background(), m).Err(); err != nil {
		t.Fatalf("failed to delete: %s", err)
	}
}

func TestCreateVolume(t *testing.T) {
	orig := core.Volume{
		ID:         "blahblah",
		WorkflowID: "workflowID",
		Name:       "test",
		Type:       "EPHEMERAL",
	}

	// Create should succeed.
	v, err := testConnection.CreateVolume(context.Background(), orig)
	if err != nil {
		t.Fatalf("failed to create: %s", err)
	}
	defer deleteTestVolume(t, testConnection.db, v.ID)

	// Verify returned volume. Force ID since it is allocated by CreateVolume.
	orig.ID = v.ID
	if _, err := primitive.ObjectIDFromHex(v.ID); err != nil {
		t.Fatalf("volume has invalid ID")
	}
	if got, want := v, orig; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Get should succeed.
	v, err = testConnection.getVolume(context.Background(), v.ID)
	if err != nil {
		t.Fatalf("failed to get: %s", err)
	}

	// Verify returned volume.
	if got, want := v, orig; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDeleteVolumesByWorkflowID(t *testing.T) {
	v := insertTestVolume(t, testConnection.db)

	// Delete should succeed.
	if err := testConnection.DeleteVolumesByWorkflowID(context.Background(), testWorkflowID); err != nil {
		t.Fatalf("failed to delete: %s", err)
	}

	// Get should fail.
	if _, err := testConnection.getVolume(context.Background(), v.ID); err == nil {
		t.Error("unexpected success")
	}

	// Delete should fail the second time.
	if err := testConnection.DeleteVolumesByWorkflowID(context.Background(), testWorkflowID); err != nil {
		t.Error("unexpected success")
	}
}
