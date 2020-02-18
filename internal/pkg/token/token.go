// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// Package token defines middleware that parses and validates a bearer token and passes it via the
// request context.
package token

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

// Token represents a JWT Token.
type Token struct {
	*jwt.Token
}

// key is an unexported type for keys defined in this package. This prevents collisions with keys
// defined in other packages.
type key int

// tokenKey is the key for Token values in Contexts. It is unexported; clients use FromContext
// instead of using this key directly.
var tokenKey key

// NewContext returns a new Context that carries value t.
func NewContext(ctx context.Context, t *Token) context.Context {
	return context.WithValue(ctx, tokenKey, t)
}

// FromContext returns the Token stored in ctx, if any.
func FromContext(ctx context.Context) (*Token, bool) {
	t, ok := ctx.Value(tokenKey).(*Token)
	return t, ok
}

// Claims returns the claims contained in the token.
func (t *Token) Claims() *Claims {
	return t.Token.Claims.(*Claims)
}
