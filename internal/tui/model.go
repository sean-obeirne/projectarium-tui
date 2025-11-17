package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
	"github.com/sean-obeirne/projectarium-tui/internal/config"
)

// ViewMode represents the current view
type ViewMode int

const (
	KanbanBoardView ViewMode = iota
	LoadingView
	ErrorView
)

// Model is the main Bubble Tea model
type Model struct {
	apiClient      *api.Client
	config         *config.Config
	viewMode       ViewMode
	projects       []api.Project
	kanbanBoard    *KanbanBoard
	todoList       *TodoList
	showTodoList   bool // Whether to show todo list overlay
	currentProject *api.Project
	width          int
	height         int
	err            error
	loading        bool
	keys           keyMap
}

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Enter   key.Binding
	Back    key.Binding
	Quit    key.Binding
	Refresh key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "l"),
		key.WithHelp("↑/l", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "k"),
		key.WithHelp("↓/k", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "j"),
		key.WithHelp("←/j", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", ";"),
		key.WithHelp("→/;", "move right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "refresh"),
	),
}

// NewModel creates a new TUI model
func NewModel() Model {
	cfg := config.Load()
	client := api.NewClient(cfg.APIBaseURL)

	return Model{
		apiClient: client,
		config:    cfg,
		viewMode:  LoadingView,
		keys:      keys,
		loading:   true,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.loadProjects
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If todo list is showing and in input mode, let it handle keys first
	if m.showTodoList && m.todoList != nil && (m.todoList.InputMode == AddingMode || m.todoList.InputMode == EditingMode) {
		var cmd tea.Cmd
		*m.todoList, cmd = m.todoList.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			// If todo list is showing, close it instead of quitting
			if m.showTodoList {
				m.showTodoList = false
				m.currentProject = nil
				m.todoList = nil
				return m, nil
			}
			return m, tea.Quit
		case key.Matches(msg, m.keys.Refresh):
			if m.viewMode == KanbanBoardView {
				m.loading = true
				m.viewMode = LoadingView
				return m, m.loadProjects
			}
		case key.Matches(msg, m.keys.Enter):
			// Toggle todo list for selected project
			if m.viewMode == KanbanBoardView && m.kanbanBoard != nil {
				if m.showTodoList {
					// Close todo list if already showing
					m.showTodoList = false
					m.currentProject = nil
					m.todoList = nil
					return m, nil
				} else {
					// Open todo list for selected project
					if project := m.kanbanBoard.GetSelectedProject(); project != nil {
						m.currentProject = project
						return m, m.loadTodos
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.kanbanBoard != nil {
			m.kanbanBoard.SetSize(msg.Width, msg.Height)
		}
		if m.todoList != nil {
			m.todoList.SetSize(msg.Width, msg.Height)
		}

	case projectsLoadedMsg:
		m.loading = false
		m.projects = msg.projects
		m.err = msg.err
		if m.err != nil {
			m.viewMode = ErrorView
			return m, nil
		}

		// Create kanban board with all projects
		m.kanbanBoard = NewKanbanBoard(m.projects)
		m.kanbanBoard.SetSize(m.width, m.height)
		m.viewMode = KanbanBoardView
		return m, nil

	case todosLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.viewMode = ErrorView
			return m, nil
		}

		projectName := "Unknown Project"
		projectID := 0
		if m.currentProject != nil {
			projectName = m.currentProject.Name
			projectID = m.currentProject.ID
		}

		m.todoList = NewTodoList(msg.todos, projectName, projectID)
		m.todoList.SetSize(m.width, m.height)
		m.showTodoList = true
		return m, nil

	case progressProjectMsg:
		return m, m.updateProjectStatus(msg.projectID, msg.status)

	case regressProjectMsg:
		return m, m.updateProjectStatus(msg.projectID, msg.status)

	case updatePriorityMsg:
		return m, m.updateProjectPriority(msg.projectID, msg.priority)

	case projectStatusUpdatedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.viewMode = ErrorView
			return m, nil
		}
		// Update the project in the kanban board
		if msg.project != nil && m.kanbanBoard != nil {
			m.kanbanBoard.UpdateProjectInBoard(*msg.project)
		}
		return m, nil

	case projectPriorityUpdatedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.viewMode = ErrorView
			return m, nil
		}
		// Update the project in the kanban board by reloading
		// (Priority changes may need reordering)
		if msg.project != nil && m.kanbanBoard != nil {
			m.kanbanBoard.UpdateProjectInBoard(*msg.project)
		}
		return m, nil

	case createTodoMsg:
		// User wants to create a new todo
		return m, m.createTodo(msg.description, msg.priority, msg.projectID)

	case updateTodoMsg:
		// User wants to update a todo
		return m, m.updateTodo(msg.id, msg.description, msg.priority, msg.projectID)

	case deleteTodoMsg:
		// User wants to delete a todo
		return m, m.deleteTodo(msg.id)

	case todoCreatedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.viewMode = ErrorView
			return m, nil
		}
		// Reload todos for the current project
		if m.currentProject != nil {
			return m, m.loadTodos
		}
		return m, nil

	case todoUpdatedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.viewMode = ErrorView
			return m, nil
		}
		// Reload todos for the current project
		if m.currentProject != nil {
			return m, m.loadTodos
		}
		return m, nil

	case todoDeletedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.viewMode = ErrorView
			return m, nil
		}
		// Reload todos for the current project
		if m.currentProject != nil {
			return m, m.loadTodos
		}
		return m, nil
	}

	// Update the appropriate view
	var cmd tea.Cmd
	if m.viewMode == KanbanBoardView && m.kanbanBoard != nil {
		// Only handle kanban updates if todo list is not showing
		if !m.showTodoList {
			*m.kanbanBoard, cmd = m.kanbanBoard.Update(msg)
			return m, cmd
		}
	}

	// If todo list is showing, let it handle updates
	if m.showTodoList && m.todoList != nil {
		*m.todoList, cmd = m.todoList.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.loading {
		return m.loadingView()
	}

	switch m.viewMode {
	case KanbanBoardView:
		board := m.kanbanBoardView()
		// Overlay todo list if showing
		if m.showTodoList && m.todoList != nil && m.currentProject != nil {
			// Calculate dimensions
			projectCardWidth := 40
			todoListWidth := int(float64(m.width) * 0.5)
			if todoListWidth < 50 {
				todoListWidth = 50
			}

			modalHeight := int(float64(m.height) * 0.7)
			if modalHeight < 20 {
				modalHeight = 20
			}

			// Render the project card
			projectCard := m.kanbanBoard.RenderProjectCard(m.currentProject, projectCardWidth)

			// Style the todo list
			todoStyle := lipgloss.NewStyle().
				Width(todoListWidth).
				Height(modalHeight).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2)

			todoView := todoStyle.Render(m.todoList.View())

			// Join project card and todo list horizontally
			combined := lipgloss.JoinHorizontal(
				lipgloss.Top,
				projectCard,
				"  ", // spacing
				todoView,
			)

			// Center the combined view
			return lipgloss.Place(
				m.width,
				m.height,
				lipgloss.Center,
				lipgloss.Center,
				combined,
			)
		}
		return board
	case ErrorView:
		return m.errorView()
	default:
		return m.loadingView()
	}
}

