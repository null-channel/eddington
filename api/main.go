// Package main implements a server for Greeter service.
package main

import (
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/null-channel/eddington/api/controllers"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port := "8080"
		router := gin.New()
		router.Use(gin.Logger())
		router.POST("/applications", controllers.CreateApplication())

		log.Fatal(router.Run(":" + port))

	}
}
