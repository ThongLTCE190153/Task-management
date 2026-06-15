package handlers

import (
	"net/http"
	"strconv"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid user id")
		return
	}

	userID := c.GetInt("user_id")
	role, _ := c.Get("role")
	userRole := role.(string)

	// User không được update ai
	if userRole == "user" {
		responses.Error(c, http.StatusForbidden, "You do not have permission")
		return
	}

	// Manager chỉ update bản thân hoặc user
	// Manager không update được admin hoặc manager khác
	if userRole == "manager" && userID != id {
		// Kiểm tra target user có phải role "user" không
		targetUser, err := h.userService.GetByID(id)
		if err != nil {
			responses.Error(c, http.StatusNotFound, "User not found")
			return
		}
		if targetUser.Role != "user" {
			responses.Error(c, http.StatusForbidden, "You can only update users with lower role")
			return
		}
	}

	var request dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.userService.Update(id, entities.User{
		FullName: request.FullName,
	})

	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses.Success(c, http.StatusOK, "Update user successfully", mappers.ToUserResponse(result))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid user id")
		return
	}

	userID := c.GetInt("user_id")
	role, _ := c.Get("role")
	userRole := role.(string)

	// User không được xóa ai
	if userRole == "user" {
		responses.Error(c, http.StatusForbidden, "You do not have permission")
		return
	}

	// Không ai được tự xóa chính mình
	if userID == id {
		responses.Error(c, http.StatusForbidden, "You cannot delete yourself")
		return
	}

	// Manager chỉ xóa được user — không xóa được admin hoặc manager khác
	if userRole == "manager" {
		targetUser, err := h.userService.GetByID(id)
		if err != nil {
			responses.Error(c, http.StatusNotFound, "User not found")
			return
		}
		if targetUser.Role != "user" {
			responses.Error(c, http.StatusForbidden, "You can only delete users with lower role")
			return
		}
	}

	err = h.userService.Delete(id)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err.Error())
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
