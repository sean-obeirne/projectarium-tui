package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

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
	// Reserve space for: empty line (1) + title (1) + title margin bottom (1) + empty line (1) + help top margin (1) + help (1) = 6 lines
	// Also account for column header (1) + empty line after header (1) + column borders (2) = 4 more lines
	maxHeight := b.height - 10 // Reserve space for UI chrome
	if maxHeight < 5 {
		maxHeight = 5 // Minimum height to show at least one project
	}

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
		// Each project card is approximately 7 lines: border(2) + padding(2) + content(2) + margin(1)
		maxProjects := maxHeight / 7
		if maxProjects < 1 {
			maxProjects = 1 // Always show at least one project
		}

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
				borderColor = lowPriorityBorder // Gray
			case 1:
				borderColor = mediumPriorityBorder // Yellow/Orange
			default: // 2 or above
				borderColor = highPriorityBorder // Red
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
