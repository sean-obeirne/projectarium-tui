package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
)

// ProjectModal represents the modal for creating a new project
type ProjectModal struct {
	inputs         []textinput.Model
	focusedIndex   int
	width          int
	height         int
	statusOptions  []string
	selectedStatus int // Index of selected status
	err            string
	isEditMode     bool
	projectID      int
}

const (
	nameField = iota
	descriptionField
	pathField
	fileField
	languageField
	priorityField
	statusField
	totalFields
)

// NewProjectModal creates a new project creation modal
func NewProjectModal() *ProjectModal {
	inputs := make([]textinput.Model, totalFields)

	// Name input
	inputs[nameField] = textinput.New()
	inputs[nameField].Placeholder = "(required)"
	inputs[nameField].CharLimit = 100
	inputs[nameField].Width = 40
	inputs[nameField].Focus()

	// Description input
	inputs[descriptionField] = textinput.New()
	inputs[descriptionField].Placeholder = ""
	inputs[descriptionField].CharLimit = 200
	inputs[descriptionField].Width = 40

	// Path input
	inputs[pathField] = textinput.New()
	inputs[pathField].Placeholder = ""
	inputs[pathField].CharLimit = 256
	inputs[pathField].Width = 40

	// File input
	inputs[fileField] = textinput.New()
	inputs[fileField].Placeholder = ""
	inputs[fileField].CharLimit = 256
	inputs[fileField].Width = 40

	// Language input
	inputs[languageField] = textinput.New()
	inputs[languageField].Placeholder = ""
	inputs[languageField].CharLimit = 50
	inputs[languageField].Width = 40

	// Priority input
	inputs[priorityField] = textinput.New()
	inputs[priorityField].Placeholder = ""
	inputs[priorityField].CharLimit = 1
	inputs[priorityField].Width = 40
	inputs[priorityField].SetValue("0")

	// Status field is not a text input - it's a selector
	inputs[statusField] = textinput.New()
	inputs[statusField].Placeholder = ""
	inputs[statusField].CharLimit = 0

	return &ProjectModal{
		inputs:         inputs,
		focusedIndex:   0,
		statusOptions:  []string{"ready", "in_progress", "finished"},
		selectedStatus: 0,
		isEditMode:     false,
		projectID:      0,
	}
}

// NewProjectModalForEdit creates a modal pre-populated with existing project data
func NewProjectModalForEdit(project *api.Project) *ProjectModal {
	modal := NewProjectModal()
	modal.isEditMode = true
	modal.projectID = project.ID
	
	// Populate fields with existing values
	modal.inputs[nameField].SetValue(project.Name)
	modal.inputs[descriptionField].SetValue(project.Description)
	modal.inputs[pathField].SetValue(project.Path)
	modal.inputs[fileField].SetValue(project.File)
	modal.inputs[languageField].SetValue(project.Language)
	modal.inputs[priorityField].SetValue(fmt.Sprintf("%d", project.Priority))
	
	// Set selected status based on project status
	for i, status := range modal.statusOptions {
		if status == project.Status {
			modal.selectedStatus = i
			break
		}
	}
	
	return modal
}

// SetSize sets the modal dimensions
func (m *ProjectModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages for the project modal
func (m ProjectModal) Update(msg tea.Msg) (ProjectModal, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel and close modal
			return m, cancelProjectCreationCmd()

		case "enter":
			// Submit if we're on the last field or status field
			if m.focusedIndex >= statusField {
				return m, m.submitProject()
			}
			// Otherwise, move to next field
			m.focusNext()
			return m, nil

		case "tab", "down":
			m.focusNext()
			return m, nil

		case "shift+tab", "up":
			m.focusPrev()
			return m, nil

		case "left", "j":
			// Navigate status options
			if m.focusedIndex == statusField {
				if m.selectedStatus > 0 {
					m.selectedStatus--
				}
			}

		case "right", ";":
			// Navigate status options
			if m.focusedIndex == statusField {
				if m.selectedStatus < len(m.statusOptions)-1 {
					m.selectedStatus++
				}
			}
		}
	}

	// Update the focused text input
	if m.focusedIndex < statusField {
		m.inputs[m.focusedIndex], cmd = m.inputs[m.focusedIndex].Update(msg)
	}

	return m, cmd
}

func (m *ProjectModal) focusNext() {
	m.inputs[m.focusedIndex].Blur()
	m.focusedIndex++
	if m.focusedIndex >= totalFields {
		m.focusedIndex = 0
	}
	if m.focusedIndex < statusField {
		m.inputs[m.focusedIndex].Focus()
	}
}

func (m *ProjectModal) focusPrev() {
	m.inputs[m.focusedIndex].Blur()
	m.focusedIndex--
	if m.focusedIndex < 0 {
		m.focusedIndex = totalFields - 1
	}
	if m.focusedIndex < statusField {
		m.inputs[m.focusedIndex].Focus()
	}
}

