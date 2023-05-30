package controllers

import (
	"github.com/gin-gonic/gin"
)

func CreateApplication() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.IndentedJSON(200, "application created successfully")

	}

}
