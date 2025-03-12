package main

import (
	"log"
	"net/http"
)

func main() {
	// Start the server
	router := http.NewServeMux()
	port := "8080"

	// Log a message indicating the server has started
	log.Println("Server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
