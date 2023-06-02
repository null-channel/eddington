package controllers

import (
	"github.com/gin-gonic/gin"
)

func CreateApplication() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(501, "Not implemented yet")
	}
}
func GetApplications() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(501, "Not implemented yet")
	}
}
