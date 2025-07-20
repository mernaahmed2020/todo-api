package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=secret123 dbname=todo_app sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening DB connection:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully!")
}

func GetDB() *sql.DB {
	return DB
}
