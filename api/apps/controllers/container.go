package controllers

import (
	"github.com/gin-gonic/gin"
)

func ContainerPOST() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(501, "Not implemented yet")
	}
}
