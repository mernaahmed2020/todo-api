package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	router := gin.Default()

	router.GET("/todos", func(c *gin.Context) {
		rows, err := DB.Query("SELECT id, title, completed FROM todos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
			return
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var t Todo
			err := rows.Scan(&t.ID, &t.Title, &t.Completed)
			if err != nil {
				continue
			}
			todos = append(todos, t)
		}
		c.JSON(http.StatusOK, todos)
	})

	router.GET("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		var t Todo
		err := DB.QueryRow("SELECT id, title, completed FROM todos WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Completed)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	router.POST("/todos", func(c *gin.Context) {
		var t Todo
		if err := c.BindJSON(&t); err != nil || t.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		err := DB.QueryRow("INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id", t.Title, t.Completed).Scan(&t.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	router.PUT("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		var t Todo
		if err := c.BindJSON(&t); err != nil || t.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		result, err := DB.Exec("UPDATE todos SET title = $1, completed = $2 WHERE id = $3", t.Title, t.Completed, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
			return
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		t.ID, _ = strconv.Atoi(id)
		c.JSON(http.StatusOK, t)
	})

	router.DELETE("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		result, err := DB.Exec("DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete todo"})
			return
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
	})

	router.Run(":8080")
}
