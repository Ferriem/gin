package main

import (
	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Key string `json:"key"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		} else {
			c.Next()
		}
	}
}

func main() {
	router := gin.Default()

	public := router.Group("/public")
	public.GET("/info", func(c *gin.Context) {
		c.String(200, "Public information")
	})
	public.GET("/products", func(c *gin.Context) {
		c.String(200, "Public products")
	})

	private := router.Group("/private")
	private.Use(AuthMiddleware())
	private.GET("/data", func(c *gin.Context) {
		c.String(200, "Private data accessible after authentication")
	})
	private.POST("/create", func(c *gin.Context) {
		var requestBody RequestBody

		// Bind JSON data from the request body to the struct
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.String(400, "Bad Request: "+err.Error())
			return
		}

		// Access the data from the struct
		key := requestBody.Key
		c.String(200, "Key: "+key)
	})

	router.Run(":8080")

}
