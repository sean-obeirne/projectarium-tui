package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
)

// KanbanBoard represents the kanban board view
type KanbanBoard struct {
	project      *api.Project
	columns      []api.Column
	selectedCol  int
	selectedTodo int
	width        int
	height       int
}

// NewKanbanBoard creates a new kanban board
func NewKanbanBoard(project *api.Project, todos []api.Todo) *KanbanBoard {
	// Organize todos into columns based on priority
	// Since the backend doesn't have a "status" field on todos,
	// we'll organize by priority for now
	columns := []api.Column{
		{Name: "To Do", Todos: []api.Todo{}},
		{Name: "In Progress", Todos: []api.Todo{}},
		{Name: "Done", Todos: []api.Todo{}},
	}

	for _, todo := range todos {
		if todo.Deleted {
			continue // Skip deleted todos
		}
		switch {
		case todo.Priority == 0:
			columns[0].Todos = append(columns[0].Todos, todo)
		case todo.Priority == 1:
			columns[1].Todos = append(columns[1].Todos, todo)
		case todo.Priority >= 2:
			columns[2].Todos = append(columns[2].Todos, todo)
		}
	}

	return &KanbanBoard{
		project:      project,
		columns:      columns,
		selectedCol:  0,
		selectedTodo: 0,
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
				b.selectedTodo = 0
			}
		case "right", "l":
			if b.selectedCol < len(b.columns)-1 {
				b.selectedCol++
				b.selectedTodo = 0
			}
		case "up", "k":
			if b.selectedTodo > 0 {
				b.selectedTodo--
			}
		case "down", "j":
			if b.selectedTodo < len(b.columns[b.selectedCol].Todos)-1 {
				b.selectedTodo++
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

	selectedColumnHeaderStyle := columnHeaderStyle.
		Background(lipgloss.Color("99"))

	taskStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	selectedTaskStyle := taskStyle.
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
		header := headerStyle.Render(fmt.Sprintf("%s (%d)", col.Name, len(col.Todos)))

		// Todos
		todoViews := []string{}
		displayedTodos := 0
		maxTodos := (maxHeight - 2) / 4 // Approximate todos that fit

		for j, todo := range col.Todos {
			if displayedTodos >= maxTodos {
				todoViews = append(todoViews, lipgloss.NewStyle().
					Foreground(lipgloss.Color("241")).
					Render(fmt.Sprintf("  ... %d more", len(col.Todos)-displayedTodos)))
				break
			}

			style := taskStyle
			if i == b.selectedCol && j == b.selectedTodo {
				style = selectedTaskStyle
			}

			// Truncate todo description if too long
			description := todo.Description
			maxDescLen := colWidth - 4
			if len(description) > maxDescLen {
				description = description[:maxDescLen-3] + "..."
			}

			todoView := style.Width(colWidth - 2).Render(description)
			todoViews = append(todoViews, todoView)
			displayedTodos++
		}

		if len(todoViews) == 0 {
			todoViews = append(todoViews, lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(1).
				Render("No todos"))
		}

		columnContent := lipgloss.JoinVertical(lipgloss.Left, todoViews...)
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
