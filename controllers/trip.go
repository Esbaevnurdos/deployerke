package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"trip-planner/db"
	"trip-planner/models"
	"trip-planner/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTrip(w http.ResponseWriter, r *http.Request) {
    var trip models.Trip
    err := json.NewDecoder(r.Body).Decode(&trip)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Get the user ID from the token
    userID, err := getUserIDFromToken(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Set the user_id for the trip
    trip.UserID = userID
    trip.ID = primitive.NewObjectID() // Ensure the ID is generated

    // Insert trip into the database
    _, err = db.TripCollection.InsertOne(context.Background(), trip)
    if err != nil {
        http.Error(w, "Failed to create trip", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(trip)
}

func getUserIDFromToken(r *http.Request) (primitive.ObjectID, error) {
    tokenString := r.Header.Get("Authorization")
    claims, err := utils.ValidateJWT(tokenString) // Validate the JWT token
    if err != nil {
        return primitive.NilObjectID, err
    }
    userID, err := primitive.ObjectIDFromHex(claims.UserID) // Extract the ObjectID from the claims
    return userID, err
}



func GetTripByID(w http.ResponseWriter, r *http.Request) {
    tripID := mux.Vars(r)["id"]
    tripObjID, err := primitive.ObjectIDFromHex(tripID)
    if err != nil {
        http.Error(w, "Invalid trip ID format", http.StatusBadRequest)
        return
    }

    userID, err := getUserIDFromToken(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var trip models.Trip
    err = db.TripCollection.FindOne(context.Background(), bson.M{"_id": tripObjID, "user_id": userID}).Decode(&trip)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "Trip not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to retrieve trip", http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(trip)
}

func GetTrips(w http.ResponseWriter, r *http.Request) {
	// Extract user from JWT token
	tokenString := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(tokenString) // Validate JWT and extract claims
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Convert the claims.UserID string to a primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Fetch all trips for the logged-in user
	cursor, err := db.TripCollection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, "Failed to fetch trips", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and append trips to a slice
	var trips []models.Trip
	for cursor.Next(context.Background()) {
		var trip models.Trip
		if err := cursor.Decode(&trip); err != nil {
			http.Error(w, "Error decoding trip", http.StatusInternalServerError)
			return
		}
		trips = append(trips, trip)
	}

	// Check for any cursor iteration error
	if err := cursor.Err(); err != nil {
		http.Error(w, "Error while fetching trips", http.StatusInternalServerError)
		return
	}

	// Respond with the trips
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

func UpdateTrip(w http.ResponseWriter, r *http.Request) {
    // Get the trip ID from the URL params
    tripID := mux.Vars(r)["id"]
    tripObjID, err := primitive.ObjectIDFromHex(tripID)
    if err != nil {
        http.Error(w, "Invalid trip ID format", http.StatusBadRequest)
        return
    }

    // Decode the incoming request body into a Trip struct
    var trip models.Trip
    err = json.NewDecoder(r.Body).Decode(&trip)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Log the incoming trip object for debugging purposes
    log.Printf("Received trip data: %+v", trip)

    // Extract user ID from the JWT token
    userID, err := getUserIDFromToken(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Log the extracted userID for debugging purposes
    log.Printf("UserID from token: %v", userID)

    // Ensure the trip belongs to the user and user has authorization to update it
    filter := bson.M{"_id": tripObjID, "user_id": userID}

    // Log the filter for debugging purposes
    log.Printf("Filter: %+v", filter)

    // Prepare the update data (don't include user_id in the update)
    updateFields := bson.M{}
    if trip.Name != "" {
        updateFields["name"] = trip.Name
    }
    if trip.Category != "" {
        updateFields["category"] = trip.Category
    }
    if trip.Region != "" {
        updateFields["region"] = trip.Region
    }
    if trip.Description != "" {
        updateFields["description"] = trip.Description
    }
    if trip.Attractions != "" {
        updateFields["attractions"] = trip.Attractions
    }

    // If no fields were specified for update, return an error
    if len(updateFields) == 0 {
        http.Error(w, "No fields to update", http.StatusBadRequest)
        return
    }

    // Perform the update operation, setting only the specified fields
    update := bson.M{"$set": updateFields}
    result := db.TripCollection.FindOneAndUpdate(context.Background(), filter, update)

    // Log the result of the update for debugging purposes
    if result.Err() != nil {
        if result.Err() == mongo.ErrNoDocuments {
            http.Error(w, "Trip not found or you do not have permission to edit", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to update trip", http.StatusInternalServerError)
        }
        return
    }

    // Prepare the updated trip data to send as response (preserve ObjectID)
    // Re-fetch the trip document to send the updated version
    var updatedTrip models.Trip
    err = db.TripCollection.FindOne(context.Background(), filter).Decode(&updatedTrip)
    if err != nil {
        http.Error(w, "Failed to retrieve updated trip", http.StatusInternalServerError)
        return
    }

    // Send the updated trip as a response, preserving the ObjectID
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedTrip)
}



func DeleteTrip(w http.ResponseWriter, r *http.Request) {
    tripID := mux.Vars(r)["id"]
    tripObjID, err := primitive.ObjectIDFromHex(tripID)
    if err != nil {
        http.Error(w, "Invalid trip ID format", http.StatusBadRequest)
        return
    }

    userID, err := getUserIDFromToken(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Delete trip from database ensuring it's the user's trip
    filter := bson.M{"_id": tripObjID, "user_id": userID}
    result, err := db.TripCollection.DeleteOne(context.Background(), filter)
    if err != nil {
        http.Error(w, "Failed to delete trip", http.StatusInternalServerError)
        return
    }

    if result.DeletedCount == 0 {
        http.Error(w, "Trip not found or you do not have permission to delete", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Trip deleted successfully"))
}
