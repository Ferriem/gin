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
