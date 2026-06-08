package main

import (
	"log"

	"trithong.com/task-golang/internal/cache"
	"trithong.com/task-golang/internal/config"
	"trithong.com/task-golang/internal/db"
	"trithong.com/task-golang/internal/handlers"
	"trithong.com/task-golang/internal/middlewares"
	"trithong.com/task-golang/internal/repositories"
	"trithong.com/task-golang/internal/routes"
	"trithong.com/task-golang/internal/services"
	"trithong.com/task-golang/internal/utils"

	appwebsocket "trithong.com/task-golang/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	utils.SetJWTConfig(cfg.JWTSecret, cfg.JWTExpireHours)

	database := db.ConnectPostgres(cfg)
	defer database.Close()

	redisClient := cache.ConnectRedis(cfg)
	defer redisClient.Close()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.RequestID())
	router.Use(middlewares.RequestLogger())

	hub := appwebsocket.NewHub()
	webSocketHandler := handlers.NewWebSocketHandler(hub)

	userRepo := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, redisClient) 
	
	healthHandler := handlers.NewHealthHandler(database, redisClient)


	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	projectRepo := repositories.NewProjectRepository(database)
	projectService := services.NewProjectService(projectRepo, redisClient, hub)
	projectHandler := handlers.NewProjectHandler(projectService)

	taskRepo := repositories.NewTaskRepository(database)
	taskService := services.NewTaskService(taskRepo, redisClient, hub)
	taskHandler := handlers.NewTaskHandler(taskService)

	commentRepo := repositories.NewCommentRepository(database)
	commentService := services.NewCommentService(commentRepo, hub)
	commentHandler := handlers.NewCommentHandler(commentService)

	routes.AuthRoutes(router, authHandler)
	routes.UserRoutes(router, userHandler)
	routes.ProjectRoutes(router, projectHandler, redisClient)
	routes.TaskRoutes(router, taskHandler, redisClient)
	routes.CommentRoutes(router, commentHandler, redisClient)
	routes.WebSocketRoutes(router, webSocketHandler)
	routes.HealthRoutes(router, healthHandler)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
