package controllers

import (
	"github.com/gin-gonic/gin"
)

type struct UserController struct{}

func (uc *UserController) GetFirst(client *redis.Client, c *gin.Context) {
	userId := c.Param("id")


}