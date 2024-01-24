package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//Create a new Gin router
	router := gin.Default()

	//define a router for the root URL
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	router.GET("/bye/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Goodbye World",
		})
	})

	//Run the router on port 8080
	router.Run(":8080")
}
