package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderSplash creates the centered splash screen view.
func (m *Model) renderSplash() string {

	// Get terminal dimensions
	w, h := m.width, m.height
	if w == 0 || h == 0 {
		w, h = 80, 24 // reasonable defaults
	}

	// Create styles
	primaryStyle := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true)

	secondaryStyle := lipgloss.NewStyle().
		Foreground(colorSecondary)

	mutedStyle := lipgloss.NewStyle().
		Foreground(colorTextMuted)

	// Content
	appName := "chroma [control]"
	version := m.version
	if version == "" {
		version = "dev"
	}
	prompt := "Press any key..."

	// Calculate vertical centering
	contentHeight := 5 // appName + version + spacing + prompt
	verticalPadding := (h - contentHeight) / 2
	if verticalPadding < 1 {
		verticalPadding = 1
	}

	// Build the content
	var b strings.Builder

	// Top padding
	for i := 0; i < verticalPadding; i++ {
		b.WriteString("\n")
	}

	// App name (centered)
	nameWidth := lipgloss.Width(appName)
	namePadding := (w - nameWidth) / 2
	if namePadding < 0 {
		namePadding = 0
	}
	b.WriteString(strings.Repeat(" ", namePadding))
	b.WriteString(primaryStyle.Render(appName))
	b.WriteString("\n")

	// Version (centered)
	versionWidth := lipgloss.Width(version)
	versionPadding := (w - versionWidth) / 2
	if versionPadding < 0 {
		versionPadding = 0
	}
	b.WriteString(strings.Repeat(" ", versionPadding))
	b.WriteString(secondaryStyle.Render(version))
	b.WriteString("\n\n")

	// Prompt (centered)
	promptWidth := lipgloss.Width(prompt)
	promptPadding := (w - promptWidth) / 2
	if promptPadding < 0 {
		promptPadding = 0
	}
	b.WriteString(strings.Repeat(" ", promptPadding))
	b.WriteString(mutedStyle.Render(prompt))

	return b.String()
}
