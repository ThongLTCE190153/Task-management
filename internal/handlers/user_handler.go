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

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid user id")
		return
	}

	user, err := h.userService.GetByID(id)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response := mappers.ToUserResponse(user)

	responses.Success(c, http.StatusOK, "Get user successfully", response)
}

func (h *UserHandler) Update(c *gin.Context) {
	var request dto.UpdateUserRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := c.GetInt("user_id")

	user := mappers.ToUpdateUserEntity(request)

	result, err := h.userService.Update(userID, user)

	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "Update user failed")
		return
	}

	response := mappers.ToUserResponse(result)

	responses.Success(c, http.StatusOK, "Update user successfully", response)
}

func (h *UserHandler) Delete(c *gin.Context) {
	userID := c.GetInt("user_id")

	err := h.userService.Delete(userID)

	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "Delete user failed")
		return
	}

	responses.Success(c, http.StatusOK, "Delete user successfully", nil)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()

	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "Get users failed")
		return
	}

	responses.Success(c, http.StatusOK, "Get users successfully", users)
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid user id")
		return
	}

	var request dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.userService.UpdateRole(id, request.Role)
	if err != nil {
		responses.Error(c, http.StatusNotFound, "User not found")
		return
	}

	responses.Success(c, http.StatusOK, "Update role successfully", nil)
}
