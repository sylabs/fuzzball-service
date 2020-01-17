// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// UserResolver resolves a user.
type UserResolver struct {
	u *model.User
}

// ID resolves the unique user ID.
func (r *UserResolver) ID() graphql.ID {
	return graphql.ID(r.u.ID)
}

// Login resolves the username used to login.
func (r *UserResolver) Login() string {
	return r.u.Login
}
