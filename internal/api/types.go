package api

// Project represents a kanban project
// Matches the backend API model in projectarium-v2/internal/models/project.go
type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	File        string `json:"file"`
	Priority    int    `json:"priority"`
	Status      string `json:"status"`
	Language    string `json:"language"`
}

// Todo represents a task/todo item in a kanban board
// Matches the backend API model in projectarium-v2/internal/models/todo.go
type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Deleted     bool   `json:"deleted"`
	ProjectID   *int   `json:"project_id"`
}

// Column represents a kanban column/status for display in the TUI
type Column struct {
	Name  string
	Todos []Todo
}
