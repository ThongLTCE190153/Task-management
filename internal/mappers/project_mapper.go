package mappers

import (
	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
)

func ToProjectEntity(request dto.CreateProjectRequest) entities.Project {
	return entities.Project{
		Name:        request.Name,
		Description: request.Description,
	}
}

func ToProjectResponse(project entities.Project) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
	}
}

func ToUpdateProjectEntity(request dto.UpdateProjectRequest) entities.Project {
	return entities.Project{
		Name:        request.Name,
		Description: request.Description,
	}
}
