package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func CreateCORSHandler(router *mux.Router) http.Handler {

	fmt.Println("Cors Sucks")
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)
}
