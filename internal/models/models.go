package models

import "time"

type GetTasksResponse struct {
	TaskID      int        `json:"task_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      int        `json:"status"`
	Priority    int        `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	DeadlineAt  *time.Time `json:"deadline_at,omitempty"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	DeadlineAt  *time.Time `json:"deadline_at"`
}

type GenricTaskResponse struct {
	TaskID  int64  `json:"taskid"`
	Message string `json:"message"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *int       `json:"status"`
	Priority    *int       `json:"priority"`
	DeadlineAt  *time.Time `json:"deadline_at"`
}
