package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
	"github.com/sean-obeirne/projectarium-tui/internal/config"
)

// ViewMode represents the current view
type ViewMode int

const (
	ProjectListView ViewMode = iota
	KanbanBoardView
	LoadingView
	ErrorView
)

// Model is the main Bubble Tea model
type Model struct {
	apiClient      *api.Client
	config         *config.Config
	viewMode       ViewMode
	projects       []api.Project
	selectedProject *api.Project
	tasks          []api.Task
	projectList    list.Model
	kanbanBoard    *KanbanBoard
	cursor         int
	width          int
	height         int
	err            error
	loading        bool
	keys           keyMap
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
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

	// Create list for projects
	delegate := list.NewDefaultDelegate()
	projectList := list.New([]list.Item{}, delegate, 0, 0)
	projectList.Title = "Projects"
	projectList.SetShowStatusBar(false)
	projectList.SetFilteringEnabled(false)

	return Model{
		apiClient:   client,
		config:      cfg,
		viewMode:    ProjectListView,
		projectList: projectList,
		keys:        keys,
		loading:     true,
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
			if m.viewMode == ProjectListView {
				m.loading = true
				return m, m.loadProjects
			} else if m.viewMode == KanbanBoardView && m.selectedProject != nil {
				m.loading = true
				return m, m.loadTasks(m.selectedProject.ID)
			}
		case key.Matches(msg, m.keys.Back):
			if m.viewMode == KanbanBoardView {
				m.viewMode = ProjectListView
				m.selectedProject = nil
				m.kanbanBoard = nil
				return m, nil
			}
		case key.Matches(msg, m.keys.Enter):
			if m.viewMode == ProjectListView && !m.loading {
				if len(m.projects) > 0 {
					selected := m.projectList.SelectedItem()
					if proj, ok := selected.(projectItem); ok {
						m.selectedProject = &proj.Project
						m.viewMode = LoadingView
						m.loading = true
						return m, m.loadTasks(proj.Project.ID)
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.projectList.SetSize(msg.Width, msg.Height-4)
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
		
		// Update project list
		items := make([]list.Item, len(m.projects))
		for i, p := range m.projects {
			items[i] = projectItem{Project: p}
		}
		m.projectList.SetItems(items)
		return m, nil

	case tasksLoadedMsg:
		m.loading = false
		m.tasks = msg.tasks
		m.err = msg.err
		if m.err != nil {
			m.viewMode = ErrorView
			return m, nil
		}
		
		// Create kanban board
		m.kanbanBoard = NewKanbanBoard(m.selectedProject, m.tasks)
		m.kanbanBoard.SetSize(m.width, m.height)
		m.viewMode = KanbanBoardView
		return m, nil
	}

	// Update the appropriate view
	var cmd tea.Cmd
	switch m.viewMode {
	case ProjectListView:
		m.projectList, cmd = m.projectList.Update(msg)
		return m, cmd
	case KanbanBoardView:
		if m.kanbanBoard != nil {
			*m.kanbanBoard, cmd = m.kanbanBoard.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.loading {
		return m.loadingView()
	}

	switch m.viewMode {
	case ProjectListView:
		return m.projectListView()
	case KanbanBoardView:
		return m.kanbanBoardView()
	case ErrorView:
		return m.errorView()
	default:
		return m.loadingView()
	}
}

func (m Model) projectListView() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	help := helpStyle("\n  ↑/k up • ↓/j down • enter select • r refresh • q quit")
	
	if len(m.projects) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(2).
			MarginLeft(2)
		return emptyStyle.Render("No projects found. Press 'r' to refresh.") + help
	}

	return m.projectList.View() + help
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
	return style.Render("Loading...")
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

type tasksLoadedMsg struct {
	tasks []api.Task
	err   error
}

// Commands

func (m Model) loadProjects() tea.Msg {
	projects, err := m.apiClient.GetProjects()
	return projectsLoadedMsg{projects: projects, err: err}
}

func (m Model) loadTasks(projectID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.apiClient.GetTasks(projectID)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

// projectItem implements list.Item interface
type projectItem struct {
	Project api.Project
}

func (p projectItem) Title() string       { return p.Project.Name }
func (p projectItem) Description() string { return p.Project.Description }
func (p projectItem) FilterValue() string { return p.Project.Name }
