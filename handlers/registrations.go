package handlers

import (
	"log"
	"net/http"
)

// HandleRegistrations handles HTTP requests for dashboard configurations.
// It supports the following methods:
// - POST: Adds a new dashboard configuration.
// - GET: Views all dashboard configurations or a specific one if an {id} parameter is provided.
// - HEAD: Checks if a specific dashboard configuration exists.
// - PUT: Replaces the existing dashboard configuration.
// - DELETE: Deletes the current dashboard configuration.
// For unsupported methods, it responds with a 405 Method Not Allowed status.
func HandleRegistrations(writer http.ResponseWriter, request *http.Request) {
	log.Println("Received " + request.Method + " request.")
	switch request.Method {
	case http.MethodPost:
		addDashboardConfiguration(writer, request)
	case http.MethodGet:
		viewDashboardConfiguration(writer, request)
	case http.MethodHead:
		handleHeadRequest(writer, request)
	case http.MethodPut:
		replaceDashboardConfiguration(writer, request)
	case http.MethodDelete:
		deleteDashboardConfiguration(writer, request)
	default:
		log.Printf("Unsupported request method: %s", request.Method)
		http.Error(writer, "Unsupported request method "+request.Method, http.StatusMethodNotAllowed)
		return
	}
}

// TODO: Implement the addDashboardConfiguration function. Issue #5
func addDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("POST method not implemented yet"))
}

// TODO: Implement the viewDashboardConfiguration function. Issue #6-7
func viewDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("GET method not implemented yet"))
}

// TODO: Implement the replaceDashboardConfiguration function. Issue #8
func replaceDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("PUT method not implemented yet"))
}

// TODO: Implement the deleteDashboardConfiguration function. Issue #9
func deleteDashboardConfiguration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("DELETE method not implemented yet"))
}

// TODO: !!ADVANCED TASK!! Implement the handleHeadRequest function. Issue #10
func handleHeadRequest(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	writer.Write([]byte("HEAD method not implemented yet"))
}
