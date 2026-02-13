package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1)

	connectedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	disconnectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
)

func (m Model) View() string {
	if len(m.effectsList.Items()) == 0 && m.width > 0 {
		return "Loading..."
	}

	var listView string
	switch m.navigationMode {
	case modeEffectsList:
		listView = m.effectsList.View()
	case modeParameterList:
		listView = m.parameterList.View()
	default:
		listView = m.effectsList.View()
	}

	if !m.showStatus {
		return appStyle.Render(listView)
	}

	connectionStatus := disconnectedStyle.Render("Disconnected")
	if m.connected {
		connectionStatus = connectedStyle.Render("Connected")
	}

	midiStatus := m.midiPort
	if midiStatus == "" {
		midiStatus = "No MIDI"
	}

	pendingCount := len(m.pendingChanges)
	pendingStatus := ""
	if pendingCount > 0 {
		pendingStatus = fmt.Sprintf(" | %d pending", pendingCount)
	}

	statusBar := statusBarStyle.Render(
		fmt.Sprintf("%s | MIDI: %s%s", connectionStatus, midiStatus, pendingStatus),
	)

	return appStyle.Render(listView) + "\n" + statusBar
}

func formatValue(value, min, max float32) string {
	if max-min >= 10 {
		return fmt.Sprintf("%5.1f", value)
	}
	return fmt.Sprintf("%5.2f", value)
}
