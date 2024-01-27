# Gin

[website](https://masteringbackend.com/posts/gin-framework#getting-started-with-gin)

## Getting Started with Gin

```go
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

	//Run the router on port 8080
	router.Run(":8080")
}
```

Created a new Gin router using `gin.Default()`. Then defined a simple router for the root URL ("/") that responds with "Hello, World!".

## Framework

### API

#### API parameters

The parameters can be passed like this `user/ferriem/score`

```go
	r.GET("/user/:name/:title", func(c *gin.Context) {
		name := c.Param("name") //name = ferriem
		title := c.Param("title") // title = score
		c.JSON(http.StatusOK,gin.H{
			"name":name, 
			"title":title,
		})
	})
```

#### Query parameters

`/user?name=ferriem&tag=score`

```go
	r.GET("/user", func(c *gin.Context) {
		name := c.DefaultQuery("name","Rick")
		tag := c.Query("tag")
		c.JSON(http.StatusOK,gin.H{
			"name":name,
			"tag":tag,
		})
	})
```

### Middleware in Gin

#### Logger middleware

```go
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
```

- The `LoggerMiddleware` function returns a `gin.HandlerFunc`. This is the signature expected by Gin for middleware functions.
- `c.Next()`: It tells the Gin framework to move to the next middleware or the stack or the final route handler. Here `router.Get()`
- The `LoggerMiddleware` is added to the Gin router using `router.Use()`. This means that every request going through this Gin router will first pass through the `LoggerMiddleware` before reaching the actual router handler.

When a request hits the server, it first goes through the `LoggerMiddleware`, which logs information about this request. After that, `c.Next()` allows the request to proceed to the defined route handlers.

#### Creating Custom Middleware

Custom middleware can handle tasks like authentication, data validation, rate limiting, and more.

```go
package main

import (
	"github.com/gin-gonic/gin"
)

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
	authGroup := router.Group("/api")
	authGroup.Use(AuthMiddleware())
	{
		authGroup.GET("/data", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Authenticated and authorized"})
		})
	}

	router.Run(":8080")
}
```

```sh
~/ curl 127.0.0.1:8080/api/data
{"error":"unauthorized"}%
~/ curl -H "X-API-KEY: ferriem" 127.0.0.1:8080/api/data
{"message":"Authenticated and authorized"}%
```

- `Middleware Definition` check for the presence of the "X-API-KEY" header in the incoming request.
- `router` A router Group created under the "/api" path, 
  - Group is a subset of routes.
  - The `AuthMiddleware` is added to the `authGroup`. This means every request to endpoints under "/api" will go through the authentication middleware.
- `api/data` A router handler is defined for the "/api/data" endpoint.

### Routing and Grouping

Routing is mapping incoming HTTP requests to specific route handles. The router mathces the URL path and HTTP method of the request to find the appropriate handler to execute.

#### Basic Routic

```go
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
```

```sh
~/ curl 127.0.0.1:8080
Hello World%
~/ curl 127.0.0.1:8080/user/ferriem
User ID ferriem%
~/ curl 127.0.0.1:8080/search
Search query : default_value%
~/ curl "127.0.0.1:8080/search?q=WhoAmI"
Search query : WhoAmI%
```

- `/user/:id` define a route with a URL parameter(`:id`). When a GET request is made to "/user/{some_id}", it captures the value of "some_id" and responds.
- `/search` define a route with query parameters. 

#### Route Groups

```go
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
```

```sh
~/ curl 127.0.0.1:8080/public/products
Public products%
~/ curl 127.0.0.1:8080/public/info
Public information%
~/ curl -H "X-API-KEY: ferriem" 127.0.0.1:8080/private/data
Private data accessible after authentication%
~/ curl -X POST -H "X-API-KEY: ferriem" -d '{"key": "value"}' 127.0.0.1:8080/private/create
Key: value%
```

### Controller and Handlers

```go
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
```

```sh
~/ curl 127.0.0.1:8080/user/ferriem
{"email":"john@example.com","id":"ferriem","name":"John Doe"}%
```

In this example, we created a `UserController` struct with a `GetUserInfo` method to handle user-related logic.

Separating business logic into controllers makes the codebase cleaner and more organized.

### Building with Gin

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	router := gin.Default()

	//Connects to an SQLite database and initializes table for Todo model
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("\nReceived termination signal. Cleaning up...")
		// Close the database connection
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println("Error getting *sql.DB:", err)
			os.Exit(1)
		}

		// Close the *sql.DB connection
		err = sqlDB.Close()
		if err != nil {
			fmt.Println("Error closing *sql.DB:", err)
		} else {
			fmt.Println("*sql.DB closed successfully.")
		}

		// Delete the todo.db file
		err = os.Remove("todo.db")
		if err != nil {
			fmt.Println("Error deleting todo.db:", err)
		} else {
			fmt.Println("todo.db deleted successfully.")
		}

		// Exit the program
		os.Exit(0)
	}()

	router.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}
		if todo.Title == "" || todo.Description == "" {
			c.JSON(400, gin.H{"error": "Title and description cannot be empty"})
		}

		db.Create(&todo)

		c.JSON(200, todo)
	})

	router.GET("/todos", func(c *gin.Context) {
		var todos []Todo
		db.Find(&todos)
		if len(todos) == 0 {
			c.JSON(404, gin.H{"error": "todos is empty"})
			return
		}
		c.JSON(200, todos)
	})

	router.GET("/todos/:id", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")
		result := db.First(&todo, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "No todo found with that ID"})
			} else {
				c.JSON(404, gin.H{"error": "Todo not found"})
			}
			return
		}
		c.JSON(200, todo)
	})

	router.PUT("/todos/:id", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")
		result := db.First(&todo, id)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}
		var updatedTodo Todo
		if err := c.ShouldBindJSON(&updatedTodo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}
		if updatedTodo.Title == "" || updatedTodo.Description == "" {
			c.JSON(400, gin.H{"error": "Title and description cannot be empty"})
		}

		todo.Title = updatedTodo.Title
		todo.Description = updatedTodo.Description
		db.Save(&todo)

		c.JSON(200, todo)
	})

	router.DELETE("/todos/:id", func(c *gin.Context) {
		var todo Todo
		id := c.Param("id")

		result := db.First(&todo, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "No todo found with that ID"})
			} else {
				c.JSON(404, gin.H{"error": "Todo not found"})
			}
			return
		}

		db.Delete(&todo)

		c.JSON(200, gin.H{"message": fmt.Sprintf("Todo with ID %s deleted", id)})
	})

	router.Run(":8080")
}
```

```sh
~/ curl -X POST -d '{"title": "key", "description": "value"}' 127.0.0.1:8080/todos
{"ID":1,"CreatedAt":"2024-01-14T18:49:22.459819+08:00","UpdatedAt":"2024-01-14T18:49:22.459819+08:00","DeletedAt":null,"title":"key","description":"value"}%
~/ curl -X GET 127.0.0.1:8080/todos
[{"ID":1,"CreatedAt":"2024-01-14T18:49:22.459819+08:00","UpdatedAt":"2024-01-14T18:49:22.459819+08:00","DeletedAt":null,"title":"key","description":"value"}]%
~/ curl -X POST -d '{"title": "apple", "description": "tech"}' 127.0.0.1:8080/todos
{"ID":2,"CreatedAt":"2024-01-14T18:50:23.856003+08:00","UpdatedAt":"2024-01-14T18:50:23.856003+08:00","DeletedAt":null,"title":"apple","description":"tech"}%
~/ curl -X GET 127.0.0.1:8080/todos
[{"ID":1,"CreatedAt":"2024-01-14T18:49:22.459819+08:00","UpdatedAt":"2024-01-14T18:49:22.459819+08:00","DeletedAt":null,"title":"key","description":"value"},{"ID":2,"CreatedAt":"2024-01-14T18:50:23.856003+08:00","UpdatedAt":"2024-01-14T18:50:23.856003+08:00","DeletedAt":null,"title":"apple","description":"tech"}]%
~/ curl -X GET 127.0.0.1:8080/todos/1
{"ID":1,"CreatedAt":"2024-01-14T18:49:22.459819+08:00","UpdatedAt":"2024-01-14T18:49:22.459819+08:00","DeletedAt":null,"title":"key","description":"value"}%
~/ curl -X PUT -d '{"title": "key_update", "description": "value_update"}' 127.0.0.1:8080/todos/2
{"ID":2,"CreatedAt":"2024-01-14T18:50:23.856003+08:00","UpdatedAt":"2024-01-14T18:53:58.312751+08:00","DeletedAt":null,"title":"key_update","description":"value_update"}%
~/ curl -X GET 127.0.0.1:8080/todos/2
{"ID":2,"CreatedAt":"2024-01-14T18:50:23.856003+08:00","UpdatedAt":"2024-01-14T18:53:58.312751+08:00","DeletedAt":null,"title":"key_update","description":"value_update"}%
~/ curl -X DELETE 127.0.0.1:8080/todos/1
{"message":"Todo with ID 1 deleted"}%
~/ curl -X GET 127.0.0.1:8080/todos
[{"ID":2,"CreatedAt":"2024-01-14T18:50:23.856003+08:00","UpdatedAt":"2024-01-14T18:53:58.312751+08:00","DeletedAt":null,"title":"key_update","description":"value_update"}]%
```

- `AutoMigrate` create **database tables** based on the defined Go struct.

- `ShouldBindJSON` to find whether a request with a data segment. Addtionally, if a filed in the JSON data **matches an exported filed** in the structure, Gin attempts to bind the value to that field. **Ignore** extra filed. 

- `db.Create` insert a new record into the database table that correspinds to the GORM model.
- `db.Save` Is used for **both** creating a new record and updating an existing record. Change the update time.
- `db.Find` retrieve records from the database match the conditions specified by the provided struct or map.

- `db.First` retrieve the first record that matches the conditions. If there is nothing  matched, return error.

### Data handling

```go
type Login struct {
	User string `form:"user" json:"user" uri:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" binding:"required"`
}


