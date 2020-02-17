// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/model"
	"gopkg.in/square/go-jose.v2"
)

const (
	pathWellKnownOAuth = "/.well-known/oauth-authorization-server"
	pathWellKnownOIDC  = "/.well-known/openid-configuration"
)

// getDiscoveryURIs returns a list of discovery URIs to try based on the supplied issuerURI.
//
// There are two specifications of concern here, OAuth 2.0 Authorization Server Metadata (RFC 8414)
// and OpenID Connect Discovery (OpenID.Discovery). Unfortunately, the construction of the
// discovery path varies when the issuer URI has a path component to it.
func getDiscoveryURIs(issuerURI string) ([]string, error) {
	issuerURI = strings.TrimSuffix(issuerURI, "/")
	u, err := url.Parse(issuerURI)
	if err != nil {
		return nil, err
	}

	paths := []string{issuerURI + pathWellKnownOIDC}
	if u.Path == "" {
		paths = append(paths, issuerURI+pathWellKnownOAuth)
	} else {
		u.Path = pathWellKnownOAuth + u.Path
		paths = append(paths, u.String())
	}
	return paths, nil
}

// discoverAuthMetadata attempts to discover metadata from an OAuth issuer using well-known URI
// discovery for OAuth 2.0 (RFC 8414) and OpenID Connect (OpenID.Discovery).
func discoverAuthMetadata(ctx context.Context, hc *http.Client, issuerURI string) (model.AuthMetadata, error) {
	uris, err := getDiscoveryURIs(issuerURI)
	if err != nil {
		return model.AuthMetadata{}, err
	}

	for _, uri := range uris {
		md, err := getAuthMetadata(ctx, hc, uri)
		if err == nil {
			return md, nil
		}
	}
	return model.AuthMetadata{}, errors.New("auth metadata discovery failed")
}

// getAuthMetadata gets metadata from uri as per the OAuth 2.0 Authorization Server Metadata
// specification (RFC 8414).
func getAuthMetadata(ctx context.Context, hc *http.Client, uri string) (md model.AuthMetadata, err error) {
	logrus.WithField("uri", uri).Info("getting auth metadata")
	defer func(t time.Time) {
		log := logrus.WithField("took", time.Since(t))
		if err != nil {
			log.WithError(err).Warn("failed to get auth metadata")
		} else {
			log.WithField("metadata", fmt.Sprintf("%#v", md)).Info("got auth metadata")
		}
	}(time.Now())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return model.AuthMetadata{}, err
	}
	res, err := hc.Do(req)
	if err != nil {
		return model.AuthMetadata{}, err
	}
	defer res.Body.Close()

	if code := res.StatusCode; (code / 100) != 2 {
		return model.AuthMetadata{}, fmt.Errorf("%d %s", code, http.StatusText(code))
	}

	if err := json.NewDecoder(res.Body).Decode(&md); err != nil {
		return model.AuthMetadata{}, err
	}
	return md, nil
}

// getKeySet gets a JSON Web Key Set from uri as per the JSON Web Key specification (RFC 7515).
func getKeySet(ctx context.Context, hc *http.Client, uri string) (ks jose.JSONWebKeySet, err error) {
	logrus.WithField("uri", uri).Info("getting key set")
	defer func(t time.Time) {
		log := logrus.WithField("took", time.Since(t))
		if err != nil {
			log.WithError(err).Warn("failed to get key set")
		} else {
			log.WithField("keySet", fmt.Sprintf("%+v", ks)).Info("got key set")
		}
	}(time.Now())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return jose.JSONWebKeySet{}, err
	}
	res, err := hc.Do(req)
	if err != nil {
		return jose.JSONWebKeySet{}, err
	}
	defer res.Body.Close()

	if code := res.StatusCode; (code / 100) != 2 {
		return jose.JSONWebKeySet{}, fmt.Errorf("%d %s", code, http.StatusText(code))
	}

	if err := json.NewDecoder(res.Body).Decode(&ks); err != nil {
		return jose.JSONWebKeySet{}, err
	}
	return ks, nil
}
