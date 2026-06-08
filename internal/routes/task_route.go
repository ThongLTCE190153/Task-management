package routes

import (
	"time"

	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func TaskRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler, redisClient *redis.Client) {
	taskGroup := router.Group("/tasks")

	taskGroup.Use(middlewares.AuthMiddleware())
	taskGroup.Use(middlewares.RateLimitMiddleware(redisClient, 60, time.Minute))

	// Admin và Manager mới được tạo/sửa/xóa task
	taskGroup.POST("", middlewares.RoleMiddleware("admin", "manager"), taskHandler.Create)
	taskGroup.PUT("/:id", middlewares.RoleMiddleware("admin", "manager"), taskHandler.Update)
	taskGroup.DELETE("/:id", middlewares.RoleMiddleware("admin", "manager"), taskHandler.Delete)

	// Tất cả role đều xem được
	taskGroup.GET("", taskHandler.GetAll)
	taskGroup.GET("/:id", taskHandler.GetByID)
}