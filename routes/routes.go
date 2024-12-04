package routes

import (
	"trip-planner/controllers"

	"github.com/gorilla/mux"
)

// InitializeRoutes configures all application routes
func InitializeRoutes() *mux.Router {
	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", controllers.LoginUser).Methods("POST")

	// Trip routes
	r.HandleFunc("/trips", controllers.CreateTrip).Methods("POST")                      // Create a new trip
	r.HandleFunc("/trips/{id}", controllers.GetTripByID).Methods("GET")                 // Get trip by ID
	r.HandleFunc("/trips", controllers.GetTrips).Methods("GET")                         // Get all trips
	r.HandleFunc("/trips/{id}", controllers.UpdateTrip).Methods("PUT")                  // Update an existing trip
	r.HandleFunc("/trips/{id}", controllers.DeleteTrip).Methods("DELETE")               // Delete a trip

	// Comment routes
	r.HandleFunc("/comments/{trip_id}/comments", controllers.CreateComment).Methods("POST")   // Create comment
	r.HandleFunc("/comments/{trip_id}/comments", controllers.GetComments).Methods("GET")    // Get all comments for a specific trip
	r.HandleFunc("/comments/{trip_id}/comments/{id}", controllers.GetCommentByID).Methods("GET") // Get comment by ID
	r.HandleFunc("/comments/{trip_id}/comments/{id}", controllers.UpdateComment).Methods("PUT")   // Update comment
	r.HandleFunc("/comments/{trip_id}/comments/{id}", controllers.DeleteComment).Methods("DELETE") // Delete comment

	return r
}
