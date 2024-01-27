package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "user",
		Value:    "admin",
		HttpOnly: true,
		MaxAge:   60,
	}
	http.SetCookie(w, &cookie)
}

func GetCookie(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("user")
	w.Write([]byte(cookie.Value))
	fmt.Println(cookie.Value)
}

func main() {
	r := gin.New()
	r.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 60, "/", "127.0.0.1", false, true)
		}
		fmt.Printf("Cookie value: %s \n", cookie)
	})
	r.Run(":8080")
}
