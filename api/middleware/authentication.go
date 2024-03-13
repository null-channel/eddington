package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// save the cookies for any upstream calls to the Ory apis
func withCookies(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, "req.cookies", v)
}
func withUser(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, "user-id", v)
}

func getCookies(ctx context.Context) string {
	return ctx.Value("req.cookies").(string)
}

func AddJwtHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		//TODO: Parse JWT token and get user id
		fmt.Println("Authentication Middleware is running")
		log.Printf("handling middleware request\n")

		// set the cookies on the ory client
		var cookies string

		ctx := withCookies(request.Context(), cookies)
		cookies = request.Header.Get("Cookie")
		tokenString := request.Header.Get("Authorization")

		// remove the Bearer prefix
		// and parse the token
		parser := &jwt.Parser{
			ValidMethods:         []string{"none"},
			UseJSONNumber:        true,
			SkipClaimsValidation: true,
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		var claims jwt.MapClaims
		_, _, err := parser.ParseUnverified(tokenString, &claims)
		if err != nil {
			fmt.Println("Error parsing token! but that is ok")
			fmt.Println(err)
			return
		}
		userId := claims["sub"].(string)

		ctx = withUser(ctx, userId)

		request.Header.Set("user-id", userId)

		// continue to the requested page (in our case the Dashboard)
		next.ServeHTTP(writer, request.WithContext(ctx))
		return
	})
}
