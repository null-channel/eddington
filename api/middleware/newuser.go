package middleware

import (
	"fmt"
	"net/http"

	user "github.com/null-channel/eddington/api/users/models"
	"github.com/uptrace/bun"
)

// NewUserMiddleware is a middleware that checks if the user is new.
type UserMiddleware struct {
	db *bun.DB
}

// NewUserMiddleware creates a new user middleware.
func NewUserMiddleware(db *bun.DB) *UserMiddleware {
	return &UserMiddleware{db: db}
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
		_, err := user.GetUserForId(userId, k.db)

		if err != nil {
			http.Error(w, "User is new, redirecting to new user page", http.StatusBadRequest)
			w.Header().Set("location", "/newuser")
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)

	})
}
