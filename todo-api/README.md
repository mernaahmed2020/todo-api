# 📝 Todo API — RESTful Web Service in Go

This is a full-stack RESTful Todo API built with **Go**, **Gin**, **GORM**, and **PostgreSQL**. It supports enhanced CRUD features including categories, priorities, tags, due dates, and priority-based sorting. Built using clean architecture and supports Docker and unit testing.

---

## ✅ Features

- Create, Read, Update, Delete todos
- Filter by:
  - ✅ Category: `/todos/category/:category`
  - ✅ Completion Status: `/todos/status/:status`
  - ✅ Title Search: `/todos/search?q=term`
- Tags support (many-to-many)
- Priority: High > Medium > Low
- Due dates (ISO 8601, UTC only)
- Completed timestamps
- Sorted todos by priority & tag count
- Bulk update all todos by category
- Unit tests using `testing` and `httptest`
- Docker support (Go App + PostgreSQL)

---

---

##  Running with Docker

1. **Start containers:**

   
   docker-compose up --build

## App runs at
http://localhost:8080

## Postgres config:

Host: localhost

Port: 5432

DB: todo_app

User: postgres

Password: postgres


## Run all tests with:
go test ./handlers