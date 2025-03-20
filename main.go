package main

import (
	"log"
	"net/http"
)

// main is the entry point for the application.
// It starts the HTTP server on port 8080 and registers the routes.
func main() {
	log.Println("Starting server...")

	// Start the server
	router := http.NewServeMux()
	port := "8080"

	// Register the routes
	router.HandleFunc("/", home)

	// Log a message indicating the server has started
	log.Println("Server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// home handles HTTP requests to the root path.
// It supports only the GET method and responds with "Hello, World!".
// For unsupported methods, it responds with a 405 Method Not Allowed status.
func home(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Hello, World!"))
		return
	default:
		log.Println("Unsupported request method " + request.Method)
		http.Error(writer, "Unsupported request method "+request.Method, http.StatusMethodNotAllowed)
		return
	}
}
