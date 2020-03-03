// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/sylabs/fuzzball-service/internal/pkg/token"
)

// tokenHandler parses and validates a JSON Web Token (JWT) in the authorization bearer of the
// request. If a valid token is found, it adds it to the request context for use by next.
func (s *Server) tokenHandler(c Config, next http.Handler) http.Handler {
	jwt := token.NewMiddleware(token.MiddlewareOptions{
		Audience: c.OAuth2Audience,
		Issuer:   c.OAuth2IssuerURI,
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			alg, ok := t.Header["alg"]
			if !ok {
				return nil, errors.New("algorithm not present")
			}
			kid, haveKID := t.Header["kid"]
			for _, k := range s.authKeys.Keys {
				if (alg == k.Algorithm) && (!haveKID || (kid == k.KeyID)) {
					return k.Key, nil
				}
			}
			return nil, errors.New("key not found")
		},
	})
	return jwt.Handler(next)
}
