package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	usercontroller "github.com/null-channel/eddington/api/users/controllers"
	"github.com/null-channel/eddington/api/users/models"
	"github.com/null-channel/eddington/api/users/types"
)

var oryIps = map[string]string{
	"34.22.170.75":   "",
	"35.242.228.133": "",
}

func NewOryController(userController *usercontroller.UserController) OryController {
	return OryController{userController: userController}
}

type OryController struct {
	userController *usercontroller.UserController
}

func (o *OryController) AddOryRoutes(router *mux.Router) {
	router.HandleFunc("/webhook", o.OryWebhook).Methods("POST")
}

func (o *OryController) OryWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New user registered from Ory!")
	forwardedFor := r.Header.Get("X-Forwarded-For")

	// Check if request is from Ory
	if _, ok := oryIps[forwardedFor]; !ok {
		http.Error(w, "You are not Ory.", http.StatusBadRequest)
		return
	}

	var user types.CreateUserRequest

	bytes := []byte{}

	r.Body.Read(bytes)

	json.Unmarshal(bytes, &user)

	fmt.Println(user)

	userDB := models.CreateUserRequestToDBModel(user)
	code, err := o.userController.CreateUser(userDB)

	if err != nil {
		http.Error(w, "Internal server error", code)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status OK, good jub")
}
