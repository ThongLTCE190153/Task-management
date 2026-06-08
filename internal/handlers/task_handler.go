package handlers

import (
	"net/http"
	"strconv"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/mappers"
	"trithong.com/task-golang/internal/responses"
	"trithong.com/task-golang/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var request dto.CreateTaskRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := c.GetInt("user_id")

	task := mappers.ToTaskEntity(request)

	result, err := h.taskService.Create(task, userID)

	if err != nil {
		// Phân biệt rõ 2 trường hợp
		if err.Error() == "project not found" {
			responses.Error(c, http.StatusNotFound, "Project not found")
			return
		}
		if err.Error() == "forbidden" {
			responses.Error(c, http.StatusForbidden, "You do not have permission")
			return
		}
		responses.Error(c, http.StatusInternalServerError, "Create task failed")
		return
	}

	response := mappers.ToTaskResponse(result)

	responses.Success(c, http.StatusCreated, "Create task successfully", response)
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	userID := c.GetInt("user_id")

	var params dto.TaskQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid query params")
		return
	}

	// Nếu không có pagination params thì dùng filter đơn giản
	if params.Page == 0 && params.Limit == 0 &&
		params.Status == "" && params.AssigneeID == nil && params.ProjectID == nil {
		tasks, err := h.taskService.GetAllByOwnerID(userID)
		if err != nil {
			responses.Error(c, http.StatusInternalServerError, "Get tasks failed")
			return
		}
		responses.Success(c, http.StatusOK, "Get tasks successfully", mappers.ToTaskResponses(tasks))
		return
	}

	tasks, total, err := h.taskService.GetAllWithFilter(userID, params)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "Get tasks failed")
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	responses.Success(c, http.StatusOK, "Get tasks successfully", dto.PaginationResponse{
		Items:      mappers.ToTaskResponses(tasks),
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	})
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid task id")
		return
	}

	userID := c.GetInt("user_id")

	task, err := h.taskService.GetByIDAndOwnerID(id, userID)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	response := mappers.ToTaskResponse(task)

	responses.Success(c, http.StatusOK, "Get task successfully", response)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid task id")
		return
	}

	var request dto.UpdateTaskRequest

	err = c.ShouldBindJSON(&request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := c.GetInt("user_id")

	task := mappers.ToUpdateTaskEntity(request)

	result, err := h.taskService.Update(id, userID, task)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	response := mappers.ToTaskResponse(result)

	responses.Success(c, http.StatusOK, "Update task successfully", response)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid task id")
		return
	}

	userID := c.GetInt("user_id")

	err = h.taskService.Delete(id, userID)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	responses.Success(c, http.StatusOK, "Delete task successfully", nil)
}
