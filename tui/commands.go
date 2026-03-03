package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Command represents an executable command.
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Handler     func(m *Model, args []string) tea.Cmd
}

// availableCommands returns all available commands.
func availableCommands() []Command {
	return []Command{
		{
			Name:        "quit",
			Aliases:     []string{},
			Description: "Exit the application",
			Handler:     cmdQuit,
		},
		{
			Name:        "help",
			Aliases:     []string{"h", "?"},
			Description: "Open help screen",
			Handler:     cmdHelp,
		},
		{
			Name:        "settings",
			Aliases:     []string{"set"},
			Description: "Open settings screen",
			Handler:     cmdSettings,
		},
	}
}

// findCommand finds a command by name or alias.
func findCommand(name string) (Command, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, cmd := range availableCommands() {
		if cmd.Name == name {
			return cmd, true
		}
		for _, alias := range cmd.Aliases {
			if alias == name {
				return cmd, true
			}
		}
	}
	return Command{}, false
}

// cmdQuit handles the quit command.
func cmdQuit(m *Model, args []string) tea.Cmd {
	return tea.Quit
}

// cmdHelp handles the help command.
func cmdHelp(m *Model, args []string) tea.Cmd {
	m.switchScreen(screenHelp)
	return nil
}

// cmdSettings handles the settings command.
func cmdSettings(m *Model, args []string) tea.Cmd {
	m.switchScreen(screenSettings)
	return nil
}
