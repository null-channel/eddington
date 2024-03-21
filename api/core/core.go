package core

import (
	"encoding/json"
	"net/http"
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

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An Unexpected Error Occured", http.StatusInternalServerError)
	}
	UnauthorizedErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
	}
)
