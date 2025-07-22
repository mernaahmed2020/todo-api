package handlers

import (
	"net/http"
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
	var todos []models.Todo
	if err := db.GetDB().Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := db.GetDB().First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func GetTodosByCategory(c *gin.Context) {
	category := c.Param("category")
	var todos []models.Todo
	if err := db.GetDB().Where("category = ?", category).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func GetTodosByStatus(c *gin.Context) {
	status := c.Param("status")
	var completed bool
	if status == "true" {
		completed = true
	} else if status == "false" {
		completed = false
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value (must be true or false)"})
		return
	}

	var todos []models.Todo
	if err := db.GetDB().Where("completed = ?", completed).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func SearchTodosByTitle(c *gin.Context) {
	q := c.Query("q")
	var todos []models.Todo
	if err := db.GetDB().Where("LOWER(title) LIKE LOWER(?)", "%"+q+"%").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
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

	if err := db.GetDB().Create(&t).Error; err != nil {
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

	var existing models.Todo
	if err := db.GetDB().First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	now := time.Now().UTC()
	if t.Completed {
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}

	// Update fields manually to avoid overwriting ID
	existing.Title = t.Title
	existing.Completed = t.Completed
	existing.Category = t.Category
	existing.Priority = t.Priority
	existing.CompletedAt = t.CompletedAt
	existing.DueDate = t.DueDate

	if err := db.GetDB().Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
		return
	}

	c.JSON(http.StatusOK, existing)
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

	if err := db.GetDB().Model(&models.Todo{}).
		Where("category = ?", category).
		Updates(map[string]interface{}{"completed": body.Completed, "completed_at": completedAt}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bulk update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todos updated"})
}

func DeleteTodoByID(c *gin.Context) {
	id := c.Param("id")
	if err := db.GetDB().Delete(&models.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func DeleteAllTodos(c *gin.Context) {
	if err := db.GetDB().Where("1 = 1").Delete(&models.Todo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete all failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All todos deleted"})
}
