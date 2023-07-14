package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
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

func (o *OryController) AddOryRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/webhook", o.OryWebhook())
}

func (o *OryController) OryWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("New user registered from Ory!")
		forwardedFor := c.Request.Header.Get("X-Forwarded-For")

		// Check if request is from Ory
		if _, ok := oryIps[forwardedFor]; !ok {
			c.AbortWithStatus(403)
			return
		}

		var user types.CreateUserRequest

		bytes := []byte{}

		c.Request.Body.Read(bytes)

		json.Unmarshal(bytes, &user)

		fmt.Println(user)

		userDB := models.CreateUserRequestToDBModel(user)
		code, err := o.userController.CreateUser(userDB)

		if err != nil {
			c.JSON(code, gin.H{"error": err.Error()})
		}

		c.JSON(code, gin.H{"status": "User created successfully!"})
	}
}
