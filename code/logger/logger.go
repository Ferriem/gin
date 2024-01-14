package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request - Method: %s | Status: %d | Duration: %v\n", c.Request.Method, c.Writer.Status(), duration)
	}
}

func main() {
	//Create a new Gin router
	router := gin.Default()

	router.Use(LoggerMiddleware())

	//define a router for the root URL
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	router.GET("/bye/", func(c *gin.Context) {
		c.String(200, "Goodbye World")
	})

	//Run the router on port 8080
	router.Run(":8080")
}
