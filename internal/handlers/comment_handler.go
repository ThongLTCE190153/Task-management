package handlers

import (
	"net/http"
	"strconv"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/responses"
	"trithong.com/task-golang/internal/services"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(
	commentService services.CommentService,
) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) Create(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("taskId"))

	if err != nil {
		responses.Error(
			c,
			http.StatusBadRequest,
			"Invalid task id",
		)
		return
	}

	userIDValue, exists := c.Get("user_id")

	if !exists {
		responses.Error(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	userID := userIDValue.(int)

	var request dto.CreateCommentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		responses.Error(
			c,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	comment := entities.Comment{
		TaskID:  taskID,
		UserID:  userID,
		Content: request.Content,
	}

	newComment, err := h.commentService.Create(comment)

	if err != nil {
		responses.Error(
			c,
			http.StatusForbidden, // ← đổi từ InternalServerError sang Forbidden
			err.Error(),
		)
		return
	}

	responses.Success(
		c,
		http.StatusCreated,
		"Comment created successfully",
		newComment,
	)
}

func (h *CommentHandler) GetByTaskID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("taskId"))

	if err != nil {
		responses.Error(
			c,
			http.StatusBadRequest,
			"Invalid task id",
		)
		return
	}

	// ← thêm phần lấy userID
	userIDValue, exists := c.Get("user_id")

	if !exists {
		responses.Error(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	userID := userIDValue.(int)

	// ← truyền thêm userID vào
	comments, err := h.commentService.GetByTaskID(taskID, userID)

	if err != nil {
		responses.Error(
			c,
			http.StatusForbidden, // ← đổi từ InternalServerError sang Forbidden
			err.Error(),
		)
		return
	}

	responses.Success(
		c,
		http.StatusOK,
		"Comments fetched successfully",
		comments,
	)
}