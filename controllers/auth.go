package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"trip-planner/db"
	"trip-planner/models"
	"trip-planner/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the database with plain text password
	_, err = db.UserCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginDetails struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decode the login request body
	err := json.NewDecoder(r.Body).Decode(&loginDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if email or username exists
	var user models.User
	var filter bson.M

	if loginDetails.Email != "" {
		filter = bson.M{"email": loginDetails.Email}
	} else if loginDetails.Username != "" {
		filter = bson.M{"username": loginDetails.Username}
	} else {
		http.Error(w, "Email or Username must be provided", http.StatusBadRequest)
		return
	}

	// Find the user in the database
	err = db.UserCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare the plaintext password
	if user.Password != loginDetails.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token using MongoDB's ObjectID as the user identifier
	token, err := utils.GenerateJWT(user.ID.Hex()) // Use Hex() to convert ObjectID to string
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the token
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	// Set response header for JSON and send the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

