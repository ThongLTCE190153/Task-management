package services

import (
	"trithong.com/task-golang/entities"
	"trithong.com/task-golang/repositories"
	"errors"
	"strings"
)

type TaskService interface {
	CreateTask(task entities.Task) (entities.Task, error)
	GetTaskByID(id int) (entities.Task, error)
	GetAllTasks() []entities.Task
	UpdateTask(id int, task entities.Task) (entities.Task, error)
	DeleteTask(id int) error
}

type taskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (s *taskService) CreateTask(task entities.Task) (entities.Task, error) {
	err := validateTask(task)

	if err != nil {
		return entities.Task{}, err
	}

	return s.repo.Create(task), nil
}

func (s *taskService) GetTaskByID(id int) (entities.Task, error) {
	return s.repo.GetByID(id)
}

func (s *taskService) GetAllTasks() []entities.Task {
	return s.repo.GetAll()
}

func (s *taskService) UpdateTask(id int, task entities.Task) (entities.Task, error) {
	err := validateTask(task)

	if err != nil {
		return entities.Task{}, err
	}

	return s.repo.Update(id, task)
}

func (s *taskService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}

func validateTask(task entities.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("title is required")
	}

	if strings.TrimSpace(task.Description) == "" {
		return errors.New("description is required")
	}

	if strings.TrimSpace(task.Assignee) == "" {
		return errors.New("assignee is required")
	}

	if task.Status != "todo" && task.Status != "doing" && task.Status != "done" {
		return errors.New("status must be todo, doing or done")
	}

	return nil
}