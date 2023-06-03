// Package main implements a server for Greeter service.
package main

import (
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/null-channel/eddington/api/controllers"
	"github.com/null-channel/eddington/api/docs"
	userController "github.com/null-channel/eddington/api/users/controllers"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"
)

func main() {
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
		v1.POST("/apps", controllers.CreateApplication())
		v1.GET("/apps", controllers.GetApplications())
		apps := v1.Group("/apps")
		{
			apps.GET("/:app_id/containers", controllers.CreateContainer())
		}

		// Users
		userController.AddAllControllers(v1)

		// AuthZ

		// AuthN

		// Space

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	log.Fatal(router.Run(":" + port))

}
