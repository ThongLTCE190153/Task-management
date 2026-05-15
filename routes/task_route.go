package routes

import (
	"trithong.com/task-golang/handlers"
	"github.com/gin-gonic/gin"
)

func TaskRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	taskGroup := router.Group("/tasks")

	taskGroup.POST("", taskHandler.CreateTask)
	taskGroup.GET("", taskHandler.GetAllTasks)
	taskGroup.GET("/:id", taskHandler.GetTaskByID)
	taskGroup.PUT("/:id", taskHandler.UpdateTask)
	taskGroup.DELETE("/:id", taskHandler.DeleteTask)
}