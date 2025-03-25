package handlers

import "net/http"

// TODO: Implement the ViewDashboard function us get method from registrations.
func ViewDashboard(writer http.ResponseWriter, request *http.Request) {
	// Set the response header to indicate the content type is JSON
	writer.Header().Set("Content-Type", "application/json")
	// Set the response status code to 200 OK
	writer.WriteHeader(http.StatusOK)
	// Write a JSON response with a message indicating the dashboard is under construction
	writer.Write([]byte(`{"message": "Dashboard is under construction"}`))
}
