package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	ory "github.com/ory/client-go"
)

type OryApp struct {
	Ory *ory.APIClient
}

// save the cookies for any upstream calls to the Ory apis
func withCookies(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, "req.cookies", v)
}

func getCookies(ctx context.Context) string {
	return ctx.Value("req.cookies").(string)
}

// save the session to display it on the dashboard
func withSession(ctx context.Context, v *ory.Session) context.Context {
	return context.WithValue(ctx, "req.session", v)
}

func getSession(ctx context.Context) *ory.Session {
	return ctx.Value("req.session").(*ory.Session)
}

func withUser(ctx context.Context, v *ory.Session) context.Context {
	return context.WithValue(ctx, "user-id", v.GetIdentity().Id)
}

func (app *OryApp) SessionMiddleware(next http.Handler) http.Handler {
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
		user_id, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			fmt.Println(token)
			return token.Claims.(jwt.MapClaims)["user-id"], nil
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		//TODO: Delete this line
		fmt.Println(user_id)

		//ctx = withSession(ctx, session)
		//ctx = withUser(ctx, session)

		// continue to the requested page (in our case the Dashboard)
		next.ServeHTTP(writer, request.WithContext(ctx))
		return
	})
}
