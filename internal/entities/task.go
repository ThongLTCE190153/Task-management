package entities

import "time"

type Task struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssigneeID  *int      `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
}
