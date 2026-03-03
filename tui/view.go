package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(2, 3)
)

func (m Model) View() string {
	// Handle tiny terminals
	if m.width > 0 && (m.width < 60 || m.height < 20) {
		return appStyle.Render("Terminal too small. Minimum: 60x20")
	}

	// Route to appropriate screen view
	switch m.screen {
	case screenSplash:
		return m.renderSplash()
	case screenSettings:
		return m.renderSettings()
	case screenHelp:
		return m.renderHelp()
	case screenMain:
		return m.renderMain()
	default:
		// Loading state or uninitialized
		if m.effectsList.Items() == nil || len(m.effectsList.Items()) == 0 {
			if m.width > 0 {
				return "Loading..."
			}
			return ""
		}
		return m.renderMain()
	}
}

func (m Model) renderMain() string {
	if m.effectsList.Items() == nil || len(m.effectsList.Items()) == 0 {
		if m.width > 0 {
			return appStyle.Render("Loading...")
		}
		return ""
	}

	mainContent := m.renderMainBase()

	if m.showCommandPalette {
		return appStyle.Render(m.renderCommandPalette())
	}

	return appStyle.Render(mainContent)
}

// renderMainBase renders the main content without checking for modal overlays
func (m Model) renderMainBase() string {
	if m.effectsList.Items() == nil || len(m.effectsList.Items()) == 0 {
		return ""
	}

	// Save and hide selection in unfocused pane
	var effectsSelection, paramsSelection int
	if m.navigationMode == modeEffectsList {
		paramsSelection = m.parameterList.Index()
		m.parameterList.Select(-1) // Hide selection in unfocused pane
	} else {
		effectsSelection = m.effectsList.Index()
		m.effectsList.Select(-1) // Hide selection in unfocused pane
	}

	// Panel styles - use muted text for unfocused pane
	var leftPanelStyle, rightPanelStyle lipgloss.Style
	if m.navigationMode == modeEffectsList {
		leftPanelStyle = lipgloss.NewStyle().
			Padding(2, 2, 1, 2)
		rightPanelStyle = lipgloss.NewStyle().
			Padding(2, 2, 1, 2).
			Foreground(colorTextMuted)
	} else {
		leftPanelStyle = lipgloss.NewStyle().
			Padding(2, 2, 1, 2).
			Foreground(colorTextMuted)
		rightPanelStyle = lipgloss.NewStyle().
			Padding(2, 2, 1, 2)
	}

	// Render panels
	leftPanel := leftPanelStyle.Render(m.effectsList.View())
	rightPanel := rightPanelStyle.Render(m.parameterList.View())

	// Restore selection in panes
	if m.navigationMode == modeEffectsList {
		m.parameterList.Select(paramsSelection)
	} else {
		m.effectsList.Select(effectsSelection)
	}

	// Join panels horizontally
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Calculate available width (accounting for app padding: 2 on each side + divider)
	availableWidth := m.width - 6

	// Render footer and status bar only if status is enabled
	var mainContent string
	if m.showStatus {
		footer := m.renderFooter(availableWidth)
		statusBar := m.renderStatusBar(availableWidth)
		spacer := lipgloss.NewStyle().Height(1).Render("")

		mainContent = lipgloss.JoinVertical(
			lipgloss.Left,
			panels,
			footer,
			spacer,
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

	return mainContent
}

func formatValue(value, min, max float32) string {
	if max-min >= 10 {
		return fmt.Sprintf("%5.1f", value)
	}
	return fmt.Sprintf("%5.2f", value)
}

func (m Model) renderFooter(width int) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted).
		Width(width).
		Padding(0, 1)

	// Build context-aware shortcuts
	var parts []string

	// Always shown shortcuts
	parts = append(parts, ":commands")
	parts = append(parts, "?:help")

	// Mode-specific shortcuts
	switch m.navigationMode {
	case modeEffectsList:
		parts = append([]string{"j/k:nav", "enter:params"}, parts...)
	case modeParameterList:
		// Check if Effects Order parameter is selected in Master section
		idx := m.parameterList.Index()
		items := m.parameterList.Items()
		isEffectsOrderSelected := false
		if idx >= 0 && idx < len(items) {
			if param, ok := items[idx].(parameterItem); ok {
				isEffectsOrderSelected = param.ctrl == ctrlEffectsOrder
			}
		}

		if isEffectsOrderSelected {
			parts = append([]string{"j/k:nav", "pgup/pgdn:reorder", "esc:back"}, parts...)
		} else {
			parts = append([]string{"j/k:nav", "h/l:adjust", "enter:toggle", "esc:back"}, parts...)
		}
	}

	text := strings.Join(parts, " | ")

	return footerStyle.Render(text)
}

func (m Model) renderStatusBar(width int) string {
	connectionStatus := lipgloss.NewStyle().Foreground(colorTextError).Render("Disconnected")
	if m.connected {
		connectionStatus = lipgloss.NewStyle().
			Foreground(colorTextSuccess).
			Render("Connected")
	}

	// MIDI status
	midiStatus := m.midiPort
	if midiStatus == "" {
		midiStatus = "No MIDI"
	}

	// Build status bar - just connection and MIDI
	leftSection := connectionStatus
	rightSection := midiStatus

	// Calculate spacing
	leftWidth := lipgloss.Width(leftSection)
	rightWidth := lipgloss.Width(rightSection)

	availableSpace := width - leftWidth - rightWidth - 4 // padding
	leftPadding := availableSpace / 2
	rightPadding := availableSpace - leftPadding

	if leftPadding < 1 {
		leftPadding = 1
	}
	if rightPadding < 1 {
		rightPadding = 1
	}

	statusBarStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted).
		Width(width).
		Padding(0, 1)

	statusText := leftSection +
		lipgloss.NewStyle().Width(leftPadding).Render("") +
		lipgloss.NewStyle().Width(rightPadding).Render("") +
		rightSection

	return statusBarStyle.Render(statusText)
}
