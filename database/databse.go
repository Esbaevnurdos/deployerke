package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is the global variable to access the MongoDB database instance
var DB *mongo.Database

// Connect initializes the connection to MongoDB
func Connect() {
	// Replace with your MongoDB URI if it's not running locally
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Create a context with a 10-second timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Set the database to use
	DB = client.Database("trip_management") // Replace "trip_management" with your database name
	log.Println("Connected to MongoDB successfully")
}
