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

	// Handle quit confirmation
	if m.showQuitConfirm {
		return m.renderQuitConfirmation()
	}

	// Route to appropriate screen view
	switch m.screen {
	case screenSplash:
		return m.renderSplash()
	case screenPresetBrowser:
		return m.renderPresetBrowser()
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
		return m.renderCommandPalette()
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
			parts = append([]string{"j/k:nav", "h/l:reorder", "enter:grab", "esc:back"}, parts...)
		} else {
			parts = append([]string{"j/k:nav", "h/l:adjust", "enter:toggle", "esc:back"}, parts...)
		}
	}

	// Always shown shortcuts at the end
	parts = append(parts, ":commands")
	parts = append(parts, "?:help")

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

	// Preset name and dirty indicator
	presetDisplay := m.currentPresetName
	if presetDisplay == "" {
		presetDisplay = "(unsaved)"
	}
	if m.isDirty {
		presetDisplay += " *"
	}
	presetSection := lipgloss.NewStyle().Foreground(colorTextMuted).Render(presetDisplay)

	// Build status bar - connection, preset, and MIDI
	leftSection := connectionStatus
	middleSection := presetSection
	rightSection := midiStatus

	// Calculate widths
	leftWidth := lipgloss.Width(leftSection)
	middleWidth := lipgloss.Width(middleSection)
	rightWidth := lipgloss.Width(rightSection)

	// Calculate spacing to center the middle section
	totalContentWidth := leftWidth + middleWidth + rightWidth
	availableSpace := width - totalContentWidth - 4 // padding

	if availableSpace < 2 {
		availableSpace = 2
	}

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
		middleSection +
		lipgloss.NewStyle().Width(rightPadding).Render("") +
		rightSection

	return statusBarStyle.Render(statusText)
}

func (m *Model) renderQuitConfirmation() string {
	modalWidth := 50
	if m.width > 0 && m.width < modalWidth+4 {
		modalWidth = m.width - 4
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(colorTextError).
		Bold(true).
		Width(modalWidth - 4)
	textStyle := lipgloss.NewStyle().
		Foreground(colorTextNormal).
		Width(modalWidth - 4)
	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorTextError).
		Padding(1, 2)

	var content strings.Builder
	content.WriteString(titleStyle.Render("Unsaved Changes"))
	content.WriteString("\n\n")
	content.WriteString(textStyle.Render("Save changes before quitting?"))
	content.WriteString("\n\n")
	content.WriteString(mutedStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("s:save  d:discard  c:cancel"))

	return m.centerModal(modalStyle.Width(modalWidth).Render(content.String()))
}

// overlay overlays foreground content on top of background content
func overlay(background, foreground string) string {
	bgLines := strings.Split(background, "\n")
	fgLines := strings.Split(foreground, "\n")

	// Find where to place the foreground (center it)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)
	startY := (bgHeight - fgHeight) / 2
	if startY < 0 {
		startY = 0
	}

	var result strings.Builder
	for i := 0; i < startY && i < len(bgLines); i++ {
		result.WriteString(bgLines[i])
		result.WriteString("\n")
	}

	// Overlay foreground lines
	for i, fgLine := range fgLines {
		if startY+i < len(bgLines) {
			bgLine := bgLines[startY+i]
			fgWidth := lipgloss.Width(fgLine)
			bgWidth := len(bgLine)

			// Center the foreground line over the background
			startX := (bgWidth - fgWidth) / 2
			if startX < 0 {
				startX = 0
			}

			// Write background up to start position
			if startX < len(bgLine) {
				result.WriteString(bgLine[:startX])
			}

			// Write foreground
			result.WriteString(fgLine)

			// Write rest of background
			endX := startX + fgWidth
			if endX < len(bgLine) {
				result.WriteString(bgLine[endX:])
			}
		} else {
			result.WriteString(fgLine)
		}
		result.WriteString("\n")
	}

	// Write remaining background lines
	for i := startY + fgHeight; i < len(bgLines); i++ {
		result.WriteString(bgLines[i])
		result.WriteString("\n")
	}

	return strings.TrimRight(result.String(), "\n")
}
