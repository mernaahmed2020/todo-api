# Todo API

A RESTful API for managing todos using Go (Gin) and PostgreSQL. The database runs in a Docker container.

## Setup Instructions

### 1. Clone the repository

git clone https://github.com/mernaahmed2020/todo-api.git 
cd todo-api

### 2. Run PostgreSQL using Docker (Windows PowerShell)
## (i know this is not best practice at all to put secrets in a readme,sorry:)

docker run --name todo-postgres `
  -e POSTGRES_PASSWORD=secret123 `
  -e POSTGRES_DB=todo_app `
  -p 5432:5432 `
  -d postgres
### 3. Update the database connection string in db/db.go if needed:

 connStr := "postgres://postgres:yourpassword@localhost:5432/todo_app?sslmode=disable"

### 4. run application
 go run .



