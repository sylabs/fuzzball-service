// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package token

import (
	"context"
	"reflect"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestContext(t *testing.T) {
	// Parse a token.
	tok, err := parseAndValidate(testToken, func(t *jwt.Token) (interface{}, error) {
		return testSigningKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	// Create a new context.
	ctx := NewContext(context.Background(), tok)

	// Obtain the token from the context, and make sure it matches.
	if got, ok := FromContext(ctx); !ok {
		t.Fatalf("token not found in context")
	} else if want := tok; !reflect.DeepEqual(got, want) {
		t.Errorf("got token %v, want %v", got, want)
	}
}

func TestGetClaims(t *testing.T) {
	// Parse a token.
	tok, err := parseAndValidate(testToken, func(t *jwt.Token) (interface{}, error) {
		return testSigningKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if got, want := tok.Claims(), testClaims; !reflect.DeepEqual(got, want) {
		t.Errorf("got claims %+v, want %+v", got, want)
	}
}
