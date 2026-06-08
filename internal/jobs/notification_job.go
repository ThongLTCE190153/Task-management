package jobs

type NotificationJob struct {
	Type       string `json:"type"`
	TaskID     int    `json:"task_id"`
	AssigneeID int    `json:"assignee_id"`
	RetryCount int    `json:"retry_count"`
	MaxRetry   int    `json:"max_retry"`
}
