package nosql

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
	err    error
)

func InitNoSQLDatabase(uri string) error {
	// Connection is already initialized
	if client != nil {
		return nil
	}
	once.Do(func() {
		if uri == "" {
			err = fmt.Errorf("database URI cannot be empty")
			return
		}

		// Set a timeout for the connection attempt
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Set up MongoDB client options
		opts := options.Client().ApplyURI(uri)

		// Establish connection to MongoDB
		client, err = mongo.Connect(ctx, opts)
		if err != nil {
			return
		}

		// Ensure the connection is established
		if err = client.Ping(ctx, nil); err != nil {
			client = nil
			err = fmt.Errorf("failed to ping MongoDB: %v", err)
		}
	})
	return err
}

func GetClient() *mongo.Client {
	return client
}

func Close() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}
