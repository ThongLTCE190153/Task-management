package routes

import (
	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(
	router *gin.Engine, authHandler *handlers.AuthHandler) {
	authGroup := router.Group("/auth")

	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.RefreshToken)
	authGroup.POST("/logout", middlewares.AuthMiddleware(), authHandler.Logout)
	authGroup.GET("/me", middlewares.AuthMiddleware(), authHandler.Me)

}
