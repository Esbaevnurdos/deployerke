package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"trip-planner/db"
	"trip-planner/models"
	"trip-planner/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateComment adds a new comment to a specific trip
func CreateComment(w http.ResponseWriter, r *http.Request) {
	// Get the trip ID from the URL path
	tripID := mux.Vars(r)["trip_id"]
	objectID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		http.Error(w, "Invalid trip ID format", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the user ID from JWT token
	token := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert userID from string to ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Set the trip ID and user ID in the comment
	comment.TripID = objectID
	comment.UserID = userID
	comment.ID = primitive.NewObjectID() // Automatically generate a new ObjectID for the comment

	// Insert the comment into the database
	_, err = db.CommentCollection.InsertOne(context.Background(), comment)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Return the created comment
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetCommentByID retrieves a comment by its ID
func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	// Get the comment ID from the URL parameters
	commentID := mux.Vars(r)["id"]

	// Convert the comment ID from string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID format", http.StatusBadRequest)
		return
	}

	// Find the comment in the database by its ObjectID
	var comment models.Comment
	err = db.CommentCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Comment not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve comment", http.StatusInternalServerError)
		}
		return
	}

	// Return the comment as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// GetComments retrieves all comments for a specific trip
func GetComments(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["trip_id"]  // Use mux.Vars to get trip_id from URL
	objectID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		http.Error(w, "Invalid trip ID format", http.StatusBadRequest)
		return
	}

	var comments []models.Comment
	cursor, err := db.CommentCollection.Find(context.Background(), bson.M{"trip_id": objectID})
	if err != nil {
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var comment models.Comment
		if err := cursor.Decode(&comment); err != nil {
			http.Error(w, "Failed to decode comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Cursor error", http.StatusInternalServerError)
		return
	}

	// Return the comments as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	// Get the comment ID from the URL
	commentID := mux.Vars(r)["id"]

	// Convert the comment ID from string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID format", http.StatusBadRequest)
		return
	}

	// Decode the request body into an updated comment object
	var updatedComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&updatedComment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the content field is properly set in the updated comment
	if updatedComment.Content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	// Get the user ID from the JWT token
	token := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert userID from string to ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Ensure the comment belongs to the user and validate if the comment exists
	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	// Get the current comment from the database to preserve the trip_id if not provided in the update
	var currentComment models.Comment
	err = db.CommentCollection.FindOne(context.Background(), filter).Decode(&currentComment)
	if err != nil {
		http.Error(w, "Comment not found or you don't have permission", http.StatusNotFound)
		return
	}

	// If trip_id is provided in the request, we will use that; otherwise, we keep the existing one.
	tripID := currentComment.TripID
	if updatedComment.TripID != primitive.NilObjectID {
		tripID = updatedComment.TripID // If new trip_id is provided, update it
	}

	// Update the comment in the database
	update := bson.M{
		"$set": bson.M{
			"content":    updatedComment.Content, // Update only the content
			"updated_at": primitive.NewDateTimeFromTime(time.Now()), // Optional: Update timestamp
			"trip_id":    tripID, // Ensure trip_id is preserved or updated
		},
	}

	// Perform the update operation
	result, err := db.CommentCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	// Check if the update was successful
	if result.MatchedCount == 0 {
		http.Error(w, "Comment not found or you don't have permission", http.StatusNotFound)
		return
	}

	// Return the updated comment as a JSON response
	updatedComment.ID = objectID // Ensure ID is set for the response
	updatedComment.UserID = userID // Preserve the userID in the response
	updatedComment.TripID = tripID // Ensure trip_id is preserved

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedComment)
}




// DeleteComment deletes a comment by its ID
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID format", http.StatusBadRequest)
		return
	}

	// Get the user ID from JWT token
	token := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert userID from string to ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Ensure the comment belongs to the user
	filter := bson.M{"_id": objectID, "user_id": userID}

	// Delete the comment from the database
	_, err = db.CommentCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content as response after successful deletion
}
