package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"user_id"` // MongoDB will automatically assign this
	Email    string             `bson:"email" json:"email"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}
