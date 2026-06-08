package services

import (
	"errors"
	"testing"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)
// ===== MOCK TASK REPOSITORY =====

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task entities.Task, ownerID int) (entities.Task, error) {
	args := m.Called(task, ownerID)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockTaskRepository) GetAllByOwnerID(ownerID int) ([]entities.Task, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockTaskRepository) GetAllByProjectID(projectID int, ownerID int) ([]entities.Task, error) {
	args := m.Called(projectID, ownerID)
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByIDAndOwnerID(id int, ownerID int) (entities.Task, error) {
	args := m.Called(id, ownerID)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(id int, ownerID int, task entities.Task) (entities.Task, error) {
	args := m.Called(id, ownerID, task)
	return args.Get(0).(entities.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(id int, ownerID int) error {
	args := m.Called(id, ownerID)
	return args.Error(0)
}

func (m *MockTaskRepository) GetOwnerIDByProjectID(projectID int) (int, error) {
	args := m.Called(projectID)
	return args.Int(0), args.Error(1)
}

func (m *MockTaskRepository) GetAllWithFilter(ownerID int, params dto.TaskQueryParams) ([]entities.Task, int, error) {
	args := m.Called(ownerID, params)
	return args.Get(0).([]entities.Task), args.Int(1), args.Error(2)
}

// ===== TEST CREATE =====

func TestCreateTask_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	
	mockRepo.On("Create", mock.AnythingOfType("entities.Task"),1).
		Return(entities.Task{
			ID:        1,
			ProjectID: 1,
			Title:     "Task A",
			Status:    "todo",
		}, nil)

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Create(entities.Task{
		ProjectID: 1,
		Title:     "Task A",
		Status:    "todo",
	}, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Task A", result.Title)
	assert.Equal(t, "todo", result.Status)

	mockRepo.AssertExpectations(t)
}

func TestCreateTask_ProjectNotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	mockRepo.On("Create", mock.AnythingOfType("entities.Task"), 1).
		Return(entities.Task{}, errors.New("project not found"))

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Create(entities.Task{
		ProjectID: 99,
		Title:     "Task A",
	}, 1)

	assert.Error(t, err)
	assert.Equal(t, "project not found", err.Error())
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCreateTask_Forbidden(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	mockRepo.On("Create", mock.AnythingOfType("entities.Task"), 99).
		Return(entities.Task{}, errors.New("forbidden"))

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Create(entities.Task{
		ProjectID: 1,
		Title:     "Task A",
	}, 99)

	assert.Error(t, err)
	assert.Equal(t, "forbidden", err.Error())
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

// ===== TEST GET ALL =====

func TestGetAllTasksByOwnerID_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	mockRepo.On("GetAllByOwnerID", 1).
		Return([]entities.Task{
			{ID: 1, Title: "Task A", Status: "todo"},
			{ID: 2, Title: "Task B", Status: "done"},
		}, nil)

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	results, err := service.GetAllByOwnerID(1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)

	mockRepo.AssertExpectations(t)
}

// ===== TEST UPDATE =====

func TestUpdateTask_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	mockRepo.On("Update", 1, 1, mock.AnythingOfType("entities.Task")).
		Return(entities.Task{
			ID:     1,
			Title:  "Updated Task",
			Status: "in_progress",
		}, nil)

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Update(1, 1, entities.Task{
		Title:  "Updated Task",
		Status: "in_progress",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", result.Title)
	assert.Equal(t, "in_progress", result.Status)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTask_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	mockRepo.On("Update", 99, 1, mock.AnythingOfType("entities.Task")).
		Return(entities.Task{}, errors.New("task not found"))

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	result, err := service.Update(99, 1, entities.Task{
		Title:  "Updated Task",
		Status: "in_progress",
	})

	assert.Error(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

// ===== TEST DELETE =====

func TestDeleteTask_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	mockRepo.On("GetByIDAndOwnerID", 1, 1).
		Return(entities.Task{ID: 1, ProjectID: 1}, nil)

	mockRepo.On("Delete", 1, 1).Return(nil)

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	err := service.Delete(1, 1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteTask_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	mockRepo.On("GetByIDAndOwnerID", 99, 1).
		Return(entities.Task{}, errors.New("task not found"))

	service := NewTaskService(mockRepo, newTestRedis(), newTestHub())

	err := service.Delete(99, 1)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}