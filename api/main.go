// Package main implements a server for Greeter service.
package main

import (
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/null-channel/eddington/api/controllers"
)

func main() {
	//TODO: Never used gin. seems like mux is archived. going to try this out.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())
	v1 := router.Group("/v1")
	{
		// Apps
		v1.POST("/apps", controllers.CreateApplication())
		v1.GET("/apps", controllers.GetApplications())
		apps := v1.Group("/apps")
		{
			apps.GET("/:app_id/containers", controllers.CreateContainer())
		}

		// Users

		// AuthZ

		// AuthN

		// Space

	}
	log.Fatal(router.Run(":" + port))

}
