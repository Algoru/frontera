package mongoadapter

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoAdapter
type MongoAdapter struct {
	User       string
	Password   string
	Host       string
	Database   string
	AuthSource string
	Timeout    time.Duration
	client     mongo.Client
}

// StartDatabase
func (ma *MongoAdapter) StartDatabase() error {
	if ma.Timeout < 1*time.Second {
		ma.Timeout = 10 * time.Second
	}

	if ma.AuthSource == "" {
		ma.AuthSource = "admin"
	}

	ctx, cancel := context.WithTimeout(context.Background(), ma.Timeout)
	defer cancel()

	connURI := fmt.Sprintf(
		"mongodb://%s:%s@%s/%s?authSource=%s",
		ma.User, ma.Password, ma.Host, ma.Database, ma.AuthSource,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI))
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	ma.client = *client

	return nil
}
