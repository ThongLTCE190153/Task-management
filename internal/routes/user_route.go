package routes

import (
	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	userGroup := router.Group("/users")

	userGroup.Use(middlewares.AuthMiddleware())

	// Chỉ Admin mới được quản lý user
	userGroup.GET("", middlewares.RoleMiddleware("admin"), userHandler.GetAll)
	userGroup.PUT("/:id/role", middlewares.RoleMiddleware("admin"), userHandler.UpdateRole)

	// User tự quản lý thông tin của mình
	userGroup.GET("/:id", userHandler.GetByID)
	userGroup.PUT("/:id", userHandler.Update)
	userGroup.DELETE("/:id", userHandler.Delete)
}