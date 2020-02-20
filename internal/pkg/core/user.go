// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

// User represents a user.
type User struct {
	ID    string `bson:"_id,omitempty"` // Unique user ID.
	Login string `bson:"login"`         // The username used to login.
}
