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

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#25A065"))

	unfocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#555555"))
)

func (m Model) View() string {
	if len(m.effectsList.Items()) == 0 && m.width > 0 {
		return "Loading..."
	}

	leftView := m.effectsList.View()
	rightView := m.parameterList.View()

	// Apply border styles based on focus
	var leftPanel, rightPanel string
	if m.navigationMode == modeEffectsList {
		leftPanel = focusedBorderStyle.Render(leftView)
		rightPanel = unfocusedBorderStyle.Render(rightView)
	} else {
		leftPanel = unfocusedBorderStyle.Render(leftView)
		rightPanel = focusedBorderStyle.Render(rightView)
	}

	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	if !m.showStatus {
		return appStyle.Render(panels)
	}

	connectionStatus := disconnectedStyle.Render("Disconnected")
	if m.connected {
		connectionStatus = connectedStyle.Render("Connected")
	}

	midiStatus := m.midiPort
	if midiStatus == "" {
		midiStatus = "No MIDI"
	}

	statusBar := statusBarStyle.Render(
		fmt.Sprintf("%s | MIDI: %s", connectionStatus, midiStatus),
	)

	return appStyle.Render(panels) + "\n" + statusBar
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
