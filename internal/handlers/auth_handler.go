package handlers

import (
	"net/http"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/mappers"
	"trithong.com/task-golang/internal/responses"
	"trithong.com/task-golang/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.authService.Register(request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response := mappers.ToUserResponse(user)

	responses.Success(c, http.StatusCreated, "Register successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Login(request)

	if err != nil {
		responses.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	responses.Success(c, http.StatusOK, "Login successfully", response)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt("user_id")

	user, err := h.authService.Me(userID)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response := mappers.ToUserResponse(user)

	responses.Success(c, http.StatusOK, "Get profile successfully", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.authService.RefreshToken(request.RefreshToken)

	if err != nil {
		responses.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	responses.Success(c, http.StatusOK, "Token refreshed successfully", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetInt("user_id")

	err := h.authService.Logout(userID)

	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "Logout failed")
		return
	}

	responses.Success(c, http.StatusOK, "Logout successfully", nil)
}