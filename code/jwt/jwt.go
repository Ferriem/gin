package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type storage struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Minute * 2

var Secret = []byte("secret")

func GenToken(username string) (string, error) {
	t := storage{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "ferriem",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	return token.SignedString(Secret)
}

func ParseToken(tokenString string) (*storage, error) {
	token, err := jwt.ParseWithClaims(tokenString, &storage{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*storage); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func authHandler(c *gin.Context) {
	user := UserInfo{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "invalid params",
		})
		return
	}
	if user.Username == "root" && user.Password == "123456" {
		tokenString, err := GenToken(user.Username)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": gin.H{"token": tokenString},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 400,
		"msg":  "auth failed",
	})
	return
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "authHeader empty",
			})
			fmt.Println("authHeader empty")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		fmt.Println(parts[0])
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 404,
				"msg":  "authHeader format error",
			})
			c.Abort()
			return
		}
		msg, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "token invalid",
			})
			c.Abort()
			return
		}
		c.Set("username", msg.Username)
		fmt.Println(msg.Username)
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.POST("/auth", authHandler)

	r.GET("/home", JWTAuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		fmt.Println(username)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": username,
		})
	})
	r.Run(":8080")
}
