package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todo-api/db"
	"todo-api/models"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/todos", GetAllTodos)
	r.GET("/todos/:id", GetTodoByID)
	r.GET("/todos/category/:category", GetTodosByCategory)
	r.GET("/todos/status/:status", GetTodosByStatus)
	r.GET("/todos/search", SearchTodosByTitle)

	r.POST("/todos", CreateTodo)
	r.PUT("/todos/:id", UpdateTodoByID)
	r.PUT("/todos/category/:category", BulkUpdateByCategory)

	r.DELETE("/todos/:id", DeleteTodoByID)
	r.DELETE("/todos", DeleteAllTodos)
}

func GetAllTodos(c *gin.Context) {
	rows, err := db.GetDB().Query("SELECT id, title, completed, category, priority, completedAt, dueDate FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Category, &t.Priority, &t.CompletedAt, &t.DueDate)
		if err != nil {
			continue
		}
		todos = append(todos, t)
	}
	c.JSON(http.StatusOK, todos)
}

func GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	var t models.Todo
	err := db.GetDB().QueryRow("SELECT id, title, completed, category, priority, completedAt, dueDate FROM todos WHERE id = $1", id).
		Scan(&t.ID, &t.Title, &t.Completed, &t.Category, &t.Priority, &t.CompletedAt, &t.DueDate)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func GetTodosByCategory(c *gin.Context) {
	category := c.Param("category")
	rows, err := db.GetDB().Query("SELECT id, title, completed, category, priority, completedAt, dueDate FROM todos WHERE category = $1", category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Category, &t.Priority, &t.CompletedAt, &t.DueDate); err == nil {
			todos = append(todos, t)
		}
	}
	c.JSON(http.StatusOK, todos)
}

func GetTodosByStatus(c *gin.Context) {
	status := c.Param("status")
	completed := false
	if status == "true" {
		completed = true
	} else if status != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value (must be true or false)"})
		return
	}

	rows, err := db.GetDB().Query("SELECT id, title, completed, category, priority, completedAt, dueDate FROM todos WHERE completed = $1", completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error"})
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Category, &t.Priority, &t.CompletedAt, &t.DueDate); err == nil {
			todos = append(todos, t)
		}
	}
	c.JSON(http.StatusOK, todos)
}

func SearchTodosByTitle(c *gin.Context) {
	q := c.Query("q")
	like := "%%%s%%"
	rows, err := db.GetDB().Query("SELECT id, title, completed, category, priority, completedAt, dueDate FROM todos WHERE LOWER(title) LIKE LOWER($1)", fmt.Sprintf(like, q))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error"})
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Category, &t.Priority, &t.CompletedAt, &t.DueDate); err == nil {
			todos = append(todos, t)
		}
	}
	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {
	var t models.Todo
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	validPriorities := map[string]bool{"Low": true, "Medium": true, "High": true}
	if !validPriorities[t.Priority] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Priority must be Low, Medium, or High"})
		return
	}

	now := time.Now().UTC()
	if t.DueDate != nil && t.DueDate.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Due date cannot be in the past"})
		return
	}

	if t.Completed {
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}

	err := db.GetDB().QueryRow(`INSERT INTO todos (title, completed, category, priority, completedAt, dueDate) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		t.Title, t.Completed, t.Category, t.Priority, t.CompletedAt, t.DueDate).Scan(&t.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func UpdateTodoByID(c *gin.Context) {
	id := c.Param("id")
	var t models.Todo
	if err := c.BindJSON(&t); err != nil || strings.TrimSpace(t.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	now := time.Now().UTC()
	if t.Completed {
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}

	result, err := db.GetDB().Exec(`UPDATE todos SET title = $1, completed = $2, category = $3, priority = $4, completedAt = $5, dueDate = $6 WHERE id = $7`,
		t.Title, t.Completed, t.Category, t.Priority, t.CompletedAt, t.DueDate, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
		return
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	t.ID = parseInt(id)
	c.JSON(http.StatusOK, t)
}

func BulkUpdateByCategory(c *gin.Context) {
	category := c.Param("category")
	var body struct {
		Completed bool `json:"completed"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	now := time.Now().UTC()
	var completedAt *time.Time
	if body.Completed {
		completedAt = &now
	}

	_, err := db.GetDB().Exec(`UPDATE todos SET completed = $1, completedAt = $2 WHERE category = $3`,
		body.Completed, completedAt, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not bulk update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todos updated"})
}

func DeleteTodoByID(c *gin.Context) {
	id := c.Param("id")
	result, err := db.GetDB().Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func DeleteAllTodos(c *gin.Context) {
	_, err := db.GetDB().Exec("DELETE FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete all todos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All todos deleted"})
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
