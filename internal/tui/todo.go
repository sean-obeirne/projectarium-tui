package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
)

// TodoInputMode represents the current input mode for todos
type TodoInputMode int

const (
	NormalMode TodoInputMode = iota
	AddingMode
	EditingMode
)

// TodoList represents the todo list view for a project
type TodoList struct {
	todos         []api.Todo
	selectedIndex int
	projectName   string
	projectID     int
	width         int
	height        int
	InputMode     TodoInputMode // Exported so model.go can check it
	textInput     textinput.Model
	editingTodoID int
}

// NewTodoList creates a new todo list view
func NewTodoList(todos []api.Todo, projectName string, projectID int) *TodoList {
	ti := textinput.New()
	ti.Placeholder = "Enter todo description..."
	ti.CharLimit = 200

	return &TodoList{
		todos:         todos,
		selectedIndex: 0,
		projectName:   projectName,
		projectID:     projectID,
		InputMode:     NormalMode,
		textInput:     ti,
	}
}

// SetSize sets the todo list dimensions
func (t *TodoList) SetSize(width, height int) {
	t.width = width
	t.height = height
}

// Update handles messages for the todo list
func (t TodoList) Update(msg tea.Msg) (TodoList, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input mode first - check for Enter/Esc BEFORE passing to textinput
		if t.InputMode == AddingMode || t.InputMode == EditingMode {
			switch msg.String() {
			case "enter":
				// Submit the todo
				description := t.textInput.Value()
				if description != "" {
					var cmd tea.Cmd
					if t.InputMode == AddingMode {
						cmd = createTodoCmd(description, 0, t.projectID)
					} else if t.InputMode == EditingMode {
						if t.selectedIndex < len(t.todos) {
							todo := t.todos[t.selectedIndex]
							cmd = updateTodoCmd(todo.ID, description, todo.Priority, todo.ProjectID)
						}
					}
					// Reset input mode and clear input after submitting
					t.InputMode = NormalMode
					t.textInput.SetValue("")
					t.textInput.Blur()
					return t, cmd
				}
				// Cancel if empty
				t.InputMode = NormalMode
				t.textInput.SetValue("")
				t.textInput.Blur()
				return t, nil
			case "esc":
				// Cancel input
				t.InputMode = NormalMode
				t.textInput.SetValue("")
				t.textInput.Blur()
				return t, nil
			}
			// Pass other keys to text input
			t.textInput, cmd = t.textInput.Update(msg)
			return t, cmd
		}

		// Normal mode key handling
		switch msg.String() {
		case "up", "k":
			if t.selectedIndex > 0 {
				t.selectedIndex--
			}
		case "down", "j":
			if t.selectedIndex < len(t.todos)-1 {
				t.selectedIndex++
			}
		case "a":
			// Start adding a new todo
			t.InputMode = AddingMode
			t.textInput.Focus()
			t.textInput.SetValue("")
			return t, textinput.Blink
		case "e":
			// Start editing selected todo
			if len(t.todos) > 0 && t.selectedIndex < len(t.todos) {
				t.InputMode = EditingMode
				t.editingTodoID = t.todos[t.selectedIndex].ID
				t.textInput.SetValue(t.todos[t.selectedIndex].Description)
				t.textInput.Focus()
				return t, textinput.Blink
			}
		case "d", "x":
			// Delete selected todo
			if len(t.todos) > 0 && t.selectedIndex < len(t.todos) {
				todoID := t.todos[t.selectedIndex].ID
				return t, deleteTodoCmd(todoID)
			}
		case "+", "=":
			// Increase priority
			if len(t.todos) > 0 && t.selectedIndex < len(t.todos) {
				todo := t.todos[t.selectedIndex]
				newPriority := todo.Priority + 1
				if newPriority > 3 {
					newPriority = 3
				}
				return t, updateTodoCmd(todo.ID, todo.Description, newPriority, todo.ProjectID)
			}
		case "-", "_":
			// Decrease priority
			if len(t.todos) > 0 && t.selectedIndex < len(t.todos) {
				todo := t.todos[t.selectedIndex]
				newPriority := todo.Priority - 1
				if newPriority < 0 {
					newPriority = 0
				}
				return t, updateTodoCmd(todo.ID, todo.Description, newPriority, todo.ProjectID)
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

	inputPromptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Bold(true).
		MarginTop(1)

	// Title
	title := titleStyle.Render(fmt.Sprintf("📝 Todos for: %s", t.projectName))

	// Header
	header := headerStyle.Render(fmt.Sprintf("\nTotal todos: %d\n", len(t.todos)))

	// Input prompt if in input mode
	var inputPrompt string
	if t.InputMode == AddingMode {
		inputPrompt = inputPromptStyle.Render("Add new todo:")
	} else if t.InputMode == EditingMode {
		inputPrompt = inputPromptStyle.Render("Edit todo:")
	}

	// Todos
	var todoViews []string
	if len(t.todos) == 0 {
		todoViews = append(todoViews, emptyStyle.Render("No todos yet! Press 'a' to add one."))
	} else {
		for i, todo := range t.todos {
			// Priority indicator
			var priorityIndicator string
			var priorityColor lipgloss.Color
			switch todo.Priority {
			case 0:
				priorityIndicator = "○" // Low
				priorityColor = lipgloss.Color("240")
			case 1:
				priorityIndicator = "◐" // Medium-Low
				priorityColor = lipgloss.Color("70")
			case 2:
				priorityIndicator = "◑" // Medium-High
				priorityColor = lipgloss.Color("214")
			default:
				priorityIndicator = "●" // High
				priorityColor = lipgloss.Color("196")
			}

			priorityStyle := lipgloss.NewStyle().Foreground(priorityColor)
			todoText := fmt.Sprintf("%s %s", priorityStyle.Render(priorityIndicator), todo.Description)

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
	var help string
	if t.InputMode == NormalMode {
		help = helpStyle.Render("↑/k ↓/j navigate • a add • e edit • d delete • +/- priority • esc close")
	} else {
		help = helpStyle.Render("enter submit • esc cancel")
	}

	// Build the view
	var sections []string
	sections = append(sections, title, header, todos)

	if t.InputMode != NormalMode {
		sections = append(sections, "", inputPrompt, t.textInput.View())
	}

	sections = append(sections, "", help)

	// Combine everything
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Message types for todo operations
type createTodoMsg struct {
	description string
	priority    int
	projectID   int
}

type updateTodoMsg struct {
	id          int
	description string
	priority    int
	projectID   *int
}

type deleteTodoMsg struct {
	id int
}

// Command functions
func createTodoCmd(description string, priority int, projectID int) tea.Cmd {
	return func() tea.Msg {
		return createTodoMsg{
			description: description,
			priority:    priority,
			projectID:   projectID,
		}
	}
}

func updateTodoCmd(id int, description string, priority int, projectID *int) tea.Cmd {
	return func() tea.Msg {
		return updateTodoMsg{
			id:          id,
			description: description,
			priority:    priority,
			projectID:   projectID,
		}
	}
}

func deleteTodoCmd(id int) tea.Cmd {
	return func() tea.Msg {
		return deleteTodoMsg{id: id}
	}
}
