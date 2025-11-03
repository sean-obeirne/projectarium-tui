package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sean-obeirne/projectarium-tui/internal/api"
)

// KanbanBoard represents the kanban board view
type KanbanBoard struct {
	columns         []ProjectColumn
	selectedCol     int
	selectedProject int
	scrollOffset    []int // scroll offset for each column
	width           int
	height          int
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
		columns:         columns,
		selectedCol:     0,
		selectedProject: 0,
		scrollOffset:    make([]int, len(columns)),
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

// Update handles messages
func (b KanbanBoard) Update(msg tea.Msg) (KanbanBoard, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			// Move left, skipping empty columns
			for i := b.selectedCol - 1; i >= 0; i-- {
				if len(b.columns[i].Projects) > 0 {
					b.selectedCol = i
					b.selectedProject = 0
					b.scrollOffset[i] = 0 // Reset scroll when switching columns
					break
				}
			}
		case "right", "l":
			// Move right, skipping empty columns
			for i := b.selectedCol + 1; i < len(b.columns); i++ {
				if len(b.columns[i].Projects) > 0 {
					b.selectedCol = i
					b.selectedProject = 0
					b.scrollOffset[i] = 0 // Reset scroll when switching columns
					break
				}
			}
		case "up", "k":
			if b.selectedProject > 0 {
				b.selectedProject--
				// Scroll up if we've scrolled above the visible area
				if b.selectedProject < b.scrollOffset[b.selectedCol] {
					b.scrollOffset[b.selectedCol] = b.selectedProject
				}
			}
		case "down", "j":
			currentCol := b.columns[b.selectedCol]
			if b.selectedProject < len(currentCol.Projects)-1 {
				b.selectedProject++

				// Calculate how many projects can fit in the visible area
				maxHeight := b.height - 8          // Reserve space for title and help
				maxProjects := (maxHeight - 2) / 5 // Approximate projects that fit

				// Check if we need to scroll down
				visibleStart := b.scrollOffset[b.selectedCol]
				visibleEnd := visibleStart + maxProjects - 1

				if b.selectedProject > visibleEnd {
					// Scroll down so the selected project is at the bottom of visible area
					b.scrollOffset[b.selectedCol] = b.selectedProject - maxProjects + 1
					if b.scrollOffset[b.selectedCol] < 0 {
						b.scrollOffset[b.selectedCol] = 0
					}
				}
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
		MarginBottom(1).
		MarginLeft(2)

	// Column colors: blue, yellow, green
	columnColors := []lipgloss.Color{
		lipgloss.Color("63"),  // Blue for "Ready"
		lipgloss.Color("214"), // Yellow for "In Progress"
		lipgloss.Color("70"),  // Green for "Finished"
	}

	columnHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.Color("15")).
		Align(lipgloss.Center)

	// Priority-based border colors
	lowPriorityBorder := lipgloss.Color("240")    // Gray
	mediumPriorityBorder := lipgloss.Color("214") // Orange/Yellow
	highPriorityBorder := lipgloss.Color("196")   // Red
	selectedBorder := lipgloss.Color("51")        // Purple

	projectCardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder())

	selectedProjectCardStyle := projectCardStyle.
		BorderForeground(selectedBorder).
		Bold(true)

	emptyColumnStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(2).
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1).
		MarginLeft(2)

	// Title
	title := titleStyle.Render("📋 Project Board")

	// Calculate column width
	colWidth := (b.width - 8) / len(b.columns)
	if colWidth < 25 {
		colWidth = 25
	}

	// Build columns
	columnViews := make([]string, len(b.columns))
	maxHeight := b.height - 8 // Reserve space for title and help

	for i, col := range b.columns {
		// Column border style
		columnBorderStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(columnColors[i]).
			Padding(0, 1)

		// Column header
		headerStyle := columnHeaderStyle.Background(columnColors[i])
		header := headerStyle.Width(colWidth - 4).Render(fmt.Sprintf("%s (%d)", col.Name, len(col.Projects)))

		// Projects
		projectViews := []string{}
		displayedProjects := 0
		maxProjects := (maxHeight - 2) / 5 // Approximate projects that fit

		// Calculate scroll range for this column
		scrollStart := b.scrollOffset[i]
		scrollEnd := scrollStart + maxProjects
		if scrollEnd > len(col.Projects) {
			scrollEnd = len(col.Projects)
		}

		for j := scrollStart; j < scrollEnd; j++ {
			project := col.Projects[j]

			// Determine border color based on priority
			var borderColor lipgloss.Color
			switch project.Priority {
			case 0:
				borderColor = highPriorityBorder
			case 1:
				borderColor = mediumPriorityBorder
			case 2:
				borderColor = lowPriorityBorder
			default:
				borderColor = lowPriorityBorder
			}

			style := projectCardStyle.BorderForeground(borderColor)
			if i == b.selectedCol && j == b.selectedProject {
				style = selectedProjectCardStyle
			}

			// Build project card content
			name := project.Name
			maxNameLen := colWidth - 6
			if len(name) > maxNameLen {
				name = name[:maxNameLen-3] + "..."
			}

			description := project.Description
			maxDescLen := colWidth - 6
			if len(description) > maxDescLen {
				description = description[:maxDescLen-3] + "..."
			}

			// Language badge (right-aligned)
			languageBadge := ""
			if project.Language != "" {
				languageBadge = project.Language
			}

			// Create header with name (left) and language (right)
			headerWidth := colWidth - 6
			nameStyle := lipgloss.NewStyle().Align(lipgloss.Left)
			langStyle := lipgloss.NewStyle().Align(lipgloss.Right).Foreground(lipgloss.Color("241"))

			header := lipgloss.JoinHorizontal(
				lipgloss.Top,
				nameStyle.Width(headerWidth-len(languageBadge)).Render(name),
				langStyle.Width(len(languageBadge)).Render(languageBadge),
			)

			// Description (centered)
			descStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(headerWidth)
			descContent := descStyle.Render(description)

			cardContent := lipgloss.JoinVertical(lipgloss.Left, header, descContent)
			projectView := style.Width(colWidth - 2).Render(cardContent)
			projectViews = append(projectViews, projectView)
			displayedProjects++
		}

		// Add scrolling indicators
		if scrollStart > 0 {
			indicator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Align(lipgloss.Center).
				Render(fmt.Sprintf("↑ %d more above", scrollStart))
			projectViews = append([]string{indicator}, projectViews...)
		}

		if scrollEnd < len(col.Projects) {
			indicator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Align(lipgloss.Center).
				Render(fmt.Sprintf("↓ %d more below", len(col.Projects)-scrollEnd))
			projectViews = append(projectViews, indicator)
		}

		if len(projectViews) == 0 {
			projectViews = append(projectViews, emptyColumnStyle.Width(colWidth-6).Render("No projects"))
		}

		columnContent := lipgloss.JoinVertical(lipgloss.Left, projectViews...)
		columnWithHeader := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			columnContent,
		)

		// Wrap entire column in border
		columnViews[i] = columnBorderStyle.Render(columnWithHeader)
	}

	// Join columns horizontally with spacing
	board := lipgloss.JoinHorizontal(lipgloss.Top, columnViews...)

	// Help text
	help := helpStyle.Render("  ←/h →/l columns • ↑/k ↓/j projects • r refresh • q quit")

	// Combine everything
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		board,
		help,
	)
}
