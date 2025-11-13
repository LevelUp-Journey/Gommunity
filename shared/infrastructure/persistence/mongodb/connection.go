package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type MongoConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoConnection creates and returns a new MongoDB connection
func NewMongoConnection(config MongoConfig) (*MongoConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(config.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(config.Database)

	log.Printf("Successfully connected to MongoDB database: %s", config.Database)

	return &MongoConnection{
		Client:   client,
		Database: database,
	}, nil
}

// Close closes the MongoDB connection
func (mc *MongoConnection) Close(ctx context.Context) error {
	if err := mc.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	log.Println("MongoDB connection closed")
	return nil
}

// GetCollection returns a collection from the database
func (mc *MongoConnection) GetCollection(collectionName string) *mongo.Collection {
	return mc.Database.Collection(collectionName)
}
