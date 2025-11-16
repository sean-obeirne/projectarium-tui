package tui

import (
	"github.com/sean-obeirne/projectarium-tui/internal/api"
)

// KanbanBoard represents the kanban board view
type KanbanBoard struct {
	columns             []ProjectColumn
	selectedCol         int
	selectedProject     int
	desiredProject      int
	scrollOffset        []int // scroll offset for each column
	desiredScrollOffset []int // desired scroll offset for each column
	width               int
	height              int
}

// ProjectColumn represents a column containing projects
type ProjectColumn struct {
	Name     string
	Projects []api.Project
}

// NewKanbanBoard creates a new kanban board with projects organized by status
func NewKanbanBoard(projects []api.Project) *KanbanBoard {
	// Organize projects into columns based on status
	columns := []ProjectColumn{
		{Name: "Ready", Projects: []api.Project{}},
		{Name: "In Progress", Projects: []api.Project{}},
		{Name: "Finished", Projects: []api.Project{}},
	}

	for _, project := range projects {
		switch project.Status {
		case "ready", "Ready", "READY", "":
			columns[0].Projects = append(columns[0].Projects, project)
		case "in_progress", "In Progress", "IN_PROGRESS", "in progress":
			columns[1].Projects = append(columns[1].Projects, project)
		case "finished", "Finished", "FINISHED", "done", "Done", "DONE":
			columns[2].Projects = append(columns[2].Projects, project)
		default:
			// Default to Ready if status is unknown
			columns[0].Projects = append(columns[0].Projects, project)
		}
	}

	kb := &KanbanBoard{
		columns:             columns,
		selectedCol:         0,
		selectedProject:     0,
		desiredProject:      0,
		scrollOffset:        make([]int, len(columns)),
		desiredScrollOffset: make([]int, len(columns)),
	}

	// Find first non-empty column to start with
	for i, col := range columns {
		if len(col.Projects) > 0 {
			kb.selectedCol = i
			break
		}
	}

	return kb
}

// SetSize sets the board dimensions
func (b *KanbanBoard) SetSize(width, height int) {
	b.width = width
	b.height = height

	// Reset scroll offsets when terminal size changes to ensure we start at the top
	for i := range b.scrollOffset {
		b.scrollOffset[i] = 0
	}

	// Reset selection to top if we're scrolled past visible area
	if len(b.columns) > 0 && len(b.columns[b.selectedCol].Projects) > 0 {
		b.selectedProject = 0
	}
}
