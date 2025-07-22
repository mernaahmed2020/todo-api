package db

import (
	"fmt"
	"log"
	"os"
	"time"
	"todo-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Read values from environment variables (with fallback defaults)
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "secret123")
	dbname := getEnv("DB_NAME", "todo_app")
	port := getEnv("DB_PORT", "5432")

	// Format the DSN string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var err error
	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Connected to PostgreSQL with GORM!")
			err = DB.AutoMigrate(&models.Todo{})
			if err != nil {
				log.Fatal("❌ AutoMigration failed:", err)
			}
			return
		}

		log.Printf("⏳ Attempt %d/%d: Could not connect to database: %v", attempt, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("❌ Could not connect to database after %d attempts: %v", maxAttempts, err)
}

func GetDB() *gorm.DB {
	return DB
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
