package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
)

// KanbanBoard represents the kanban board view
type KanbanBoard struct {
	project       *api.Project
	columns       []api.Column
	selectedCol   int
	selectedTask  int
	width         int
	height        int
}

// NewKanbanBoard creates a new kanban board
func NewKanbanBoard(project *api.Project, tasks []api.Task) *KanbanBoard {
	// Organize tasks into columns
	columns := []api.Column{
		{Name: "To Do", Tasks: []api.Task{}},
		{Name: "In Progress", Tasks: []api.Task{}},
		{Name: "Done", Tasks: []api.Task{}},
	}

	for _, task := range tasks {
		switch task.Status {
		case "todo", "":
			columns[0].Tasks = append(columns[0].Tasks, task)
		case "in_progress":
			columns[1].Tasks = append(columns[1].Tasks, task)
		case "done":
			columns[2].Tasks = append(columns[2].Tasks, task)
		}
	}

	return &KanbanBoard{
		project:      project,
		columns:      columns,
		selectedCol:  0,
		selectedTask: 0,
	}
}

// SetSize sets the board dimensions
func (b *KanbanBoard) SetSize(width, height int) {
	b.width = width
	b.height = height
}

// Update handles messages
func (b KanbanBoard) Update(msg tea.Msg) (KanbanBoard, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if b.selectedCol > 0 {
				b.selectedCol--
				b.selectedTask = 0
			}
		case "right", "l":
			if b.selectedCol < len(b.columns)-1 {
				b.selectedCol++
				b.selectedTask = 0
			}
		case "up", "k":
			if b.selectedTask > 0 {
				b.selectedTask--
			}
		case "down", "j":
			if b.selectedTask < len(b.columns[b.selectedCol].Tasks)-1 {
				b.selectedTask++
			}
		}
	}
	return b, nil
}

// View renders the kanban board
func (b *KanbanBoard) View() string {
	if b.width == 0 {
		return "Initializing..."
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		MarginBottom(1)

	columnHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("63"))

	selectedColumnHeaderStyle := columnHeaderStyle.Copy().
		Background(lipgloss.Color("99"))

	taskStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	selectedTaskStyle := taskStyle.Copy().
		BorderForeground(lipgloss.Color("63")).
		Foreground(lipgloss.Color("63"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	// Title
	title := titleStyle.Render(fmt.Sprintf("📋 %s", b.project.Name))

	// Calculate column width
	colWidth := (b.width - 4) / len(b.columns)
	if colWidth < 20 {
		colWidth = 20
	}

	// Build columns
	columnViews := make([]string, len(b.columns))
	maxHeight := b.height - 8 // Reserve space for title and help

	for i, col := range b.columns {
		// Column header
		headerStyle := columnHeaderStyle
		if i == b.selectedCol {
			headerStyle = selectedColumnHeaderStyle
		}
		header := headerStyle.Render(fmt.Sprintf("%s (%d)", col.Name, len(col.Tasks)))

		// Tasks
		taskViews := []string{}
		displayedTasks := 0
		maxTasks := (maxHeight - 2) / 4 // Approximate tasks that fit

		for j, task := range col.Tasks {
			if displayedTasks >= maxTasks {
				taskViews = append(taskViews, lipgloss.NewStyle().
					Foreground(lipgloss.Color("241")).
					Render(fmt.Sprintf("  ... %d more", len(col.Tasks)-displayedTasks)))
				break
			}

			style := taskStyle
			if i == b.selectedCol && j == b.selectedTask {
				style = selectedTaskStyle
			}

			// Truncate task title if too long
			title := task.Title
			maxTitleLen := colWidth - 4
			if len(title) > maxTitleLen {
				title = title[:maxTitleLen-3] + "..."
			}

			taskView := style.Width(colWidth - 2).Render(title)
			taskViews = append(taskViews, taskView)
			displayedTasks++
		}

		if len(taskViews) == 0 {
			taskViews = append(taskViews, lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(1).
				Render("No tasks"))
		}

		columnContent := lipgloss.JoinVertical(lipgloss.Left, taskViews...)
		columnViews[i] = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			columnContent,
		)
	}

	// Join columns horizontally
	board := lipgloss.JoinHorizontal(lipgloss.Top, columnViews...)

	// Help text
	help := helpStyle.Render("  ←/h → /l move columns • ↑/k ↓/j move tasks • esc back • r refresh • q quit")

	// Combine everything
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		board,
		help,
	)
}
