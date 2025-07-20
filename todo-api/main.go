package main

import (
	"log"
	"todo-api/db"
	"todo-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()
	handlers.RegisterRoutes(r)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
