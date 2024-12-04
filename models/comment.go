package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Comment represents a comment on a trip
type Comment struct {
    ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID primitive.ObjectID `bson:"user_id" json:"user_id"`  // User ID associated with the comment
    TripID primitive.ObjectID `bson:"trip_id" json:"trip_id"`  // Trip ID associated with the comment
    Content string            `bson:"content" json:"content"`  // Content of the comment
}
