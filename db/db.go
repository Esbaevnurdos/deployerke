package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var TripCollection *mongo.Collection
var CommentCollection *mongo.Collection

// InitDB initializes MongoDB connection
func InitDB() error {
	// Set up a context with a timeout for MongoDB client connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf("Error creating MongoDB client: %v", err)
		return err
	}

	// Attempt to connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		return err
	}

	// Test the connection by pinging the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Error pinging MongoDB: %v", err)
		return err
	}

	// Initialize the collections
	UserCollection = client.Database("trip-planner").Collection("users")
	TripCollection = client.Database("trip-planner").Collection("trips")
	CommentCollection = client.Database("trip-planner").Collection("comments")

	log.Println("Connected to MongoDB successfully!")
	return nil
}
