package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	core "github.com/null-channel/eddington/api/core"
)

func withUser(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, "user-id", v)
}

func AddJwtHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		tokenString := request.Header.Get("Authorization")

		if tokenString == "" {
			fmt.Println("Missing token in header")
			core.UnauthorizedErrorHandler(writer)
			return
		}
		// remove the Bearer prefix
		parser := &jwt.Parser{
			ValidMethods:         []string{"none"},
			UseJSONNumber:        true,
			SkipClaimsValidation: true,
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		if tokenString == "" {
			fmt.Println("Missing token in header")
			core.UnauthorizedErrorHandler(writer)
			return
		}
		var claims jwt.MapClaims
		// parse the token
		_, _, err := parser.ParseUnverified(tokenString, &claims)
		if err != nil {
			fmt.Println("Error parsing token! ")
			fmt.Println(err)
			core.UnauthorizedErrorHandler(writer)
			return
		}
		userId := claims["sub"].(string)

		ctx := withUser(request.Context(), userId)

		request.Header.Set("user-id", userId)

		// continue to the requested page (in our case the Dashboard)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
