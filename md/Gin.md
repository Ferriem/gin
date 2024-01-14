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