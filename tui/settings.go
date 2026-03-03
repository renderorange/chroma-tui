package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsItem struct {
	name        string
	description string
	isSelected  bool
}

func (m *Model) renderSettings() string {
	background := m.renderMain()
	modalContent := m.renderSettingsModal()
	return m.overlayModal(background, modalContent)
}

func (m *Model) renderSettingsModal() string {
	modalWidth := 50
	if m.width > 0 && m.width < modalWidth+4 {
		modalWidth = m.width - 4
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	itemStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted).
		Width(modalWidth - 4)

	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 1)

	var content strings.Builder

	content.WriteString(itemStyle.Render("Settings coming soon..."))

	content.WriteString("\n\n")
	content.WriteString(mutedStyle.Render(strings.Repeat("─", modalWidth-2)))
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("esc: back"))

	return modalStyle.Width(modalWidth).Render(content.String())
}

func (m *Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.goBack()
			return m, nil

		case "?":
			m.switchScreen(screenHelp)
			return m, nil
		}
	}
	return m, nil
}
