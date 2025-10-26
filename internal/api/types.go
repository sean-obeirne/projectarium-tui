package api

import "time"

// Project represents a kanban project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Task represents a task in a kanban board
type Task struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // todo, in_progress, done
	Priority    string    `json:"priority"` // low, medium, high
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Column represents a kanban column/status
type Column struct {
	Name  string
	Tasks []Task
}
