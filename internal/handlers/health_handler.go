package handlers

import (
	"database/sql"
	"net/http"

	"trithong.com/task-golang/internal/responses"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"context"
)

type HealthHandler struct {
	db          *sql.DB
	redisClient *redis.Client
}

func NewHealthHandler(db *sql.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	status := map[string]string{
		"status":   "ok",
		"postgres": "ok",
		"redis":    "ok",
	}

	httpStatus := http.StatusOK

	// Kiểm tra PostgreSQL
	if err := h.db.Ping(); err != nil {
		status["status"] = "degraded"
		status["postgres"] = "unreachable"
		httpStatus = http.StatusServiceUnavailable
	}

	// Kiểm tra Redis
	if err := h.redisClient.Ping(context.Background()).Err(); err != nil {
		status["status"] = "degraded"
		status["redis"] = "unreachable"
		httpStatus = http.StatusServiceUnavailable
	}

	responses.Success(c, httpStatus, "Health check", status)
}