package services

import (
	"errors"
	"testing"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
	appwebsocket "trithong.com/task-golang/internal/websocket"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== MOCK PROJECT REPOSITORY =====

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(project entities.Project) (entities.Project, error) {
	args := m.Called(project)
	return args.Get(0).(entities.Project), args.Error(1)
}

func (m *MockProjectRepository) GetAllByOwnerID(ownerID int) ([]entities.Project, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]entities.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByIDAndOwnerID(id int, ownerID int) (entities.Project, error) {
	args := m.Called(id, ownerID)
	return args.Get(0).(entities.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(id int, ownerID int, project entities.Project) (entities.Project, error) {
	args := m.Called(id, ownerID, project)
	return args.Get(0).(entities.Project), args.Error(1)
}

func (m *MockProjectRepository) Delete(id int, ownerID int) error {
	args := m.Called(id, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) GetAllWithPagination(ownerID int, params dto.ProjectQueryParams) ([]entities.Project, int, error) {
	args := m.Called(ownerID, params)
	return args.Get(0).([]entities.Project), args.Int(1), args.Error(2)
}

// ===== HELPER =====

func newTestRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func newTestHub() *appwebsocket.Hub {
	return appwebsocket.NewHub()
}

// ===== TEST CREATE =====

func TestCreateProject_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)

	mockRepo.On("Create", mock.AnythingOfType("entities.Project")).
		Return(entities.Project{
			ID:      1,
			Name:    "Project A",
			OwnerID: 1,
		}, nil)

	service := NewProjectService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Create(entities.Project{
		Name:    "Project A",
		OwnerID: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Project A", result.Name)
	assert.Equal(t, 1, result.OwnerID)

	mockRepo.AssertExpectations(t)
}

func TestCreateProject_Failed(t *testing.T) {
	mockRepo := new(MockProjectRepository)

	mockRepo.On("Create", mock.AnythingOfType("entities.Project")).
		Return(entities.Project{}, errors.New("database error"))

	service := NewProjectService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Create(entities.Project{
		Name:    "Project A",
		OwnerID: 1,
	})

	assert.Error(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

// ===== TEST GET ALL =====

func TestGetAllProjectsByOwnerID_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)

	mockRepo.On("GetAllByOwnerID", 1).
		Return([]entities.Project{
			{ID: 1, Name: "Project A", OwnerID: 1},
			{ID: 2, Name: "Project B", OwnerID: 1},
		}, nil)

	service := NewProjectService(mockRepo, newTestRedis(), newTestHub())

	results, err := service.GetAllByOwnerID(1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)

	mockRepo.AssertExpectations(t)
}

// ===== TEST UPDATE =====

func TestUpdateProject_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)

	mockRepo.On("Update", 1, 1, mock.AnythingOfType("entities.Project")).
		Return(entities.Project{
			ID:      1,
			Name:    "New Name",
			OwnerID: 1,
		}, nil)

	service := NewProjectService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Update(1, 1, entities.Project{Name: "New Name"})

	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProject_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)

	mockRepo.On("Update", 1, 99, mock.AnythingOfType("entities.Project")).
		Return(entities.Project{}, errors.New("project not found"))

	service := NewProjectService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Update(1, 99, entities.Project{Name: "New Name"})

	assert.Error(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

// ===== TEST DELETE =====

func TestDeleteProject_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)

	mockRepo.On("Delete", 1, 1).Return(nil)

	service := NewProjectService(mockRepo, newTestRedis(), newTestHub())

	err := service.Delete(1, 1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteProject_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)

	mockRepo.On("Delete", 1, 99).
		Return(errors.New("project not found"))

	service := NewProjectService(mockRepo, newTestRedis(), newTestHub())

	err := service.Delete(1, 99)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}