package routes

import (
	"trithong.com/task-golang/internal/handlers"

	"github.com/gin-gonic/gin"
)

func HealthRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler) {
	router.GET("/health", healthHandler.Check)
}