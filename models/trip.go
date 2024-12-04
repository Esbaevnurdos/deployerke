package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Trip represents a trip entry
type Trip struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Category    string             `json:"category" bson:"category"`
	Region      string             `json:"region" bson:"region"`
	Description string             `json:"description" bson:"description"`
	Attractions string             `json:"attractions" bson:"attractions"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
}
