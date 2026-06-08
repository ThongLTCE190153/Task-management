package routes

import (
	"time"

	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func CommentRoutes(router *gin.Engine, commentHandler *handlers.CommentHandler, redisClient *redis.Client) {
	commentGroup := router.Group("/comments")

	commentGroup.Use(middlewares.AuthMiddleware())
	commentGroup.Use(middlewares.RateLimitMiddleware(redisClient, 30, time.Minute))

	// Tất cả role đều comment được
	commentGroup.POST("/task/:taskId", commentHandler.Create)
	commentGroup.GET("/task/:taskId", commentHandler.GetByTaskID)
}