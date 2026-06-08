package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"trithong.com/task-golang/internal/cache"
	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/repositories"
	"trithong.com/task-golang/internal/utils"
)

type AuthService interface {
	Register(request dto.RegisterRequest) (entities.User, error)
	Login(request dto.LoginRequest) (dto.LoginResponse, error)
	Me(userID int) (entities.User, error)
	RefreshToken(refreshToken string) (dto.LoginResponse, error)
	Logout(userID int) error
}

type authService struct {
	userRepo    repositories.UserRepository
	redisClient *redis.Client
}

func NewAuthService(
	userRepo repositories.UserRepository,
	redisClient *redis.Client,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *authService) Register(request dto.RegisterRequest) (entities.User, error) {
	_, err := s.userRepo.GetByEmail(request.Email)

	if err == nil {
		return entities.User{}, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		return entities.User{}, err
	}

	user := entities.User{
		Email:        request.Email,
		PasswordHash: hashedPassword,
		FullName:     request.FullName,
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(request dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(request.Email)

	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	err = utils.CheckPassword(user.PasswordHash, request.Password)

	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateJWT(user.ID, user.Role)

	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)

	if err != nil {
		return dto.LoginResponse{}, err
	}

	key := fmt.Sprintf("refresh_token:%d", user.ID)
	err = s.redisClient.Set(cache.Ctx, key, refreshToken, 7*24*time.Hour).Err()

	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}

func (s *authService) Me(userID int) (entities.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *authService) RefreshToken(refreshToken string) (dto.LoginResponse, error) {
	token, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil || !token.Valid {
		return dto.LoginResponse{}, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return dto.LoginResponse{}, errors.New("invalid token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return dto.LoginResponse{}, errors.New("invalid token type")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return dto.LoginResponse{}, errors.New("invalid user id in token")
	}

	userID := int(userIDFloat)

	key := fmt.Sprintf("refresh_token:%d", userID)
	savedToken, err := s.redisClient.Get(cache.Ctx, key).Result()
	if err != nil || savedToken != refreshToken {
		return dto.LoginResponse{}, errors.New("refresh token has been revoked")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return dto.LoginResponse{}, errors.New("user not found")
	}

	newAccessToken, err := utils.GenerateJWT(userID, user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	s.redisClient.Set(cache.Ctx, key, newRefreshToken, 7*24*time.Hour)

	return dto.LoginResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}

func (s *authService) Logout(userID int) error {
	key := fmt.Sprintf("refresh_token:%d", userID)
	return s.redisClient.Del(cache.Ctx, key).Err()
}