func (m Model) kanbanBoardView() string {
	if m.kanbanBoard == nil {
		return "Loading board..."
	}
	return m.kanbanBoard.View()
}

func (m Model) loadingView() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		MarginTop(2).
		MarginLeft(2)
	return style.Render("Loading projects...")
}

func (m Model) errorView() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		MarginTop(2).
		MarginLeft(2)
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1).
		MarginLeft(2)

	errMsg := errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	help := helpStyle.Render("\nPress 'R' to retry, 'q' to quit")
	return errMsg + help
}

// Messages

type projectsLoadedMsg struct {
	projects []api.Project
	err      error
}

type todosLoadedMsg struct {
	todos []api.Todo
	err   error
}

type progressProjectMsg struct {
	projectID int
	status    string
}

type regressProjectMsg struct {
	projectID int
	status    string
}

type updatePriorityMsg struct {
	projectID int
	priority  int
}

type projectStatusUpdatedMsg struct {
	project *api.Project
	err     error
}

type projectPriorityUpdatedMsg struct {
	project *api.Project
	err     error
}

type todoCreatedMsg struct {
	todo *api.Todo
	err  error
}

type todoUpdatedMsg struct {
	todo *api.Todo
	err  error
}

type todoDeletedMsg struct {
	err error
}

// Commands

func (m Model) loadProjects() tea.Msg {
	projects, err := m.apiClient.GetProjects()
	return projectsLoadedMsg{projects: projects, err: err}
}

func (m Model) loadTodos() tea.Msg {
	if m.currentProject == nil {
		return todosLoadedMsg{todos: []api.Todo{}, err: nil}
	}
	todos, err := m.apiClient.GetTodosByProject(m.currentProject.ID)
	return todosLoadedMsg{todos: todos, err: err}
}

func (m Model) updateProjectStatus(projectID int, status string) tea.Cmd {
	return func() tea.Msg {
		project, err := m.apiClient.UpdateProjectStatus(projectID, status)
		return projectStatusUpdatedMsg{project: project, err: err}
	}
}

func (m Model) updateProjectPriority(projectID int, priority int) tea.Cmd {
	return func() tea.Msg {
		project, err := m.apiClient.UpdateProjectPriority(projectID, priority)
		return projectPriorityUpdatedMsg{project: project, err: err}
	}
}

func (m Model) createTodo(description string, priority int, projectID int) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.apiClient.CreateTodo(description, priority, &projectID)
		return todoCreatedMsg{todo: todo, err: err}
	}
}

func (m Model) updateTodo(id int, description string, priority int, projectID *int) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.apiClient.UpdateTodo(id, description, priority, projectID)
		return todoUpdatedMsg{todo: todo, err: err}
	}
}

func (m Model) deleteTodo(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.apiClient.DeleteTodo(id)
		return todoDeletedMsg{err: err}
	}
}
