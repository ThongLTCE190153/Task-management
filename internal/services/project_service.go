package services

import (
	"encoding/json"
	"fmt"
	"time"

	"trithong.com/task-golang/internal/cache"
	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/repositories"
	appwebsocket "trithong.com/task-golang/internal/websocket"

	"github.com/redis/go-redis/v9"
)

type ProjectService interface {
	Create(project entities.Project) (entities.Project, error)
	GetAllByOwnerID(ownerID int) ([]entities.Project, error)
	GetAllByRole(userID int, role string) ([]entities.Project, error)
	GetAllWithPagination(ownerID int, params dto.ProjectQueryParams) ([]entities.Project, int, error)
	GetAllWithPaginationByRole(userID int, role string, params dto.ProjectQueryParams) ([]entities.Project, int, error)
	GetByIDAndOwnerID(id int, ownerID int) (entities.Project, error)
	Update(id int, ownerID int, project entities.Project) (entities.Project, error)
	Delete(id int, ownerID int) error
}

type projectService struct {
	projectRepo repositories.ProjectRepository
	redisClient *redis.Client
	hub         *appwebsocket.Hub
}

func NewProjectService(
	projectRepo repositories.ProjectRepository,
	redisClient *redis.Client,
	hub *appwebsocket.Hub,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		redisClient: redisClient,
		hub:         hub,
	}
}

func (s *projectService) Create(project entities.Project) (entities.Project, error) {
	newProject, err := s.projectRepo.Create(project)

	if err != nil {
		return entities.Project{}, err
	}

	cacheKey := fmt.Sprintf("projects:user:%d", project.OwnerID)
	s.redisClient.Del(cache.Ctx, cacheKey)

	// ← đổi BroadcastJSON sang SendToUser
	s.hub.SendToUser(
		project.OwnerID,
		"project.created",
		newProject,
		"Project created successfully",
	)

	return newProject, nil
}

func (s *projectService) GetAllByOwnerID(ownerID int) ([]entities.Project, error) {
	cacheKey := fmt.Sprintf("projects:user:%d", ownerID)

	cachedData, err := s.redisClient.Get(cache.Ctx, cacheKey).Result()

	if err == nil {
		var projects []entities.Project

		err = json.Unmarshal([]byte(cachedData), &projects)

		if err == nil {
			return projects, nil
		}
	}

	projects, err := s.projectRepo.GetAllByOwnerID(ownerID)

	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(projects)

	if err == nil {
		s.redisClient.Set(cache.Ctx, cacheKey, data, 5*time.Minute)
	}

	return projects, nil
}

func (s *projectService) GetAllByRole(userID int, role string) ([]entities.Project, error) {
	cacheKey := fmt.Sprintf("projects:user:%d:role:%s", userID, role)

	cachedData, err := s.redisClient.Get(cache.Ctx, cacheKey).Result()

	if err == nil {
		var projects []entities.Project
		err = json.Unmarshal([]byte(cachedData), &projects)
		if err == nil {
			return projects, nil
		}
	}

	projects, err := s.projectRepo.GetAllByRole(userID, role)

	if err != nil {
		return nil, err
	}

	cacheBytes, _ := json.Marshal(projects)
	s.redisClient.Set(cache.Ctx, cacheKey, cacheBytes, 10*time.Minute)

	return projects, nil
}

func (s *projectService) GetAllWithPaginationByRole(userID int, role string, params dto.ProjectQueryParams) ([]entities.Project, int, error) {
	cacheKey := fmt.Sprintf("projects:user:%d:role:%s:page:%d:limit:%d", userID, role, params.Page, params.Limit)

	cachedData, err := s.redisClient.Get(cache.Ctx, cacheKey).Result()

	if err == nil {
		var result struct {
			Projects []entities.Project `json:"projects"`
			Total    int                `json:"total"`
		}
		err = json.Unmarshal([]byte(cachedData), &result)
		if err == nil {
			return result.Projects, result.Total, nil
		}
	}

	projects, total, err := s.projectRepo.GetAllWithPaginationByRole(userID, role, params)

	if err != nil {
		return nil, 0, err
	}

	result := struct {
		Projects []entities.Project `json:"projects"`
		Total    int                `json:"total"`
	}{
		Projects: projects,
		Total:    total,
	}
	cacheBytes, _ := json.Marshal(result)
	s.redisClient.Set(cache.Ctx, cacheKey, cacheBytes, 10*time.Minute)

	return projects, total, nil
}

func (s *projectService) GetByIDAndOwnerID(id int, ownerID int) (entities.Project, error) {
	cacheKey := fmt.Sprintf("project:%d:user:%d", id, ownerID)

	cachedData, err := s.redisClient.Get(cache.Ctx, cacheKey).Result()

	if err == nil {
		var project entities.Project

		err = json.Unmarshal([]byte(cachedData), &project)

		if err == nil {
			return project, nil
		}
	}

	project, err := s.projectRepo.GetByIDAndOwnerID(id, ownerID)

	if err != nil {
		return entities.Project{}, err
	}

	data, err := json.Marshal(project)

	if err == nil {
		s.redisClient.Set(cache.Ctx, cacheKey, data, 5*time.Minute)
	}

	return project, nil
}

func (s *projectService) Update(id int, ownerID int, project entities.Project) (entities.Project, error) {
	updatedProject, err := s.projectRepo.Update(id, ownerID, project)

	if err != nil {
		return entities.Project{}, err
	}

	listCacheKey := fmt.Sprintf("projects:user:%d", ownerID)
	detailCacheKey := fmt.Sprintf("project:%d:user:%d", id, ownerID)

	s.redisClient.Del(cache.Ctx, listCacheKey)
	s.redisClient.Del(cache.Ctx, detailCacheKey)

	// ← đổi BroadcastJSON sang SendToUser
	s.hub.SendToUser(
		ownerID,
		"project.updated",
		updatedProject,
		"Project updated successfully",
	)

	return updatedProject, nil
}

func (s *projectService) Delete(id int, ownerID int) error {
	err := s.projectRepo.Delete(id, ownerID)

	if err != nil {
		return err
	}

	listCacheKey := fmt.Sprintf("projects:user:%d", ownerID)
	detailCacheKey := fmt.Sprintf("project:%d:user:%d", id, ownerID)

	s.redisClient.Del(cache.Ctx, listCacheKey)
	s.redisClient.Del(cache.Ctx, detailCacheKey)

	// ← đổi BroadcastJSON sang SendToUser
	s.hub.SendToUser(
		ownerID,
		"project.deleted",
		map[string]int{"id": id},
		"Project deleted successfully",
	)

	return nil
}

func (s *projectService) GetAllWithPagination(ownerID int, params dto.ProjectQueryParams) ([]entities.Project, int, error) {
	return s.projectRepo.GetAllWithPagination(ownerID, params)
}