func main() {
	r := gin.Default()
	login := Login{}
	//{"user":"ferriem","password":"123456"}
	r.POST("/json", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&login); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"user": login.User,
			"password": login.Password,
		})
	})


	// form
	r.POST("/form", func(c *gin.Context) {
		if err := c.ShouldBind(&login);err != nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"user":login.User,
			"password":login.Password,
		})
	})

	// query
	r.GET("/query", func(c *gin.Context) {
		if err := c.ShouldBindQuery(&login); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"user":login.User,
			"password":login.Password,
		})
	})

	// api
	r.GET("/api/:user/:password", func(c *gin.Context) {
		if err := c.ShouldBindUri(&login); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"user":login.User,
			"password":login.Password,
		})
	})

	r.Run()
}
```

### Log And Logrus

#### Configuration

```sh
go get -u github.com/sirupsen/logrus
```

Implement a server to log request into local Redis database.

```go
package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisHook struct {
	client *redis.Client
	key    string
}

func NewRedisHook(client *redis.Client, key string) *RedisHook {
	return &RedisHook{
		client: client,
		key:    key,
	}
}

func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	b := entry.Data

	logData, err := json.Marshal(b)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = hook.client.RPush(ctx, hook.key, logData).Err()
	if err != nil {
		return err
	}
	return nil
}

