package handlers

import (
	"net/http"
)

func HandleNotifications(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		registerWebhook(writer, request)
	case http.MethodDelete:
		deleteWebhook(writer, request)
	case http.MethodGet:
		retrieveWebhook(writer, request)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func registerWebhook(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Service not implemented"))
}

func deleteWebhook(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Service not implemented"))
}

func retrieveWebhook(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Service not implemented"))
}
