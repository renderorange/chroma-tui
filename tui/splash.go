package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/renderorange/chroma/chroma-control/config"
)

func (m *Model) renderSplash() string {
	w, h := m.width, m.height
	if w == 0 || h == 0 {
		w, h = 80, 24
	}

	primaryStyle := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	secondaryStyle := lipgloss.NewStyle().Foreground(colorSecondary)
	mutedStyle := lipgloss.NewStyle().Foreground(colorTextMuted)
	highlightStyle := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)

	appName := "chroma [control]"
	version := m.version
	if version == "" {
		version = "dev"
	}

	// Build horizontal selector like intensity setting
	options := []string{"last", "new", "load", "quit"}
	var optionParts []string
	for i, opt := range options {
		if i == int(m.splashSelection) {
			optionParts = append(optionParts, highlightStyle.Render("["+opt+"]"))
		} else {
			optionParts = append(optionParts, mutedStyle.Render(opt))
		}
	}
	optionsLine := strings.Join(optionParts, "  ")

	// Content height: appName + version + 2 spaces + selector
	contentHeight := 5
	verticalPadding := (h - contentHeight) / 2
	if verticalPadding < 1 {
		verticalPadding = 1
	}

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

	// Selector (centered)
	optionsWidth := lipgloss.Width(optionsLine)
	optionsPadding := (w - optionsWidth) / 2
	if optionsPadding < 0 {
		optionsPadding = 0
	}
	b.WriteString(strings.Repeat(" ", optionsPadding))
	b.WriteString(optionsLine)

	return b.String()
}

func (m *Model) updateSplash(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "left", "h":
			if m.splashSelection > 0 {
				m.splashSelection--
			}
			return m, nil

		case "right", "l":
			if m.splashSelection < splashQuit {
				m.splashSelection++
			}
			return m, nil

		case "enter":
			return m.handleSplashSelection()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m *Model) handleSplashSelection() (tea.Model, tea.Cmd) {
	// Initialize lists first if we have dimensions
	if m.width > 0 && m.effectsList.Items() == nil {
		m.InitLists(m.width, m.height)
	}

	switch m.splashSelection {
	case splashLast:
		// Try to load autosave first, otherwise use last preset name
		if autosave, err := config.LoadAutosave(); err == nil {
			m.applyPreset(autosave)
			m.currentPresetName = config.LoadLastPresetName()
		} else if lastName := config.LoadLastPresetName(); lastName != "" {
			if preset, err := config.LoadPreset(lastName); err == nil {
				m.applyPreset(preset)
				m.currentPresetName = lastName
			} else {
				m.applyPreset(config.DefaultPreset())
				m.currentPresetName = ""
			}
		} else {
			m.applyPreset(config.DefaultPreset())
			m.currentPresetName = ""
		}
		m.switchScreen(screenMain)

	case splashNew:
		// Start with factory defaults
		m.applyPreset(config.DefaultPreset())
		m.currentPresetName = ""
		m.switchScreen(screenMain)

	case splashLoad:
		// Go to preset browser
		m.refreshPresetList()
		m.switchScreen(screenPresetBrowser)

	case splashQuit:
		// Quit the application
		return m, tea.Quit
	}

	return m, nil
}
