package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/renderorange/chroma/chroma-control/config"
)

func (m *Model) renderPresetBrowser() string {
	switch m.presetBrowser.mode {
	case browserModeList:
		return m.renderPresetList()
	case browserModeConfirmDelete:
		return m.renderDeleteConfirmation()
	case browserModeSaveAs:
		return m.renderSaveAsDialog()
	}
	return ""
}

func (m *Model) renderPresetList() string {
	modalWidth := 60
	if m.width > 0 && m.width < modalWidth+4 {
		modalWidth = m.width - 4
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Width(modalWidth - 4)
	itemStyle := lipgloss.NewStyle().
		Foreground(colorTextNormal).
		Width(modalWidth - 4)
	selectedStyle := lipgloss.NewStyle().
		Foreground(colorTextHighlight).
		Bold(true).
		Width(modalWidth - 4)
	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)
	errorStyle := lipgloss.NewStyle().
		Foreground(colorTextError)
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2)

	var content strings.Builder

	content.WriteString(titleStyle.Render("Load Preset"))
	content.WriteString("\n\n")

	if m.presetBrowser.errorMsg != "" {
		content.WriteString(errorStyle.Render(m.presetBrowser.errorMsg))
		content.WriteString("\n\n")
	}

	if len(m.presetBrowser.presets) == 0 {
		content.WriteString(mutedStyle.Render("No presets saved yet"))
	} else {
		for i, preset := range m.presetBrowser.presets {
			style := itemStyle
			prefix := "  "
			if i == m.presetBrowser.selectedIdx {
				style = selectedStyle
				prefix = "> "
			}
			content.WriteString(style.Render(prefix + preset))
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(mutedStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("enter:load  d:delete  n:new  esc:back"))

	return m.centerModal(modalStyle.Width(modalWidth).Render(content.String()))
}

func (m *Model) renderDeleteConfirmation() string {
	modalWidth := 50

	titleStyle := lipgloss.NewStyle().
		Foreground(colorTextError).
		Bold(true)
	textStyle := lipgloss.NewStyle().
		Foreground(colorTextNormal)
	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorTextError).
		Padding(1, 2)

	var content strings.Builder
	content.WriteString(titleStyle.Render("Delete Preset?"))
	content.WriteString("\n\n")
	content.WriteString(textStyle.Render(fmt.Sprintf("Delete '%s'?", m.presetBrowser.confirmName)))
	content.WriteString("\n\n")
	content.WriteString(mutedStyle.Render("y:confirm  n:cancel"))

	return m.centerModal(modalStyle.Width(modalWidth).Render(content.String()))
}

func (m *Model) renderSaveAsDialog() string {
	modalWidth := 50

	titleStyle := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true)
	inputStyle := lipgloss.NewStyle().
		Foreground(colorTextHighlight)
	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2)

	var content strings.Builder
	content.WriteString(titleStyle.Render("Save Preset As"))
	content.WriteString("\n\n")
	content.WriteString(inputStyle.Render(m.presetBrowser.inputBuffer + "_"))
	content.WriteString("\n\n")
	content.WriteString(mutedStyle.Render("enter:save  esc:cancel"))

	return m.centerModal(modalStyle.Width(modalWidth).Render(content.String()))
}

