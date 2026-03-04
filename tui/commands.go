package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-control/config"
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
			Aliases:     []string{"exit"},
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
		{
			Name:        "save",
			Aliases:     []string{},
			Description: "Save current preset (optionally with name: :save my-preset)",
			Handler:     cmdSave,
		},
		{
			Name:        "load",
			Aliases:     []string{},
			Description: "Load a preset by name (:load my-preset) or open browser",
			Handler:     cmdLoad,
		},
		{
			Name:        "presets",
			Aliases:     []string{"browser"},
			Description: "Open preset browser",
			Handler:     cmdPresets,
		},
		{
			Name:        "reset",
			Aliases:     []string{"defaults"},
			Description: "Reset all settings to factory defaults",
			Handler:     cmdReset,
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
	return func() tea.Msg {
		if m.isDirty {
			m.showQuitConfirm = true
			// Return a message to trigger re-render
			type redrawMsg struct{}
			return redrawMsg{}
		}
		return tea.Quit()
	}
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

// cmdSave handles the save command.
func cmdSave(m *Model, args []string) tea.Cmd {
	return func() tea.Msg {
		if len(args) > 0 {
			name := strings.Join(args, " ")
			preset := m.buildCurrentPreset()
			if err := config.SavePreset(preset, name); err == nil {
				m.currentPresetName = name
				config.SaveLastPresetName(name)
				m.loadedPresetHash = preset.Hash()
				m.isDirty = false
			}
		} else if m.currentPresetName != "" && m.currentPresetName != "_last" {
			preset := m.buildCurrentPreset()
			config.SavePreset(preset, m.currentPresetName)
			m.loadedPresetHash = preset.Hash()
			m.isDirty = false
		} else {
			// Open save-as dialog
			m.presetBrowser.mode = browserModeSaveAs
			m.presetBrowser.inputBuffer = ""
			m.refreshPresetList()
			m.switchScreen(screenPresetBrowser)
		}
		return nil
	}
}

// cmdLoad handles the load command.
func cmdLoad(m *Model, args []string) tea.Cmd {
	if len(args) > 0 {
		return func() tea.Msg {
			name := strings.Join(args, " ")
			if preset, err := config.LoadPreset(name); err == nil {
				m.applyPreset(preset)
				m.currentPresetName = name
				config.SaveLastPresetName(name)
			}
			return nil
		}
	}
	// Open preset browser immediately
	m.refreshPresetList()
	m.switchScreen(screenPresetBrowser)
	return nil
}

// cmdPresets handles the presets command.
func cmdPresets(m *Model, args []string) tea.Cmd {
	m.refreshPresetList()
	m.switchScreen(screenPresetBrowser)
	return nil
}

// cmdReset handles the reset command.
func cmdReset(m *Model, args []string) tea.Cmd {
	m.applyPreset(config.DefaultPreset())
	m.currentPresetName = ""
	m.isDirty = true
	return nil
}
