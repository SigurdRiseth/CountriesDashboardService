package main

import (
	"CountriesDashboardService/consts"
	"CountriesDashboardService/firebase"
	"CountriesDashboardService/handlers"
	"CountriesDashboardService/utils"
	"log"
	"net/http"
)

// main is the entry point for the application.
// It starts the HTTP server on port 8080 and registers the routes.
func main() {
	log.Println(consts.LogStartingServer)

	// Initialize Firebase and Firestore
	firebase.InitFirebase()
	defer firebase.CloseFirebase()

	// Start the server
	router := http.NewServeMux()
	port := consts.PortInUse

	// Register the routes
	router.HandleFunc(consts.Dash, handlers.Home)
	router.HandleFunc(consts.BaseRegistrationPath, handlers.HandleRegistrations)
	router.HandleFunc(consts.BaseNotificationPath, handlers.HandleNotifications)
	router.HandleFunc(consts.BaseDashboardPath, handlers.ViewDashboard)
	router.HandleFunc(consts.BaseDashboardPath, handlers.HandleStatus)

	//Calling start time
	utils.InitStartTime()

	// Log a message indicating the server has started
	log.Println(consts.LogServerStarted + port)
	log.Fatal(http.ListenAndServe(consts.LogColon+port, router))
}
