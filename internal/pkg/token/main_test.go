// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package token

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	testSigningMethod = jwt.SigningMethodHS256
	testSigningKey    = []byte("AllYourBase")

	testClaims               *Claims
	testToken                string
	testTokenNoClaims        string
	testTokenInvalidAudience string
	testTokenInvalidIssuer   string
	testTokenExpired         string
)

func makeToken(c *Claims) string {
	s, err := jwt.NewWithClaims(testSigningMethod, c).SignedString(testSigningKey)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func TestMain(m *testing.M) {
	// Claims for testing.
	testClaims = &Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        "id",
			Subject:   "subject",
			Audience:  "api://default",
			Issuer:    "https://example.com",
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).UTC().Unix(),
		},
	}
	testToken = makeToken(testClaims)

	// No claims.
	testTokenNoClaims = makeToken(&Claims{})

	// Bad audience value.
	testTokenInvalidAudience = makeToken(&Claims{
		StandardClaims: jwt.StandardClaims{
			Audience: "bad",
		},
	})

	// Bad issuer value.
	testTokenInvalidIssuer = makeToken(&Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: "bad",
		},
	})

	// Expired.
	testTokenExpired = makeToken(&Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 1000000000,
		},
	})

	os.Exit(m.Run())
}
