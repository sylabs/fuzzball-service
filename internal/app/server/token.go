// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"errors"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// tokenHandler parses and validates a JSON Web Token (JWT) in the authorization bearer of the
// request. If a valid token is found, it adds it to the request context for use by next.
func (s *Server) tokenHandler(c Config, next http.Handler) http.Handler {
	jwt := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
			// Audience and issuer must match our configuration.
			if !t.Claims.(jwt.MapClaims).VerifyAudience(c.OAuth2Audience, false) {
				return nil, errors.New("invalid audience claim")
			}
			if !t.Claims.(jwt.MapClaims).VerifyIssuer(c.OAuth2IssuerURI, false) {
				return nil, errors.New("invalid issuer claim")
			}

			// Algorithm and key ID must match a key in the server state.
			alg, ok := t.Header["alg"]
			if !ok {
				return nil, errors.New("algorithm not present")
			}
			kid, ok := t.Header["kid"]
			if !ok {
				return nil, errors.New("key ID not present")
			}
			for _, k := range s.authKeys.Keys {
				if alg == k.Algorithm && kid == k.KeyID {
					return k.Key, nil
				}
			}
			return nil, errors.New("key not found")
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			logrus.WithField("error", err).Warn("failed to validate JWT")
			http.Error(w, err, http.StatusUnauthorized)
		},
	})
	return jwt.Handler(next)
}
