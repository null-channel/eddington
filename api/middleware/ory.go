package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
	return context.WithValue(ctx, "user-name", v.GetIdentity().Id)
}

func (app *OryApp) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {


		fmt.Println("Authentication Middleware is running")
		log.Printf("handling middleware request\n")

		// set the cookies on the ory client
		var cookies string

		// this example passes all request.Cookies
		// to `ToSession` function
		//
		// However, you can pass only the value of
		// ory_session_projectid cookie to the endpoint
		cookies = request.Header.Get("Cookie")

		// check if we have a session
		session, _, err := app.Ory.FrontendApi.ToSession(request.Context()).Cookie(cookies).Execute()
		if (err != nil && session == nil) || (err == nil && !*session.Active) {
			// this will redirect the user to the managed Ory Login UI
			http.Redirect(writer, request, "/.ory/self-service/login/browser", http.StatusSeeOther)
			return
		}

		ctx := withCookies(request.Context(), cookies)
		ctx = withSession(ctx, session)
		ctx = withUser(ctx,session)

		// continue to the requested page (in our case the Dashboard)
		next.ServeHTTP(writer, request.WithContext(ctx))
		return
	})
}
