package main

import (
	"trithong.com/task-golang/handlers"
	"trithong.com/task-golang/middlewares"
	"trithong.com/task-golang/repositories"
	"trithong.com/task-golang/routes"
	"trithong.com/task-golang/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.RequestID())

	taskRepo := repositories.NewTaskRepository()
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	routes.TaskRoutes(router, taskHandler)

	router.Run(":8080")
}