package main

import (
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	userId := c.Param("id")

	c.JSON(200, gin.H{"id": userId, "name": "John Doe", "email": "john@example.com"})
}

func main() {
	router := gin.Default()

	userController := &UserController{}

	router.GET("/user/:id", userController.GetUserInfo)

	router.Run(":8080")

}
