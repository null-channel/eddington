package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	user "github.com/null-channel/eddington/api/users/models"
	"github.com/uptrace/bun"
)

// NewUserMiddleware is a middleware that checks if the user is new.
type AuthzMiddleware struct {
	db *bun.DB
}

// NewAuthzMiddleware creates a new user middleware.
func NewAuthzMiddleware(db *bun.DB) *AuthzMiddleware {
	return &AuthzMiddleware{db: db}
}

func (k *AuthzMiddleware) CheckAuthz(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user-id header is set
		userId, ok := r.Context().Value("user-id").(string)
		if !ok {
			fmt.Println("User is new, redirecting to new user page")
			http.Error(w, "User is new, redirecting to new user page", http.StatusTemporaryRedirect)
			w.Header().Set("location", "/newuser")
			return
		}
		fmt.Println("Checking if user is new...")
		// Check database for user
		_, err := user.GetUserForId(userId, k.db)

		if err != nil {
			fmt.Println("User is new, redirecting to new user page")
			http.Error(w, "User is new, redirecting to new user page", http.StatusTemporaryRedirect)
			w.Header().Set("location", "/newuser")
			return
		}
		org, err := user.GetOrgByOwnerId(userId, k.db)

		if err != nil {
			fmt.Println("Org not found for user. Failing.")
			http.Redirect(w, r, "error", http.StatusSeeOther)
			w.Header().Set("location", "error")
		}

		// check url for org id
		vars := mux.Vars(r)
		orgId := vars["orgId"]

		if err != nil {
			fmt.Println("Org id not found in url. Failing Authz.")
			http.Redirect(w, r, "/error", http.StatusSeeOther)
		}

		// check if org id matches org id in url
		if orgId != fmt.Sprintf("%d", org.ID) {
			fmt.Println("Org id in url does not match org id in database. Failing Authz.")
			http.Redirect(w, r, "/error", http.StatusSeeOther)
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
