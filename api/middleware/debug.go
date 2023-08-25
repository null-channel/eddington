package middleware

import (
	"context"
	"net/http"
)

type DebugAuth struct {
}

func (app *DebugAuth) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), "user-id", "1234")
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
