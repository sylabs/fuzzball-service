// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package token

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// MiddlewareOptions control the behaviour of the token middleware.
type MiddlewareOptions struct {
	// The value to verify in the "aud" claim (if claim present).
	Audience string
	// The value to verify in the "iss" claim (if claim present).
	Issuer string
	// This callback function is used to supply the key for verification. The function receives the
	// parsed, but unverified Token. This allows you to use properties in the Header of the token
	// (such as `kid`) to identify which key to use. The algorithm specified in the token should be
	// verified to match the key.
	KeyFunc jwt.Keyfunc
}

// Middleware is a token middleware. Use the Handler method to obtain a http.Handler.
type Middleware struct {
	o MiddlewareOptions
}

// NewMiddleware returns a new token middleware.
func NewMiddleware(o MiddlewareOptions) *Middleware {
	return &Middleware{o: o}
}

// Handler returns a new token handler.
func (m *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := m.verifyJWT(r)
		if err != nil {
			logrus.WithError(err).Print("invalid token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// getBearerToken returns the bearer token found in the authorization header.
func getBearerToken(r *http.Request) (string, error) {
	// Extract token from request.
	ah := r.Header.Get("Authorization")
	if ah == "" {
		return "", errors.New("authorization header not present")
	}
	f := strings.Fields(ah)
	if len(f) != 2 || strings.ToLower(f[0]) != "bearer" {
		return "", errors.New("authorization header not valid")
	}
	return f[1], nil
}

// parseAndValidate parses and validate tokenString using the key supplied by kf.
func parseAndValidate(tokenString string, kf jwt.Keyfunc) (*Token, error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, kf)
	if err != nil {
		return nil, err
	}
	return &Token{t}, nil
}

// verifyJWT attempts to extract a bearer token from the authorization header of r.
//
// If a valid token is found, it is added to r.Context(). If no bearer token is present in r, this
// is not considered to be an error. If a bearer token is present but cannot be parsed/validated,
// an appropriate error is returned.
func (m *Middleware) verifyJWT(r *http.Request) error {
	// Get token from authorization header.
	tokenString, err := getBearerToken(r)
	if err != nil {
		return nil
	}

	// Parse the token.
	t, err := parseAndValidate(tokenString, m.o.KeyFunc)
	if err != nil {
		return err
	}

	// Validate the audience and issuer.
	if !t.Claims().VerifyAudience(m.o.Audience) {
		return errors.New("invalid audience in token")
	}
	if !t.Claims().VerifyIssuer(m.o.Issuer) {
		return errors.New("invalid issuer in token")
	}

	// Add the token to the request context.
	nr := r.WithContext(NewContext(r.Context(), t))
	*r = *nr
	return nil
}
