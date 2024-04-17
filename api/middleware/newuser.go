package middleware

import (
	"fmt"
	"net/http"

	"github.com/null-channel/eddington/api/users/controllers"
)

// NewUserMiddleware is a middleware that checks if the user is new.
type UserMiddleware struct {
	membersDatastore controllers.MembersDatastore
}

// NewUserMiddleware creates a new user middleware.
func NewUserMiddleware(members controllers.MembersDatastore) *UserMiddleware {
	return &UserMiddleware{membersDatastore: members}
}

func (k *UserMiddleware) NewUserMiddlewareCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, userExists := r.Context().Value("user-id").(string)
		if !userExists {
			http.Error(w, "User not found in context", http.StatusBadRequest)
			return
		}
		fmt.Println("Checking if user is new: ", userId)
		// Check database for user
		_, err := k.membersDatastore.GetUserByID(r.Context(), userId)

		if err != nil {
			http.Error(w, "User is new, redirecting to new user page", http.StatusBadRequest)
			w.Header().Set("location", "/newuser")
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)

	})
}
