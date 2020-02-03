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

	"github.com/sylabs/compute-service/internal/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var testWorkflowID = "testWorkflowID"

// getTestJob generates a job for use in testing. The attributes of the job are
// varied based on the value of i.
func getTestJob(i int32) model.Job {
	return model.Job{
		Name:       fmt.Sprintf("job-%02d", i),
		WorkflowID: testWorkflowID,
	}
}

// insertTestJob inserts a job into the DB.
func insertTestJob(t *testing.T, db *mongo.Database) model.Job {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt32)))
	if err != nil {
		t.Fatalf("failed to generate random int: %v", err)
	}
	j := getTestJob(int32(n.Int64()))
	sr, err := db.Collection(jobCollectionName).InsertOne(context.Background(), j)
	if err != nil {
		t.Fatalf("failed to insert: %s", err)
	}
	j.ID = sr.InsertedID.(primitive.ObjectID).Hex()
	return j
}

// deleteTestJob deletes a job.
func deleteTestJob(t *testing.T, db *mongo.Database, id string) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Fatalf("failed to parse object ID: %v", err)
	}
	m := bson.M{"_id": oid}
	if err := db.Collection(jobCollectionName).FindOneAndDelete(context.Background(), m).Err(); err != nil {
		t.Fatalf("failed to delete: %s", err)
	}
}

func TestCreateJob(t *testing.T) {
	orig := model.Job{
		ID:         "blahblah",
		WorkflowID: "workflowID",
		Name:       "test",
		Image:      "blah://test",
		Command: []string{
			"test",
		},
	}

	// Create should succeed.
	j, err := testConnection.CreateJob(context.Background(), orig)
	if err != nil {
		t.Fatalf("failed to create: %s", err)
	}
	defer deleteTestJob(t, testConnection.db, j.ID)

	// Verify returned workflow. Force ID since it is allocated by CreateWorkflow.
	orig.ID = j.ID
	if _, err := primitive.ObjectIDFromHex(j.ID); err != nil {
		t.Fatalf("job has invalid ID")
	}
	if got, want := j, orig; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Get should succeed.
	j, err = testConnection.GetJob(context.Background(), j.ID)
	if err != nil {
		t.Fatalf("failed to get: %s", err)
	}

	// Verify returned job.
	if got, want := j, orig; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDeleteJob(t *testing.T) {
	j := insertTestJob(t, testConnection.db)

	// Delete should succeed.
	if _, err := testConnection.deleteJob(context.Background(), j.ID); err != nil {
		t.Fatalf("failed to delete: %s", err)
	}

	// Get should fail.
	if _, err := testConnection.GetJob(context.Background(), j.ID); err == nil {
		t.Error("unexpected success")
	}

	// Delete should fail the second time.
	if _, err := testConnection.deleteJob(context.Background(), j.ID); err == nil {
		t.Error("unexpected success")
	}

	// Delete should fail with bad BSON ID.
	if _, err := testConnection.deleteJob(context.Background(), "oops"); err == nil {
		t.Error("unexpected success")
	}
}

func TestDeleteJobsByWorkflowID(t *testing.T) {
	j := insertTestJob(t, testConnection.db)

	// Delete should succeed.
	if err := testConnection.DeleteJobsByWorkflowID(context.Background(), testWorkflowID); err != nil {
		t.Fatalf("failed to delete: %s", err)
	}

	// Get should fail.
	if _, err := testConnection.GetJob(context.Background(), j.ID); err == nil {
		t.Error("unexpected success")
	}

	// Delete should fail the second time.
	if err := testConnection.DeleteJobsByWorkflowID(context.Background(), testWorkflowID); err != nil {
		t.Error("unexpected success")
	}
}

func TestGetJob(t *testing.T) {
	j := insertTestJob(t, testConnection.db)
	defer deleteTestJob(t, testConnection.db, j.ID)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"Found", j.ID, false},
		{"NotFound", primitive.NewObjectID().Hex(), true},
		{"BadID", "1234", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testConnection.GetJob(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetJobsByID(t *testing.T) {
	j1 := insertTestJob(t, testConnection.db)
	defer deleteTestJob(t, testConnection.db, j1.ID)
	j2 := insertTestJob(t, testConnection.db)
	defer deleteTestJob(t, testConnection.db, j2.ID)

	tests := []struct {
		name    string
		jobs    []string
		wantNum int
	}{
		{"FoundOne", []string{j1.ID}, 1},
		{"FoundNone", []string{"000000000000000000000000"}, 0},
		{"FoundTwo", []string{j1.ID, j2.ID}, 2},
		{"Search3Find2", []string{j1.ID, j2.ID, "000000000000000000000000"}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := testConnection.GetJobsByID(context.Background(), model.PageArgs{}, testWorkflowID, tt.jobs)
			if err != nil {
				t.Fatal(err)
			}

			if p.TotalCount != tt.wantNum {
				t.Fatalf("want %d results, got %d", tt.wantNum, p.TotalCount)
			}

		})
	}
}

func TestSetJobStatus(t *testing.T) {
	j := insertTestJob(t, testConnection.db)
	defer deleteTestJob(t, testConnection.db, j.ID)

	// Get should return an unfinished job.
	j, err := testConnection.GetJob(context.Background(), j.ID)
	if err != nil {
		t.Error("unexpected failure")
	}

	if j.Status != "" {
		t.Errorf("unexpected status: got %q, want %q", j.Status, "")
	}
	if j.ExitCode != 0 {
		t.Errorf("unexpected exit code: got %d, want %d", j.ExitCode, 0)
	}

	if err := testConnection.SetJobStatus(context.Background(), j.ID, "newStatus", 1); err != nil {
		t.Error("unexpected failure")
	}

	// Get should return a finished job.
	j, err = testConnection.GetJob(context.Background(), j.ID)
	if err != nil {
		t.Error("unexpected failure")
	}

	if j.Status != "newStatus" {
		t.Errorf("unexpected status: got %q, want %q", j.Status, "newStatus")
	}
	if j.ExitCode != 1 {
		t.Errorf("unexpected exit code: got %d, want %d", j.ExitCode, 1)
	}

}
