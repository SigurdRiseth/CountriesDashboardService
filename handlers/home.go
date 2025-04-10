package handlers

import (
	"CountriesDashboardService/consts"
	"log"
	"net/http"
)

// Home handles HTTP requests to the root path.
// It supports only the GET method and responds with "Hello, World!".
// For unsupported methods, it responds with a 405 Method Not Allowed status.
//func Home(writer http.ResponseWriter, request *http.Request) {
//	switch request.Method {
//	case http.MethodGet:
//		writer.WriteHeader(http.StatusOK)
//		writer.Write([]byte("Hello, World!"))
//		return
//	default:
//		log.Println("Unsupported request method " + request.Method)
//		http.Error(writer, "Unsupported request method "+request.Method, http.StatusMethodNotAllowed)
//		return
//	}
//}

// HomeHandler serves the static API documentation page located in /static/index.html.
func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		http.ServeFile(writer, request, consts.StaticFilePath)
	default:
		log.Println(consts.UnsopprtedReqMethod + request.Method)
		http.Error(writer, consts.UnsopprtedReqMethod+request.Method, http.StatusMethodNotAllowed)
	}
}
