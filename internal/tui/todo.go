package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
)

// TodoList represents the todo list view for a project
type TodoList struct {
	todos         []api.Todo
	selectedIndex int
	projectName   string
	width         int
	height        int
}

// NewTodoList creates a new todo list view
func NewTodoList(todos []api.Todo, projectName string) *TodoList {
	return &TodoList{
		todos:         todos,
		selectedIndex: 0,
		projectName:   projectName,
	}
}

// SetSize sets the todo list dimensions
func (t *TodoList) SetSize(width, height int) {
	t.width = width
	t.height = height
}

// Update handles messages for the todo list
func (t TodoList) Update(msg tea.Msg) (TodoList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.selectedIndex > 0 {
				t.selectedIndex--
			}
		case "down", "j":
			if t.selectedIndex < len(t.todos)-1 {
				t.selectedIndex++
			}
		}
	}
	return t, nil
}

// View renders the todo list
func (t *TodoList) View() string {
	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214"))

	todoItemStyle := lipgloss.NewStyle().
		Padding(0, 1)

	selectedTodoStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("51")).
		Bold(true)

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginLeft(2).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	// Title
	title := titleStyle.Render(fmt.Sprintf("📝 Todos for: %s", t.projectName))

	// Header
	header := headerStyle.Render(fmt.Sprintf("\nTotal todos: %d\n", len(t.todos)))

	// Todos
	var todoViews []string
	if len(t.todos) == 0 {
		todoViews = append(todoViews, emptyStyle.Render("No todos yet!"))
	} else {
		for i, todo := range t.todos {
			// Priority indicator
			priorityIndicator := ""
			switch todo.Priority {
			case 0:
				priorityIndicator = "○" // Low
			case 1:
				priorityIndicator = "◐" // Medium-Low
			case 2:
				priorityIndicator = "◑" // Medium-High
			default:
				priorityIndicator = "●" // High
			}

			todoText := fmt.Sprintf("%s %s", priorityIndicator, todo.Description)

			style := todoItemStyle
			if i == t.selectedIndex {
				todoText = "▸ " + todoText
				style = selectedTodoStyle
			} else {
				todoText = "  " + todoText
			}

			todoViews = append(todoViews, style.Render(todoText))
		}
	}

	todos := lipgloss.JoinVertical(lipgloss.Left, todoViews...)

	// Help text
	help := helpStyle.Render("↑/k ↓/j navigate • esc/backspace close")

	// Combine everything
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		header,
		todos,
		"",
		help,
	)
}
