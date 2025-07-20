package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{}
var nextID = 1

func main() {
	router := gin.Default()

	router.GET("/todos", func(c *gin.Context) {
		c.JSON(200, todos)
	})

	router.GET("/todos/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for _, todo := range todos {
			if todo.ID == id {
				c.JSON(200, todo)
				return
			}
		}
		c.JSON(404, gin.H{"error": "didn't find todo"})
	})

	router.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.BindJSON(&todo); err != nil || todo.Title == "" {
			c.JSON(400, gin.H{"error": "wrong or empty title"})
			return
		}
		todo.ID = nextID
		nextID++
		todos = append(todos, todo)
		c.JSON(200, todo)
	})

	router.PUT("/todos/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var data Todo
		if err := c.BindJSON(&data); err != nil || data.Title == "" {
			c.JSON(400, gin.H{"error": "wrong or empty title"})
			return
		}
		for i := range todos {
			if todos[i].ID == id {
				todos[i].Title = data.Title
				todos[i].Completed = data.Completed
				c.JSON(200, todos[i])
				return
			}
		}
		c.JSON(404, gin.H{"error": "Todo not found"})
	})

	router.DELETE("/todos/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for i := range todos {
			if todos[i].ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				c.JSON(200, gin.H{"message": "deleted todo successfully"})
				return
			}
		}
		c.JSON(404, gin.H{"error": "Todo not found"})
	})

	router.Run(":8080")
}
