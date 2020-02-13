// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/nats-io/nats.go"
	"gopkg.in/square/go-jose.v2"
)

var (
	natsURIs = flag.String("nats-uris", nats.DefaultURL, "Comma-separated list of NATS server URIs")

	testKeySet jose.JSONWebKeySet
)

func TestMain(m *testing.M) {
	flag.Parse()

	k, err := rsa.GenerateKey(rand.New(rand.NewSource(0)), 256)
	if err != nil {
		log.Fatal(err)
	}
	testKeySet = jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:          k.Public(),
				Certificates: []*x509.Certificate{},
				KeyID:        "0123456789abcdef",
				Algorithm:    string(jose.RS256),
				Use:          "sig",
			},
		},
	}

	os.Exit(m.Run())
}