func (m *ProjectModal) submitProject() tea.Cmd {
	// Validate required fields
	name := m.inputs[nameField].Value()
	if name == "" {
		m.err = "Project name is required"
		return nil
	}

	// Parse priority
	priorityStr := m.inputs[priorityField].Value()
	if priorityStr == "" {
		priorityStr = "0"
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil || priority < 0 || priority > 3 {
		m.err = "Priority must be a number between 0 and 3"
		return nil
	}

	if m.isEditMode {
		return updateProjectCmd(
			m.projectID,
			name,
			m.inputs[descriptionField].Value(),
			m.inputs[pathField].Value(),
			m.inputs[fileField].Value(),
			m.inputs[languageField].Value(),
			priority,
			m.statusOptions[m.selectedStatus],
		)
	}

	return createProjectCmd(
		name,
		m.inputs[descriptionField].Value(),
		m.inputs[pathField].Value(),
		m.inputs[fileField].Value(),
		m.inputs[languageField].Value(),
		priority,
		m.statusOptions[m.selectedStatus],
	)
}

// View renders the project modal
func (m *ProjectModal) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(20)

	focusedLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Bold(true).
		Width(20)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	statusStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	selectedStatusStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Bold(true)

	// Title
	var title string
	if m.isEditMode {
		title = titleStyle.Render("✏️  Edit Project")
	} else {
		title = titleStyle.Render("✨ Create New Project")
	}

	// Build form fields
	var formFields []string

	fieldLabels := []string{
		"Project Name*:",
		"Description:",
		"Project Dir:",
		"Default File:",
		"Language:",
		"Priority:",
		"Status:",
	}

	for i := 0; i < totalFields-1; i++ {
		var label string
		if i == m.focusedIndex {
			label = focusedLabelStyle.Render(fieldLabels[i])
		} else {
			label = labelStyle.Render(fieldLabels[i])
		}

		field := lipgloss.JoinHorizontal(
			lipgloss.Top,
			label,
			m.inputs[i].View(),
		)
		formFields = append(formFields, field)
	}

	// Status selector
	var statusLabel string
	if m.focusedIndex == statusField {
		statusLabel = focusedLabelStyle.Render(fieldLabels[statusField])
	} else {
		statusLabel = labelStyle.Render(fieldLabels[statusField])
	}

	var statusButtons []string
	for i, status := range m.statusOptions {
		var displayText string
		switch status {
		case "ready":
			displayText = "Ready"
		case "in_progress":
			displayText = "In Progress"
		case "finished":
			displayText = "Finished"
		}

		var button string
		if i == m.selectedStatus {
			button = selectedStatusStyle.Render(displayText)
		} else {
			button = statusStyle.Render(displayText)
		}
		statusButtons = append(statusButtons, button)
	}

	statusSelector := lipgloss.JoinHorizontal(lipgloss.Top, statusButtons...)
	statusField := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusLabel,
		statusSelector,
	)
	formFields = append(formFields, statusField)

	form := lipgloss.JoinVertical(lipgloss.Left, formFields...)

	// Help text
	help := helpStyle.Render("tab/↓ next • shift+tab/↑ prev • ←/j →/; navigate status • enter submit • esc cancel")

	// Error message
	var errorMsg string
	if m.err != "" {
		errorMsg = errorStyle.Render(fmt.Sprintf("Error: %s", m.err))
	}

	// Combine everything
	sections := []string{title, "", form}
	if errorMsg != "" {
		sections = append(sections, "", errorMsg)
	}
	sections = append(sections, "", help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// Message types for project modal
type createProjectMsg struct {
	name        string
	description string
	path        string
	file        string
	language    string
	priority    int
	status      string
}

type updateProjectMsg struct {
	id          int
	name        string
	description string
	path        string
	file        string
	language    string
	priority    int
	status      string
}

type cancelProjectCreationMsg struct{}

// Command functions
func createProjectCmd(name, description, path, file, language string, priority int, status string) tea.Cmd {
	return func() tea.Msg {
		return createProjectMsg{
			name:        name,
			description: description,
			path:        path,
			file:        file,
			language:    language,
			priority:    priority,
			status:      status,
		}
	}
}

func updateProjectCmd(id int, name, description, path, file, language string, priority int, status string) tea.Cmd {
	return func() tea.Msg {
		return updateProjectMsg{
			id:          id,
			name:        name,
			description: description,
			path:        path,
			file:        file,
			language:    language,
			priority:    priority,
			status:      status,
		}
	}
}

func cancelProjectCreationCmd() tea.Cmd {
	return func() tea.Msg {
		return cancelProjectCreationMsg{}
	}
}
