package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages
func (b KanbanBoard) Update(msg tea.Msg) (KanbanBoard, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			// Progress: move project to next status
			if project := b.GetSelectedProject(); project != nil {
				nextStatus := b.GetNextStatus()
				if nextStatus != "" {
					return b, func() tea.Msg {
						// This will be handled by the parent Model
						return progressProjectMsg{projectID: project.ID, status: nextStatus}
					}
				}
			}
		case "r":
			// Regress: move project to previous status
			if project := b.GetSelectedProject(); project != nil {
				prevStatus := b.GetPrevStatus()
				if prevStatus != "" {
					return b, func() tea.Msg {
						// This will be handled by the parent Model
						return regressProjectMsg{projectID: project.ID, status: prevStatus}
					}
				}
			}
		case "+", "=":
			// Increase priority (maximum 3)
			if project := b.GetSelectedProject(); project != nil {
				newPriority := project.Priority + 1
				if newPriority > 3 {
					newPriority = 3
				}
				return b, func() tea.Msg {
					// This will be handled by the parent Model
					return updatePriorityMsg{projectID: project.ID, priority: newPriority}
				}
			}
		case "-", "_":
			// Decrease priority (minimum 0)
			if project := b.GetSelectedProject(); project != nil {
				newPriority := project.Priority - 1
				if newPriority < 0 {
					newPriority = 0
				}
				return b, func() tea.Msg {
					// This will be handled by the parent Model
					return updatePriorityMsg{projectID: project.ID, priority: newPriority}
				}
			}
		case "left", "h":
			// Move left, skipping empty columns
			for i := b.selectedCol - 1; i >= 0; i-- {
				if len(b.columns[i].Projects) > 0 {
					b.selectedCol = i
					// Try to maintain position using desiredProject, but clamp to available range
					b.selectedProject = min(b.desiredProject, len(b.columns[i].Projects)-1)
					// Restore the desired scroll offset for this column
					b.scrollOffset[i] = b.desiredScrollOffset[i]

					// Ensure selected project is visible - scroll up if needed
					b.scrollOffset[i] = min(b.selectedProject, b.scrollOffset[i])

					// Also check if we need to scroll down to keep it visible
					maxHeight := b.height - 10
					if maxHeight < 5 {
						maxHeight = 5
					}
					maxProjects := maxHeight / 7
					if maxProjects < 1 {
						maxProjects = 1
					}
					visibleEnd := b.scrollOffset[i] + maxProjects - 1
					if b.selectedProject > visibleEnd {
						b.scrollOffset[i] = b.selectedProject - maxProjects + 1
						if b.scrollOffset[i] < 0 {
							b.scrollOffset[i] = 0
						}
					}
					break
				}
			}
		case "right", "l":
			// Move right, skipping empty columns
			for i := b.selectedCol + 1; i < len(b.columns); i++ {
				if len(b.columns[i].Projects) > 0 {
					b.selectedCol = i
					// Try to maintain position using desiredProject, but clamp to available range
					b.selectedProject = min(b.desiredProject, len(b.columns[i].Projects)-1)
					// Restore the desired scroll offset for this column
					b.scrollOffset[i] = b.desiredScrollOffset[i]

					// Ensure selected project is visible - scroll up if needed
					b.scrollOffset[i] = min(b.selectedProject, b.scrollOffset[i])

					// Also check if we need to scroll down to keep it visible
					maxHeight := b.height - 10
					if maxHeight < 5 {
						maxHeight = 5
					}
					maxProjects := maxHeight / 7
					if maxProjects < 1 {
						maxProjects = 1
					}
					visibleEnd := b.scrollOffset[i] + maxProjects - 1
					if b.selectedProject > visibleEnd {
						b.scrollOffset[i] = b.selectedProject - maxProjects + 1
						b.scrollOffset[i] = min(b.scrollOffset[i], 0)
					}
					break
				}
			}
		case "up", "k":
			if b.selectedProject > 0 {
				b.selectedProject--
				// Update desiredProject to track the maximum index reached
				b.desiredProject = b.selectedProject
				// Scroll up if we've scrolled above the visible area
				b.scrollOffset[b.selectedCol] = min(b.selectedProject, b.scrollOffset[b.selectedCol])
				// Save the current scroll position as desired
				b.desiredScrollOffset[b.selectedCol] = b.scrollOffset[b.selectedCol]
			}
		case "down", "j":
			currentCol := b.columns[b.selectedCol]
			if b.selectedProject < len(currentCol.Projects)-1 {
				b.selectedProject++
				// Update desiredProject to track the maximum index reached
				b.desiredProject = b.selectedProject

				// Calculate how many projects can fit in the visible area
				maxHeight := b.height - 10 // Reserve space for title and help
				if maxHeight < 5 {
					maxHeight = 5
				}
				maxProjects := maxHeight / 7 // Each project card is ~7 lines
				if maxProjects < 1 {
					maxProjects = 1
				}

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
				// Save the current scroll position as desired
				b.desiredScrollOffset[b.selectedCol] = b.scrollOffset[b.selectedCol]
			}
		}
	}
	return b, nil
}
