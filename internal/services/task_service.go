package services

import (
	"encoding/json"
	"fmt"
	"time"

	"trithong.com/task-golang/internal/cache"
	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/jobs"
	"trithong.com/task-golang/internal/repositories"
	"trithong.com/task-golang/internal/dto"
	appwebsocket "trithong.com/task-golang/internal/websocket"

	"github.com/redis/go-redis/v9"
)

type TaskService interface {
	Create(task entities.Task, ownerID int) (entities.Task, error)
	GetAllByOwnerID(ownerID int) ([]entities.Task, error)
	GetAllWithFilter(ownerID int, params dto.TaskQueryParams) ([]entities.Task, int, error)
	GetByIDAndOwnerID(id int, ownerID int) (entities.Task, error)
	GetAllByProjectID(projectID int, ownerID int) ([]entities.Task, error)
	Update(id int, ownerID int, task entities.Task) (entities.Task, error)
	Delete(id int, ownerID int) error
}

type taskService struct {
	taskRepo    repositories.TaskRepository
	redisClient *redis.Client
	hub         *appwebsocket.Hub
}

func NewTaskService(
	taskRepo repositories.TaskRepository,
	redisClient *redis.Client,
	hub *appwebsocket.Hub,
) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		redisClient: redisClient,
		hub:         hub,
	}
}

func (s *taskService) Create(task entities.Task, ownerID int) (entities.Task, error) {
	newTask, err := s.taskRepo.Create(task, ownerID)

	if err != nil {
		return entities.Task{}, err
	}

	listCacheKey := fmt.Sprintf("tasks:user:%d", ownerID)
	projectCacheKey := fmt.Sprintf("tasks:project:%d:user:%d", newTask.ProjectID, ownerID)
	s.redisClient.Del(cache.Ctx, listCacheKey)
	s.redisClient.Del(cache.Ctx, projectCacheKey)

	if newTask.AssigneeID != nil {
		err := jobs.PushNotificationJob(s.redisClient, jobs.NotificationJob{
			Type:       "task.assigned",
			TaskID:     newTask.ID,
			AssigneeID: *newTask.AssigneeID,
			RetryCount: 0,
			MaxRetry:   3,
		})

		if err != nil {
			fmt.Println("Failed to push notification job:", err)
		}
	}

	// ← đổi BroadcastJSON sang SendToUser
	s.hub.SendToUser(
		ownerID,
		"task.created",
		newTask,
		"Task created successfully",
	)

	return newTask, nil
}


func (s *taskService) GetAllByOwnerID(ownerID int) ([]entities.Task, error) {
	cacheKey := fmt.Sprintf("tasks:user:%d", ownerID)

	cachedData, err := s.redisClient.Get(cache.Ctx, cacheKey).Result()

	if err == nil {
		var tasks []entities.Task

		err = json.Unmarshal([]byte(cachedData), &tasks)

		if err == nil {
			return tasks, nil
		}
	}

	tasks, err := s.taskRepo.GetAllByOwnerID(ownerID)

	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(tasks)

	if err == nil {
		s.redisClient.Set(cache.Ctx, cacheKey, data, 5*time.Minute)
	}

	return tasks, nil
}

func (s *taskService) GetByIDAndOwnerID(id int, ownerID int) (entities.Task, error) {
	cacheKey := fmt.Sprintf("task:%d:user:%d", id, ownerID)

	cachedData, err := s.redisClient.Get(cache.Ctx, cacheKey).Result()

	if err == nil {
		var task entities.Task

		err = json.Unmarshal([]byte(cachedData), &task)

		if err == nil {
			return task, nil
		}
	}

	task, err := s.taskRepo.GetByIDAndOwnerID(id, ownerID)

	if err != nil {
		return entities.Task{}, err
	}

	data, err := json.Marshal(task)

	if err == nil {
		s.redisClient.Set(cache.Ctx, cacheKey, data, 5*time.Minute)
	}

	return task, nil
}

func (s *taskService) Update(id int, ownerID int, task entities.Task) (entities.Task, error) {
	updatedTask, err := s.taskRepo.Update(id, ownerID, task)

	if err != nil {
		return entities.Task{}, err
	}

	listCacheKey := fmt.Sprintf("tasks:user:%d", ownerID)
	detailCacheKey := fmt.Sprintf("task:%d:user:%d", id, ownerID)
	projectCacheKey := fmt.Sprintf("tasks:project:%d:user:%d", updatedTask.ProjectID, ownerID)

	s.redisClient.Del(cache.Ctx, listCacheKey)
	s.redisClient.Del(cache.Ctx, detailCacheKey)
	s.redisClient.Del(cache.Ctx, projectCacheKey)

	// ← đổi BroadcastJSON sang SendToUser
	s.hub.SendToUser(
		ownerID,
		"task.status_updated",
		updatedTask,
		"Task updated successfully",
	)

	return updatedTask, nil
}

func (s *taskService) Delete(id int, ownerID int) error {
	task, err := s.taskRepo.GetByIDAndOwnerID(id, ownerID)
	if err != nil {
		return err
	}

	err = s.taskRepo.Delete(id, ownerID)

	if err != nil {
		return err
	}

	listCacheKey := fmt.Sprintf("tasks:user:%d", ownerID)
	detailCacheKey := fmt.Sprintf("task:%d:user:%d", id, ownerID)
	projectCacheKey := fmt.Sprintf("tasks:project:%d:user:%d", task.ProjectID, ownerID)

	s.redisClient.Del(cache.Ctx, listCacheKey)
	s.redisClient.Del(cache.Ctx, detailCacheKey)
	s.redisClient.Del(cache.Ctx, projectCacheKey)

	// ← đổi BroadcastJSON sang SendToUser
	s.hub.SendToUser(
		ownerID,
		"task.deleted",
		map[string]int{"id": id},
		"Task deleted successfully",
	)

	return nil
}


func (s *taskService) GetAllByProjectID(projectID int, ownerID int) ([]entities.Task, error) {
	cacheKey := fmt.Sprintf("tasks:project:%d:user:%d", projectID, ownerID)

	cachedData, err := s.redisClient.Get(cache.Ctx, cacheKey).Result()

	if err == nil {
		var tasks []entities.Task

		err = json.Unmarshal([]byte(cachedData), &tasks)

		if err == nil {
			return tasks, nil
		}
	}

	tasks, err := s.taskRepo.GetAllByProjectID(projectID, ownerID)

	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(tasks)

	if err == nil {
		s.redisClient.Set(cache.Ctx, cacheKey, data, 5*time.Minute)
	}

	return tasks, nil
}

func (s *taskService) GetAllWithFilter(ownerID int, params dto.TaskQueryParams) ([]entities.Task, int, error) {
	return s.taskRepo.GetAllWithFilter(ownerID, params)
}
