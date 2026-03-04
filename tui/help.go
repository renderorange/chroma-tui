package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type helpSection struct {
	Title string
	Items []helpItem
}

type helpItem struct {
	Key         string
	Description string
}

func helpContent() []helpSection {
	return []helpSection{
		{
			Title: "Global",
			Items: []helpItem{
				{Key: ":", Description: "Open command palette"},
				{Key: "?", Description: "Toggle help"},
				{Key: "ctrl+c", Description: "Quit"},
			},
		},
		{
			Title: "Effects List",
			Items: []helpItem{
				{Key: "j/k", Description: "Navigate effects"},
				{Key: "enter", Description: "Open parameters"},
				{Key: "1/2/3", Description: "Set blend mode"},
			},
		},
		{
			Title: "Parameters",
			Items: []helpItem{
				{Key: "j/k", Description: "Navigate parameters"},
				{Key: "h/l", Description: "Adjust value"},
				{Key: "enter", Description: "Toggle/cycle"},
				{Key: "esc", Description: "Back to effects"},
			},
		},
		{
			Title: "Effects Order",
			Items: []helpItem{
				{Key: "enter", Description: "Toggle reorder mode"},
				{Key: "j/k", Description: "Select effect"},
				{Key: "h/l", Description: "Move effect"},
			},
		},
		{
			Title: "Commands",
			Items: []helpItem{
				{Key: ":quit", Description: "Exit application"},
				{Key: "help/h/?", Description: "Show help"},
				{Key: "settings/set", Description: "Open settings"},
			},
		},
	}
}

func (m *Model) renderHelp() string {
	return m.renderHelpModal()
}

func (m *Model) renderHelpModal() string {
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

	sectionStyle := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Width(modalWidth - 4)

	keyStyle := lipgloss.NewStyle().
		Foreground(colorTextNormal).
		Width(16)

	descStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted).
		Width(modalWidth - 20)

	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2)

	var content strings.Builder

	content.WriteString(titleStyle.Render("Help"))
	content.WriteString("\n\n")

	for i, section := range helpContent() {
		if i > 0 {
			content.WriteString("\n")
		}
		content.WriteString(sectionStyle.Render(section.Title))
		content.WriteString("\n")

		for _, item := range section.Items {
			key := keyStyle.Render(item.Key)
			desc := descStyle.Render(item.Description)
			content.WriteString(key)
			content.WriteString(desc)
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(mutedStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n")
	content.WriteString(mutedStyle.Render("esc:back"))

	return m.centerModal(modalStyle.Width(modalWidth).Render(content.String()))
}

func (m *Model) overlayModal(background, modal string) string {
	bgLines := strings.Split(background, "\n")
	modalLines := strings.Split(modal, "\n")

	modalHeight := len(modalLines)

	verticalPadding := (m.height - modalHeight) / 2
	if verticalPadding < 1 {
		verticalPadding = 1
	}

	var result strings.Builder

	for i := 0; i < verticalPadding && i < len(bgLines); i++ {
		result.WriteString(bgLines[i])
		result.WriteString("\n")
	}

	availableWidth := m.width
	if availableWidth <= 0 {
		availableWidth = 80
	}

	for _, modalLine := range modalLines {
		lineWidth := lipgloss.Width(modalLine)
		leftPadding := (availableWidth - lineWidth) / 2
		if leftPadding < 0 {
			leftPadding = 0
		}

		result.WriteString(strings.Repeat(" ", leftPadding))
		result.WriteString(modalLine)
		result.WriteString("\n")
	}

	for i := verticalPadding + modalHeight; i < len(bgLines); i++ {
		result.WriteString(bgLines[i])
		if i < len(bgLines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
