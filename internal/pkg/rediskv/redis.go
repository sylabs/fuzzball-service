// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package rediskv

import (
	"github.com/go-redis/redis"
)

// Connection is an active connection to a Redis key value store.
type Connection struct {
	rc *redis.Client
}

// NewConnection opens a new connection to a Redis key value store.
func NewConnection(a string) (c *Connection, err error) {
	rc := redis.NewClient(&redis.Options{
		Addr: a,
	})
	return &Connection{
		rc,
	}, nil
}

// Disconnect disconnects from the Redis key value store.
func (c *Connection) Disconnect() error {
	return c.rc.Close()
}

// Set will store the value at the supplied key.
func (c *Connection) Set(key, value string) error {
	return c.rc.Set(key, value, 0).Err()

}

// Append will append the value to the existing entry for the
// supplied key, or create a new one.
func (c *Connection) Append(key, value string) error {
	return c.rc.Append(key, value).Err()
}

// Get will retrieve the value at the supplied key.
// If the key is not found, "" is returned without an error.
func (c *Connection) Get(key string) (string, error) {
	v, err := c.rc.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return v, nil
}

// GetJobOutput retrieves the stored output of the job with the supplied id.
func (c *Connection) GetJobOutput(id string) (string, error) {
	return c.Get(id)
}
