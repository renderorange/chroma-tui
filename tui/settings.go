package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/renderorange/chroma/chroma-control/config"
)

func (m *Model) renderSettings() string {
	modalWidth := 50
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
		Foreground(colorPrimary).
		Bold(true).
		Width(modalWidth - 4)
	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)
	infoStyle := lipgloss.NewStyle().
		Foreground(colorSecondary)
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2)

	var content strings.Builder

	content.WriteString(titleStyle.Render("Settings"))
	content.WriteString("\n\n")

	// Show current preset name
	currentName := m.currentPresetName
	if currentName == "" {
		currentName = "(unsaved)"
	}
	content.WriteString(infoStyle.Render("Current: " + currentName))
	if m.isDirty {
		content.WriteString(infoStyle.Render(" *"))
	}
	content.WriteString("\n\n")

	items := []struct {
		id   settingsMenuItem
		name string
	}{
		{settingsSave, "Save"},
		{settingsSaveAs, "Save As"},
		{settingsLoad, "Load"},
		{settingsReset, "Reset to Defaults"},
	}

	for _, item := range items {
		style := itemStyle
		prefix := "  "
		if item.id == m.settingsSelection {
			style = selectedStyle
			prefix = "> "
		}
		content.WriteString(style.Render(prefix + item.name))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(mutedStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("j/k:nav  enter:select  esc:back"))

	return m.centerModal(modalStyle.Width(modalWidth).Render(content.String()))
}

func (m *Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.switchScreen(screenMain)
			return m, nil

		case "j", "down":
			if m.settingsSelection < settingsReset {
				m.settingsSelection++
			}
			return m, nil

		case "k", "up":
			if m.settingsSelection > settingsSave {
				m.settingsSelection--
			}
			return m, nil

		case "enter":
			return m.handleSettingsSelection()

		case "?":
			m.switchScreen(screenHelp)
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) handleSettingsSelection() (tea.Model, tea.Cmd) {
	switch m.settingsSelection {
	case settingsSave:
		if m.currentPresetName != "" && m.currentPresetName != "_last" {
			preset := m.buildCurrentPreset()
			config.SavePreset(preset, m.currentPresetName)
			m.loadedPresetHash = preset.Hash()
			m.isDirty = false
			m.switchScreen(screenMain)
		} else {
			// Show save-as dialog
			m.presetBrowser.mode = browserModeSaveAs
			m.presetBrowser.inputBuffer = ""
			m.refreshPresetList()
			m.switchScreen(screenPresetBrowser)
		}

	case settingsSaveAs:
		m.presetBrowser.mode = browserModeSaveAs
		m.presetBrowser.inputBuffer = ""
		m.refreshPresetList()
		m.switchScreen(screenPresetBrowser)

	case settingsLoad:
		m.refreshPresetList()
		m.switchScreen(screenPresetBrowser)

	case settingsReset:
		m.applyPreset(config.DefaultPreset())
		m.currentPresetName = ""
		m.isDirty = true
	}
	return m, nil
}
