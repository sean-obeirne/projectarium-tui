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

// GetSelectedProject returns the currently selected project, or nil if none
func (b *KanbanBoard) GetSelectedProject() *api.Project {
	if len(b.columns) == 0 || b.selectedCol >= len(b.columns) {
		return nil
	}
	col := b.columns[b.selectedCol]
	if b.selectedProject >= len(col.Projects) {
		return nil
	}
	return &col.Projects[b.selectedProject]
}

// GetNextStatus returns the status name for progressing a project forward
func (b *KanbanBoard) GetNextStatus() string {
	if b.selectedCol < len(b.columns)-1 {
		return b.statusNameToAPIStatus(b.columns[b.selectedCol+1].Name)
	}
	return "" // Already at the end
}

// GetPrevStatus returns the status name for regressing a project backward
func (b *KanbanBoard) GetPrevStatus() string {
	if b.selectedCol > 0 {
		return b.statusNameToAPIStatus(b.columns[b.selectedCol-1].Name)
	}
	return "" // Already at the start
}

// statusNameToAPIStatus converts display name to API status value
func (b *KanbanBoard) statusNameToAPIStatus(displayName string) string {
	switch displayName {
	case "Ready":
		return "ready"
	case "In Progress":
		return "in_progress"
	case "Finished":
		return "finished"
	default:
		return "ready"
	}
}

// UpdateProjectInBoard updates a project in the board after an API change
func (b *KanbanBoard) UpdateProjectInBoard(updatedProject api.Project) {
	// Remove the project from its current column
	for colIdx := range b.columns {
		for projIdx, proj := range b.columns[colIdx].Projects {
			if proj.ID == updatedProject.ID {
				// Remove from current column
				b.columns[colIdx].Projects = append(
					b.columns[colIdx].Projects[:projIdx],
					b.columns[colIdx].Projects[projIdx+1:]...,
				)

				// Find the new column for the project
				newColIdx := b.getColumnIndexForStatus(updatedProject.Status)

				// Insert into new column in sorted order (priority DESC, then name ASC)
				inserted := false
				for i, proj := range b.columns[newColIdx].Projects {
					// Insert before this project if:
					// - Our priority is higher, OR
					// - Same priority but our name comes first alphabetically
					if updatedProject.Priority > proj.Priority ||
						(updatedProject.Priority == proj.Priority && updatedProject.Name < proj.Name) {
						// Insert at position i
						b.columns[newColIdx].Projects = append(
							b.columns[newColIdx].Projects[:i],
							append([]api.Project{updatedProject}, b.columns[newColIdx].Projects[i:]...)...,
						)
						b.selectedProject = i
						b.desiredProject = i
						inserted = true
						break
					}
				}

				// If not inserted yet, append to end
				if !inserted {
					b.columns[newColIdx].Projects = append(b.columns[newColIdx].Projects, updatedProject)
					b.selectedProject = len(b.columns[newColIdx].Projects) - 1
					b.desiredProject = b.selectedProject
				}

				// Update selection to follow the project
				b.selectedCol = newColIdx

				// Ensure the selected project is visible
				b.scrollOffset[b.selectedCol] = min(b.selectedProject, b.scrollOffset[b.selectedCol])
				b.desiredScrollOffset[b.selectedCol] = b.scrollOffset[b.selectedCol]

				return
			}
		}
	}
}

// getColumnIndexForStatus returns the column index for a given status
func (b *KanbanBoard) getColumnIndexForStatus(status string) int {
	switch status {
	case "ready", "Ready", "READY", "":
		return 0
	case "in_progress", "In Progress", "IN_PROGRESS", "in progress":
		return 1
	case "finished", "Finished", "FINISHED", "done", "Done", "DONE":
		return 2
	default:
		return 0
	}
}
