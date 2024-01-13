package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//Basic route
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	//Route with URL parameters
	router.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(200, "User ID "+id)
	})

	//Route with query parameters
	router.GET("/search", func(c *gin.Context) {
		query := c.DefaultQuery("q", "default_value")
		c.String(200, "Search query : "+query)
	})
	router.Run(":8080")
}
