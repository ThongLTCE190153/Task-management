package mappers

import (
	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
)

func ToTaskEntity(request dto.CreateTaskRequest) entities.Task {
	return entities.Task{
		ProjectID:   request.ProjectID,
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
		AssigneeID:  request.AssigneeID,
	}
}

func ToUpdateTaskEntity(request dto.UpdateTaskRequest) entities.Task {
	return entities.Task{
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
		AssigneeID:  request.AssigneeID,
	}
}

func ToTaskResponse(task entities.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:          task.ID,
		ProjectID:   task.ProjectID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		AssigneeID:  task.AssigneeID,
		CreatedAt:   task.CreatedAt,
	}
}

func ToTaskResponses(tasks []entities.Task) []dto.TaskResponse {
	var responses []dto.TaskResponse

	for _, task := range tasks {
		responses = append(responses, ToTaskResponse(task))
	}

	return responses
}
