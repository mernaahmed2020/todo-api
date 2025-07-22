package handlers

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()
	RegisterRoutes(r)
	return r
}
