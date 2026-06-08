package routes

import (
	"time"

	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ProjectRoutes(router *gin.Engine, projectHandler *handlers.ProjectHandler, redisClient *redis.Client) {
	projectGroup := router.Group("/projects")

	projectGroup.Use(middlewares.AuthMiddleware())
	projectGroup.Use(middlewares.RateLimitMiddleware(redisClient, 60, time.Minute))

	// Admin và Manager mới được tạo/sửa/xóa project
	projectGroup.POST("", middlewares.RoleMiddleware("admin", "manager"), projectHandler.Create)
	projectGroup.PUT("/:id", middlewares.RoleMiddleware("admin", "manager"), projectHandler.Update)
	projectGroup.DELETE("/:id", middlewares.RoleMiddleware("admin", "manager"), projectHandler.Delete)

	// Tất cả role đều xem được
	projectGroup.GET("", projectHandler.GetAll)
	projectGroup.GET("/:id", projectHandler.GetByID)
}