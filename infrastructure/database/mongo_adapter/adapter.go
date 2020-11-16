package mongoadapter

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAdapter
type MongoAdapter struct {
	User     string
	Password string
	Host     string
	Database string
	Timeout  int64
	client   mongo.Client
}

// StartDatabase
func (ma MongoAdapter) StartDatabase() error {
	if ma.Timeout < 1 {
		ma.Timeout = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ma.Timeout))
	defer cancel()

	connURI := fmt.Sprintf("mongodb://%s:%s@%s/%s", ma.User, ma.Password, ma.Host, ma.Database)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI))
	if err != nil {
		return err
	}
	ma.client = *client

	return nil
}
