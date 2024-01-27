package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//http
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.google.com")
	})

	//router

	router.GET("/first", func(c *gin.Context) {
		c.Request.URL.Path = "/second"
		router.HandleContext(c)
	})

	router.GET("/second", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "second"})
	})

	router.Run(":8080")
}
