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
		// this example passes all request.Cookies
		// to `ToSession` function
		//
		// However, you can pass only the value of
		// ory_session_projectid cookie to the endpoint
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
		userId, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//fmt.Println("claims: " + token.Claims.(jwt.MapClaims)["sub"])
			claims := token.Claims.(jwt.MapClaims)
			// You can now extract any data from the token's payload
			return claims["sub"], nil
		})
		user_id := fmt.Sprintf("%v", userId)
		if err != nil {
			fmt.Println("Error parsing token! but that is ok")
			// can fail if the token is invalid but we don't want to validate it here for now
			//return
		}
		//TODO: Delete this line
		fmt.Println("request userId: %s" + user_id)
		fmt.Println("next line")
		ctx = withUser(ctx, user_id)

		//ctx = withSession(ctx, session)
		request.Header.Set("user-id", fmt.Sprintf("%v", user_id))

		// continue to the requested page (in our case the Dashboard)
		next.ServeHTTP(writer, request.WithContext(ctx))
		return
	})
}
