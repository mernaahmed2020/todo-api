package handlers

import (
	"net/http"
	"sort"
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
	r.GET("/todos/sorted", GetSortedTodos)
}

func GetAllTodos(c *gin.Context) {
	var todos []models.Todo
	if err := db.GetDB().Preload("Tags").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := db.GetDB().Preload("Tags").First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func GetTodosByCategory(c *gin.Context) {
	category := c.Param("category")
	var todos []models.Todo
	if err := db.GetDB().Preload("Tags").Where("category = ?", category).Find(&todos).Error; err != nil {
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
	if err := db.GetDB().Preload("Tags").Where("completed = ?", completed).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func SearchTodosByTitle(c *gin.Context) {
	q := c.Query("q")
	var todos []models.Todo
	if err := db.GetDB().Preload("Tags").Where("LOWER(title) LIKE LOWER(?)", "%"+q+"%").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {
	var input struct {
		Title     string     `json:"title"`
		Completed bool       `json:"completed"`
		Category  string     `json:"category"`
		Priority  string     `json:"priority"`
		DueDate   *time.Time `json:"dueDate"`
		Tags      []string   `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	validPriorities := map[string]bool{"Low": true, "Medium": true, "High": true}
	if !validPriorities[input.Priority] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Priority must be Low, Medium, or High"})
		return
	}

	now := time.Now().UTC()
	if input.DueDate != nil && input.DueDate.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Due date cannot be in the past"})
		return
	}

	tagMap := map[string]bool{}
	var tagModels []*models.Tag

	for _, tag := range input.Tags {
		tag = strings.TrimSpace(tag)

		if tag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tags must not be empty"})
			return
		}
		if len(tag) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag cannot exceed 50 characters"})
			return
		}
		if tagMap[tag] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate tags are not allowed"})
			return
		}
		tagMap[tag] = true

		var tagModel models.Tag
		if err := db.GetDB().FirstOrCreate(&tagModel, models.Tag{Tag: tag}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tag: " + tag})
			return
		}
		tagModels = append(tagModels, &tagModel)
	}

	todo := models.Todo{
		Title:     input.Title,
		Completed: input.Completed,
		Category:  input.Category,
		Priority:  input.Priority,
		DueDate:   input.DueDate,
		Tags:      tagModels,
	}

	if input.Completed {
		todo.CompletedAt = &now
	}

	if err := db.GetDB().Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func UpdateTodoByID(c *gin.Context) {
	id := c.Param("id")

	// Define a separate input struct to cleanly accept JSON
	var input struct {
		Title     string     `json:"title"`
		Completed bool       `json:"completed"`
		Category  string     `json:"category"`
		Priority  string     `json:"priority"`
		DueDate   *time.Time `json:"dueDate"`
		Tags      []string   `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Step 1: Fetch the existing todo
	var existing models.Todo
	if err := db.GetDB().Preload("Tags").First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// Step 2: Validate priority
	validPriorities := map[string]bool{"Low": true, "Medium": true, "High": true}
	if !validPriorities[input.Priority] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Priority must be Low, Medium, or High"})
		return
	}

	// Step 3: Validate due date
	now := time.Now().UTC()
	if input.DueDate != nil && input.DueDate.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Due date cannot be in the past"})
		return
	}

	// Step 4: Handle tags validation and association
	tagMap := map[string]bool{}
	var tagModels []*models.Tag

	for _, tag := range input.Tags {
		tag = strings.TrimSpace(tag)

		if tag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tags must not be empty"})
			return
		}
		if len(tag) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag cannot exceed 50 characters"})
			return
		}
		if tagMap[tag] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate tags are not allowed"})
			return
		}
		tagMap[tag] = true

		var tagModel models.Tag
		if err := db.GetDB().FirstOrCreate(&tagModel, models.Tag{Tag: tag}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tag: " + tag})
			return
		}
		tagModels = append(tagModels, &tagModel)
	}

	// Step 5: Apply updates to todo fields
	existing.Title = strings.TrimSpace(input.Title)
	existing.Completed = input.Completed
	existing.Category = input.Category
	existing.Priority = input.Priority
	existing.DueDate = input.DueDate

	if input.Completed {
		existing.CompletedAt = &now
	} else {
		existing.CompletedAt = nil
	}

	// Step 6: Save todo
	if err := db.GetDB().Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
		return
	}

	// Step 7: Update tags (many-to-many association)
	if err := db.GetDB().Model(&existing).Association("Tags").Replace(tagModels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update tags"})
		return
	}

	// Step 8: Return updated todo
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

func GetSortedTodos(c *gin.Context) {
	var todos []models.Todo

	if err := db.GetDB().Preload("Tags").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch todos"})
		return
	}

	priorityOrder := map[string]int{
		"High":   3,
		"Medium": 2,
		"Low":    1,
	}

	sort.SliceStable(todos, func(i, j int) bool {
		pi := priorityOrder[todos[i].Priority]
		pj := priorityOrder[todos[j].Priority]

		if pi != pj {
			return pi > pj
		}

		return len(todos[i].Tags) > len(todos[j].Tags)
	})

	c.JSON(http.StatusOK, todos)
}
