package controllers

import "github.com/gin-gonic/gin"

// OK sends a standardized success JSON response wrapping the payload in a "data" key.
func OK(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, gin.H{
		"data": data,
	})
}

// Fail sends a standardized error JSON response with an error message.
func Fail(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}
