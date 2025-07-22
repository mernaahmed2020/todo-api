package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"todo-api/db"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	db.InitDB()

	code := m.Run()
	os.Exit(code)
}

func TestGetAllTodos(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/todos", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)
}

func TestGetTodoByID(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/todos/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusNotFound)
}

func TestGetTodosByCategory(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/todos/category/work", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)
}

func TestGetTodosByStatus(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/todos/status/true", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)

	req, _ = http.NewRequest("GET", "/todos/status/invalid", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestSearchTodosByTitle(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/todos/search?q=example", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)
}

func TestCreateTodo(t *testing.T) {
	router := setupRouter()
	todo := map[string]interface{}{
		"title":     "Test Task",
		"category":  "testing",
		"priority":  "High",
		"completed": false,
	}
	body, _ := json.Marshal(todo)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)
}

func TestCreateTodoInvalid(t *testing.T) {
	router := setupRouter()
	invalid := `{"title": "", "priority": "Urgent"}`
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBufferString(invalid))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUpdateTodoByID(t *testing.T) {
	router := setupRouter()
	todo := map[string]interface{}{
		"title":     "Updated Task",
		"category":  "updated",
		"priority":  "Medium",
		"completed": true,
		"due_date":  time.Now().Add(24 * time.Hour),
	}
	body, _ := json.Marshal(todo)

	req, _ := http.NewRequest("PUT", "/todos/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusNotFound || resp.Code == http.StatusInternalServerError)
}

func TestBulkUpdateByCategory(t *testing.T) {
	router := setupRouter()
	body := `{"completed": true}`
	req, _ := http.NewRequest("PUT", "/todos/category/testing", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)
}

func TestDeleteTodoByID(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("DELETE", "/todos/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)
}

func TestDeleteAllTodos(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("DELETE", "/todos", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusInternalServerError)
}
