package middleware

import "net/http"

type AuthMiddleware interface {
	SessionMiddleware(http.Handler) http.Handler
}

