package dto

import "time"

type CreateTaskRequest struct {
	ProjectID   int    `json:"project_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	AssigneeID  *int   `json:"assignee_id"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	AssigneeID  *int   `json:"assignee_id"`
}

type TaskResponse struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssigneeID  *int      `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskQueryParams struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	Status     string `form:"status"`
	AssigneeID *int   `form:"assignee_id"`
	ProjectID  *int   `form:"project_id"`
}

type PaginationResponse struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