func (m *Model) centerModal(content string) string {
	// Simple centering based on terminal dimensions
	if m.width == 0 || m.height == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	contentWidth := 0
	for _, line := range lines {
		if w := lipgloss.Width(line); w > contentWidth {
			contentWidth = w
		}
	}

	horizontalPad := (m.width - contentWidth) / 2
	if horizontalPad < 0 {
		horizontalPad = 0
	}

	verticalPad := (m.height - len(lines)) / 3
	if verticalPad < 0 {
		verticalPad = 0
	}

	var result strings.Builder
	for i := 0; i < verticalPad; i++ {
		result.WriteString("\n")
	}

	pad := strings.Repeat(" ", horizontalPad)
	for _, line := range lines {
		result.WriteString(pad)
		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}

func (m *Model) updatePresetBrowser(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.presetBrowser.mode {
		case browserModeList:
			return m.handlePresetListKeys(msg)
		case browserModeConfirmDelete:
			return m.handleDeleteConfirmKeys(msg)
		case browserModeSaveAs:
			return m.handleSaveAsKeys(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *Model) handlePresetListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.presetBrowser.errorMsg = "" // Clear error on any key

	switch msg.String() {
	case "esc", "q":
		// Go back to splash or settings depending on where we came from
		if m.prevScreen == screenSplash {
			m.switchScreen(screenSplash)
		} else {
			m.switchScreen(screenSettings)
		}
		return m, nil

	case "j", "down":
		if m.presetBrowser.selectedIdx < len(m.presetBrowser.presets)-1 {
			m.presetBrowser.selectedIdx++
		}
		return m, nil

	case "k", "up":
		if m.presetBrowser.selectedIdx > 0 {
			m.presetBrowser.selectedIdx--
		}
		return m, nil

	case "enter":
		if len(m.presetBrowser.presets) > 0 && m.presetBrowser.selectedIdx < len(m.presetBrowser.presets) {
			name := m.presetBrowser.presets[m.presetBrowser.selectedIdx]
			if preset, err := config.LoadPreset(name); err == nil {
				// Initialize lists first if needed
				if m.width > 0 && m.effectsList.Items() == nil {
					m.InitLists(m.width, m.height)
				}
				m.applyPreset(preset)
				m.currentPresetName = name
				config.SaveLastPresetName(name)
				m.switchScreen(screenMain)
			} else {
				m.presetBrowser.errorMsg = "Failed to load preset"
			}
		}
		return m, nil

	case "d":
		if len(m.presetBrowser.presets) > 0 {
			m.presetBrowser.confirmName = m.presetBrowser.presets[m.presetBrowser.selectedIdx]
			m.presetBrowser.mode = browserModeConfirmDelete
		}
		return m, nil

	case "n":
		// New preset - enter save-as mode
		m.presetBrowser.inputBuffer = ""
		m.presetBrowser.mode = browserModeSaveAs
		return m, nil
	}
	return m, nil
}

func (m *Model) handleDeleteConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		config.DeletePreset(m.presetBrowser.confirmName)
		m.refreshPresetList()
		m.presetBrowser.mode = browserModeList

	case "n", "esc":
		m.presetBrowser.mode = browserModeList
	}
	return m, nil
}

func (m *Model) handleSaveAsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.presetBrowser.mode = browserModeList
		return m, nil

	case tea.KeyEnter:
		if m.presetBrowser.inputBuffer != "" {
			preset := m.buildCurrentPreset()
			if err := config.SavePreset(preset, m.presetBrowser.inputBuffer); err != nil {
				m.presetBrowser.errorMsg = "Failed to save preset"
				m.presetBrowser.mode = browserModeList
			} else {
				m.currentPresetName = m.presetBrowser.inputBuffer
				config.SaveLastPresetName(m.presetBrowser.inputBuffer)
				m.isDirty = false
				// Return to main screen after saving
				m.switchScreen(screenMain)
			}
		}
		return m, nil

	case tea.KeyBackspace:
		if len(m.presetBrowser.inputBuffer) > 0 {
			m.presetBrowser.inputBuffer = m.presetBrowser.inputBuffer[:len(m.presetBrowser.inputBuffer)-1]
		}
		return m, nil

	case tea.KeyRunes:
		m.presetBrowser.inputBuffer += string(msg.Runes)
		return m, nil
	}
	return m, nil
}

func (m *Model) refreshPresetList() {
	m.presetBrowser.presets, _ = config.ListPresets()
	if m.presetBrowser.selectedIdx >= len(m.presetBrowser.presets) {
		m.presetBrowser.selectedIdx = len(m.presetBrowser.presets) - 1
	}
	if m.presetBrowser.selectedIdx < 0 {
		m.presetBrowser.selectedIdx = 0
	}
}
