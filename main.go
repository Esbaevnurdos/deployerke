package main

import (
	"log"
	"net/http"
	"trip-planner/db"
	"trip-planner/routes"
)

func main() {
	// Initialize DB
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize routes
	r := routes.InitializeRoutes()

	// Start server
	log.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", r)
}
