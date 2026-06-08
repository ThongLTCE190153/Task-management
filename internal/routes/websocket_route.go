package routes

import (
	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(
	router *gin.Engine,
	webSocketHandler *handlers.WebSocketHandler,
) {
	router.GET("/ws", middlewares.WebSocketAuthMiddleware(), webSocketHandler.HandleConnection)
}