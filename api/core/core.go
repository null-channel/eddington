package core

import (
	"encoding/json"
	"net/http"

	"github.com/null-channel/eddington/api/users/types"
)

type Error struct {
	Code    int
	Message string
}

func writeError(writer http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	json.NewEncoder(writer).Encode(resp)
}
func writeValidationError(writer http.ResponseWriter, apiError []types.ApiError) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(apiError)
}

var (
	ValidationErrors = func(w http.ResponseWriter, apiError []types.ApiError) {
		writeValidationError(w, apiError)
	}
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An Unexpected Error Occured", http.StatusInternalServerError)
	}
	UnauthorizedErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
	}
	UserRegistrationError = func(w http.ResponseWriter) {
		writeError(w, "Please complete the registration process.", http.StatusNotFound)
	}
)
