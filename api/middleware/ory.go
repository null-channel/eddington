package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (app *OryApp) SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("handling session middleware for request\n")

		// set the cookies on the ory client
		var cookies string

		// this example passes all request.Cookies
		// to `ToSession` function
		// However, you can pass only the value of
		// ory_session_projectid cookie to the endpoint
		cookies = c.Request.Header.Get("Cookie")

		fmt.Println("cookies: ", cookies)

		// check if we have a session
		session, _, err := app.Ory.FrontendApi.ToSession(c.Request.Context()).Cookie(cookies).Execute()
		if (err != nil && session == nil) || (err == nil && !*session.Active) {
			// if we don't have a session, we need to fail the request
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Not Authorized")
			return
		}

		ctx := withCookies(c.Request.Context(), cookies)
		_ = withSession(ctx, session)

		// continue to the requested page
		c.Next()
	}
}
