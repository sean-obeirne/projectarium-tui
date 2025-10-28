package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

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
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	// Sample data
	projects := []Project{
		{
			ID:          "1",
			Name:        "Website Redesign",
			Description: "Redesign the company website with modern UI",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "Mobile App",
			Description: "Develop mobile application for iOS and Android",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "3",
			Name:        "API Migration",
			Description: "Migrate REST API to GraphQL",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}

	tasksMap := map[string][]Task{
		"1": {
			{ID: "1", ProjectID: "1", Title: "Design homepage mockup", Description: "Create mockup in Figma", Status: "done", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "2", ProjectID: "1", Title: "Implement responsive layout", Description: "Make site responsive", Status: "in_progress", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "3", ProjectID: "1", Title: "Add contact form", Description: "Create contact form", Status: "todo", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "4", ProjectID: "1", Title: "SEO optimization", Description: "Optimize for search engines", Status: "todo", Priority: "low", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		"2": {
			{ID: "5", ProjectID: "2", Title: "Setup React Native", Description: "Initialize RN project", Status: "done", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "6", ProjectID: "2", Title: "Design app screens", Description: "Create UI designs", Status: "in_progress", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "7", ProjectID: "2", Title: "Implement authentication", Description: "Add login/signup", Status: "todo", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "8", ProjectID: "2", Title: "App store submission", Description: "Submit to stores", Status: "todo", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		"3": {
			{ID: "9", ProjectID: "3", Title: "Define GraphQL schema", Description: "Design schema", Status: "done", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "10", ProjectID: "3", Title: "Setup Apollo Server", Description: "Configure server", Status: "in_progress", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "11", ProjectID: "3", Title: "Migrate endpoints", Description: "Convert REST to GraphQL", Status: "in_progress", Priority: "high", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "12", ProjectID: "3", Title: "Update documentation", Description: "Document new API", Status: "todo", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "13", ProjectID: "3", Title: "Write integration tests", Description: "Test API endpoints", Status: "todo", Priority: "low", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	// API handlers
	http.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	})

	http.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		// Remove the base path
		path := strings.TrimPrefix(r.URL.Path, "/api/projects/")
		parts := strings.Split(path, "/")

		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		projectID := parts[0]

		// Check if it's a tasks request
		if len(parts) > 1 && parts[1] == "tasks" {
			w.Header().Set("Content-Type", "application/json")
			tasks, ok := tasksMap[projectID]
			if !ok {
				tasks = []Task{}
			}
			json.NewEncoder(w).Encode(tasks)
			return
		}

		// Single project request
		w.Header().Set("Content-Type", "application/json")
		for _, p := range projects {
			if p.ID == projectID {
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		http.Error(w, "Project not found", http.StatusNotFound)
	})

	addr := ":8888"
	fmt.Printf("Mock projectarium-v2 API server starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET /api/projects")
	fmt.Println("  GET /api/projects/:id")
	fmt.Println("  GET /api/projects/:id/tasks")
	fmt.Println("\nPress Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(addr, nil))
}
