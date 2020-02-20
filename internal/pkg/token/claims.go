// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package token

import "github.com/dgrijalva/jwt-go"

// Claims type.
type Claims struct {
	jwt.StandardClaims
	UserID string `json:"uid,omitempty"`
}

// VerifyAudience compares the "aud" claim (if present) against cmp.
func (c Claims) VerifyAudience(s string) bool {
	if c.Audience != "" {
		return c.Audience == s
	}
	return true
}

// VerifyIssuer compares the "iss" claim (if present) against cmp.
func (c Claims) VerifyIssuer(s string) bool {
	if c.Issuer != "" {
		return c.Issuer == s
	}
	return true
}