func (hook *RedisHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func LoggerToRedis(rdb *redis.Client) gin.HandlerFunc {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	redisHook := NewRedisHook(rdb, "logrus")

	logger.AddHook(redisHook)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	router := gin.Default()
	router.Use(LoggerToRedis(redisClient))

	router.Any("/product", func(c *gin.Context) {
		price := c.Query("price")
		if price != "" {
			c.JSON(200, gin.H{"message": "product", "price": price})
			return
		}
		c.JSON(200, gin.H{"message": "product"})
	})

	router.Any("/order", func(c *gin.Context) {
		c.JSON(400, gin.H{"error": "order"})
	})

	router.Run(":8080")

}

```

```sh
~/ curl -X GET "127.0.0.1:8080/product?price=10"
{"message":"product","price":"10"}%
~/ curl -X DELETE "127.0.0.1:8080/product?price=10"
{"message":"product","price":"10"}%
~/ curl -X DELETE "127.0.0.1:8080/product?price"
{"message":"product"}%

127.0.0.1:6379> lrange logrus 0 -1
1) "{\"client_ip\":\"127.0.0.1\",\"latency_time\":80625,\"req_method\":\"GET\",\"req_uri\":\"/product?price=10\",\"status_code\":200}"
2) "{\"client_ip\":\"127.0.0.1\",\"latency_time\":19208,\"req_method\":\"DELETE\",\"req_uri\":\"/product?price=10\",\"status_code\":200}"
3) "{\"client_ip\":\"127.0.0.1\",\"latency_time\":35000,\"req_method\":\"DELETE\",\"req_uri\":\"/product?price\",\"status_code\":200}"
```

- `Fire()` method marshal data to json and push it to Redis list.
- `Level()` records the level need to be stored.
- `LoggerToRedis` is a middleware logging request information and sends it to Redis.
- `Logrus.New()` create a new instance of the `Logrus.Logger` type. This logger is the main interface through which you interact with the Logrus logging framework.
- `Logger.Info`  used to log an informational message

```go
logger := logrus.New()
logger.SetOutput(myCustomWriter)
logger.SetLevel(logrus.DebugLevel)
logger.SetFormatter(&logrus.TextFormatter{})
```

- `Logger.AddHook` add a Hook to the logger, this hook should be executed whenever a log enrty is made.

- `WithFields` method in Logrus is used to attach additional fields(key-value pairs) to a log entry.

### Cookie

```go
//http/net package
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

```go
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
```

- `c.Cookie("name")` finds the cookie with name, `err = nil` if found none.
- `c.SetCookie()` sets cookie to domain/path.

### JWT

```go
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

```

```sh
~/ curl -X POST -H "Content-Type: application/json" -d '{"username": "root", "password": "123456"}' http://127.0.0.1:8080/auth
{"code":200,"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvb3QiLCJleHAiOjE3MDYzNDM4MDYsImlzcyI6ImZlcnJpZW0ifQ.mD8gHkXvjAi15gFJh28QC6092dRUQYnt_G2UGq8ws28"},"msg":"success"}%

~/ curl -X GET -H "Authorization: Bearer token" http://127.0.0.1:8080/home
{"code":200,"data":"root","msg":"success"}%
```

