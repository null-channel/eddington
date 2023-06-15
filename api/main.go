// Package main implements a server for Greeter service.
package main

import (
	"fmt"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/null-channel/eddington/api/controllers"
	"github.com/null-channel/eddington/api/docs"
	marketing "github.com/null-channel/eddington/api/marketing/controllers"
	userController "github.com/null-channel/eddington/api/users/controllers"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"
)

func main() {

	from := mail.NewEmail("Marek Counts", "mcounts@nulllabs.io")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "mcounts@nulllabs.io")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	//TODO: Never used gin. seems like mux is archived. going to try this out.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())
	docs.SwaggerInfo.BasePath = "/api/v1"
	userController := userController.New()

	v1 := router.Group("api/v1")
	{
		// Apps
		v1.POST("/apps", controllers.AppPOST())
		v1.GET("/apps", controllers.AppGET())
		apps := v1.Group("/apps")
		{
			apps.GET("/:app_id/containers", controllers.ContainerPOST())
		}

		// Users
		userController.AddAllControllers(v1)

		// AuthZ

		// AuthN

		// Space

		// Marketing
		marketingGroup := v1.Group("/marketing")
		{
			_ = marketing.New(os.Getenv("SENDGRID_API_KEY"), marketingGroup)
		}

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	log.Fatal(router.Run(":" + port))

}
