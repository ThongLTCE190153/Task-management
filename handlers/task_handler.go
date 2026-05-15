package handlers

import (
	"trithong.com/task-golang/entities"
	"trithong.com/task-golang/responses"
	"trithong.com/task-golang/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task entities.Task

	err := c.ShouldBindJSON(&task)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	newTask, err := h.service.CreateTask(task)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.Success(c, http.StatusCreated, "Create task successfully", newTask)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks := h.service.GetAllTasks()

	responses.Success(c, http.StatusOK, "Get all tasks successfully", tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid task id")
		return
	}

	task, err := h.service.GetTaskByID(id)

	if err != nil {
		responses.Error(c, http.StatusNotFound, err.Error())
		return
	}

	responses.Success(c, http.StatusOK, "Get task successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid task id")
		return
	}

	var task entities.Task

	err = c.ShouldBindJSON(&task)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedTask, err := h.service.UpdateTask(id, task)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.Success(c, http.StatusOK, "Update task successfully", updatedTask)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid task id")
		return
	}

	err = h.service.DeleteTask(id)

	if err != nil {
		responses.Error(c, http.StatusNotFound, err.Error())
		return
	}

	responses.Success(c, http.StatusOK, "Delete task successfully", nil)
}