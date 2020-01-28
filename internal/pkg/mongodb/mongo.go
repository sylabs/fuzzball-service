// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxPageSize = 100

// Connection is an active connection to a MongoDB database.
type Connection struct {
	db *mongo.Database
}

// NewConnection opens a new connection to a MongoDB database.
func NewConnection(ctx context.Context, mongoURI, dbName string) (c *Connection, err error) {
	o := options.Client().ApplyURI(mongoURI)
	if err := o.Validate(); err != nil {
		return nil, err
	}
	mc, err := mongo.NewClient(o)
	if err != nil {
		return nil, err
	}
	if err := mc.Connect(ctx); err != nil {
		return nil, err
	}
	if err := mc.Ping(ctx, nil); err != nil {
		return nil, err
	}

	c = &Connection{
		db: mc.Database(dbName),
	}
	return c, nil
}

// Disconnect disconnects from the MongoDB database.
func (c *Connection) Disconnect(ctx context.Context) error {
	return c.db.Client().Disconnect(ctx)
}
