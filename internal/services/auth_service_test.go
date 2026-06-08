package services

import (
	"errors"
	"testing"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/utils"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== MOCK =====

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByEmail(email string) (entities.User, error) {
	args := m.Called(email)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) Create(user entities.User) (entities.User, error) {
	args := m.Called(user)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (entities.User, error) {
	args := m.Called(id)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(id int, user entities.User) (entities.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll() ([]entities.User, error) {
	args := m.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *MockUserRepository) UpdateRole(id int, role string) error {
	args := m.Called(id, role)
	return args.Error(0)
}

// ===== REDIS MOCK =====

func newTestRedisForAuth(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	t.Cleanup(func() {
		client.Close()
		mr.Close()
	})

	return client
}

// ===== TEST REGISTER =====

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetByEmail", "test@gmail.com").
		Return(entities.User{}, errors.New("not found"))

	mockRepo.On("Create", mock.AnythingOfType("entities.User")).
		Return(entities.User{
			ID:       1,
			Email:    "test@gmail.com",
			FullName: "Test User",
			Role:     "user",
		}, nil)

	service := NewAuthService(mockRepo, newTestRedisForAuth(t))

	result, err := service.Register(dto.RegisterRequest{
		Email:    "test@gmail.com",
		Password: "123456",
		FullName: "Test User",
	})

	assert.NoError(t, err)
	assert.Equal(t, "test@gmail.com", result.Email)
	assert.Equal(t, "Test User", result.FullName)

	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetByEmail", "existed@gmail.com").
		Return(entities.User{ID: 1, Email: "existed@gmail.com"}, nil)

	service := NewAuthService(mockRepo, newTestRedisForAuth(t))

	result, err := service.Register(dto.RegisterRequest{
		Email:    "existed@gmail.com",
		Password: "123456",
		FullName: "Test User",
	})

	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

// ===== TEST LOGIN =====

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	hashedPassword, _ := utils.HashPassword("123456")

	mockRepo.On("GetByEmail", "test@gmail.com").
		Return(entities.User{
			ID:           1,
			Email:        "test@gmail.com",
			PasswordHash: hashedPassword,
			FullName:     "Test User",
			Role:         "user",
		}, nil)

	utils.SetJWTConfig("test-secret", 24)

	service := NewAuthService(mockRepo, newTestRedisForAuth(t))

	result, err := service.Login(dto.LoginRequest{
		Email:    "test@gmail.com",
		Password: "123456",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "test@gmail.com", result.User.Email)

	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)

	hashedPassword, _ := utils.HashPassword("correctpassword")

	mockRepo.On("GetByEmail", "test@gmail.com").
		Return(entities.User{
			ID:           1,
			Email:        "test@gmail.com",
			PasswordHash: hashedPassword,
			Role:         "user",
		}, nil)

	service := NewAuthService(mockRepo, newTestRedisForAuth(t))

	result, err := service.Login(dto.LoginRequest{
		Email:    "test@gmail.com",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid email or password", err.Error())
	assert.Empty(t, result.Token)

	mockRepo.AssertExpectations(t)
}

func TestLogin_EmailNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetByEmail", "notfound@gmail.com").
		Return(entities.User{}, errors.New("not found"))

	service := NewAuthService(mockRepo, newTestRedisForAuth(t))

	result, err := service.Login(dto.LoginRequest{
		Email:    "notfound@gmail.com",
		Password: "123456",
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid email or password", err.Error())
	assert.Empty(t, result.Token)

	mockRepo.AssertExpectations(t)
}