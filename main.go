package main

import (
	"CountriesDashboardService/firebase"
	"CountriesDashboardService/handlers"
	"CountriesDashboardService/utils"
	"log"
	"net/http"
)

// main is the entry point for the application.
// It starts the HTTP server on port 8080 and registers the routes.
func main() {
	log.Println("Starting server...")

	// Initialize Firebase and Firestore
	firebase.InitFirebase()
	defer firebase.CloseFirebase()

	// Start the server
	router := http.NewServeMux()
	port := "8080"

	// Register the routes
	router.HandleFunc("/", handlers.Home)
	router.HandleFunc("/dashboard/v1/registrations/", handlers.HandleRegistrations)
	router.HandleFunc("/dashboard/v1/notifications/", handlers.HandleNotifications)
	router.HandleFunc("/dashboard/v1/dashboards/", handlers.ViewDashboard)
	router.HandleFunc("/dashboard/v1/status/", handlers.HandleStatus)

	//Calling start time
	utils.InitStartTime()

	// Log a message indicating the server has started
	log.Println("Server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
