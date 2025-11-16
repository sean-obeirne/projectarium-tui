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
	apiClient   *api.Client
	config      *config.Config
	viewMode    ViewMode
	projects    []api.Project
	kanbanBoard *KanbanBoard
	width       int
	height      int
	err         error
	loading     bool
	keys        keyMap
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
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
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
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Refresh):
			if m.viewMode == KanbanBoardView {
				m.loading = true
				m.viewMode = LoadingView
				return m, m.loadProjects
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.kanbanBoard != nil {
			m.kanbanBoard.SetSize(msg.Width, msg.Height)
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
	}

	// Update the kanban board view
	var cmd tea.Cmd
	if m.viewMode == KanbanBoardView && m.kanbanBoard != nil {
		*m.kanbanBoard, cmd = m.kanbanBoard.Update(msg)
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
		return m.kanbanBoardView()
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
	help := helpStyle.Render("\nPress 'r' to retry, 'q' to quit")
	return errMsg + help
}

// Messages

type projectsLoadedMsg struct {
	projects []api.Project
	err      error
}

// Commands

func (m Model) loadProjects() tea.Msg {
	projects, err := m.apiClient.GetProjects()
	return projectsLoadedMsg{projects: projects, err: err}
}
