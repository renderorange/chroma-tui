package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1)

	connectedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	disconnectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#25A065"))

	unfocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#555555"))
)

func (m Model) View() string {
	// Handle tiny terminals
	if m.width > 0 && (m.width < 60 || m.height < 20) {
		return appStyle.Render("Terminal too small. Minimum: 60x20")
	}

	// Loading state
	if m.effectsList.Items() == nil || len(m.effectsList.Items()) == 0 {
		if m.width > 0 {
			return "Loading..."
		}
		return ""
	}

	// Determine border styles based on focus
	var leftBorderStyle, rightBorderStyle lipgloss.Style
	if m.navigationMode == modeEffectsList {
		leftBorderStyle = focusedBorderStyle
		rightBorderStyle = unfocusedBorderStyle
	} else {
		leftBorderStyle = unfocusedBorderStyle
		rightBorderStyle = focusedBorderStyle
	}

	// Add right margin to left panel for spacing
	leftBorderStyle = leftBorderStyle.MarginRight(2)

	// Render panels with borders
	leftPanel := leftBorderStyle.Render(m.effectsList.View())
	rightPanel := rightBorderStyle.Render(m.parameterList.View())

	// Join panels horizontally
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Calculate available width (accounting for app padding)
	availableWidth := m.width - 4

	// Render footer and status bar only if status is enabled
	var mainContent string
	if m.showStatus {
		footer := m.renderFooter(availableWidth)
		statusBar := m.renderStatusBar(availableWidth)

		mainContent = lipgloss.JoinVertical(
			lipgloss.Left,
			panels,
			footer,
			statusBar,
		)
	} else {
		footer := m.renderFooter(availableWidth)
		mainContent = lipgloss.JoinVertical(
			lipgloss.Left,
			panels,
			footer,
		)
	}

	return appStyle.Render(mainContent)
}

func formatValue(value, min, max float32) string {
	if max-min >= 10 {
		return fmt.Sprintf("%5.1f", value)
	}
	return fmt.Sprintf("%5.2f", value)
}

func (m Model) renderFooter(width int) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(width).
		Padding(0, 1)

	var text string
	switch m.navigationMode {
	case modeEffectsList:
		text = "j/k: navigate | enter: open params | i: grain intensity | 1/2/3: blend mode | q: quit"
	case modeParameterList:
		if m.currentSection == "global" {
			text = "j/k: navigate | pgup/pgdn: reorder | r: reset order | esc: back | q: quit"
		} else {
			text = "j/k: navigate | h/l: adjust value | enter: toggle | esc: back | q: quit"
		}
	}

	return footerStyle.Render(text)
}

func (m Model) renderStatusBar(width int) string {
	connectionStatus := disconnectedStyle.Render("Disconnected")
	if m.connected {
		connectionStatus = connectedStyle.Render("Connected")
	}

	midiStatus := m.midiPort
	if midiStatus == "" {
		midiStatus = "No MIDI"
	}

	return statusBarStyle.Render(
		fmt.Sprintf("%s | MIDI: %s", connectionStatus, midiStatus),
	)
}
