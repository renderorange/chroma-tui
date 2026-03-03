package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-control/osc"
)

func TestPalette_Toggle(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Initially palette should be closed
	if model.showCommandPalette {
		t.Error("expected palette to be initially closed")
	}

	// Toggle palette on with ':' key
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	m := result.(*Model)

	if !m.showCommandPalette {
		t.Error("expected palette to be open after ':' key")
	}

	// Toggle palette off again
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	m = result.(*Model)

	if m.showCommandPalette {
		t.Error("expected palette to be closed after second ':' key")
	}
}

func TestPalette_EscCloses(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Open palette
	model.showCommandPalette = true
	model.commandPaletteText = ":quit"

	// Press Esc
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m := result.(*Model)

	if m.showCommandPalette {
		t.Error("expected palette to be closed after Esc")
	}

	if m.commandPaletteText != "" {
		t.Error("expected palette text to be cleared")
	}
}

func TestPalette_QuitCommand(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Open palette and type quit
	model.showCommandPalette = true
	model.commandPaletteText = ":quit"

	// Press Enter
	result, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = result.(*Model)

	// Should return quit command
	if cmd == nil {
		t.Error("expected quit command, got nil")
	}
}

func TestPalette_HelpCommand(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Open palette and type help
	model.showCommandPalette = true
	model.commandPaletteText = ":help"

	// Press Enter
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := result.(*Model)

	// Should transition to help screen
	if m.screen != screenHelp {
		t.Errorf("expected screen to be help, got %d", m.screen)
	}

	// Palette should be closed
	if m.showCommandPalette {
		t.Error("expected palette to be closed after executing command")
	}
}

func TestPalette_SettingsCommand(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Open palette and type settings
	model.showCommandPalette = true
	model.commandPaletteText = ":settings"

	// Press Enter
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := result.(*Model)

	// Should transition to settings screen
	if m.screen != screenSettings {
		t.Errorf("expected screen to be settings, got %d", m.screen)
	}
}

func TestPalette_FuzzyMatch(t *testing.T) {
	// Test fuzzy matching for commands
	matches := fuzzyMatchCommands("qu")

	foundQuit := false
	for _, match := range matches {
		if match.command.Name == "quit" {
			foundQuit = true
			break
		}
	}

	if !foundQuit {
		t.Error("expected 'qu' to match 'quit' command")
	}
}

func TestPalette_UnknownCommand(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Open palette and type unknown command
	model.showCommandPalette = true
	model.commandPaletteText = ":unknown"

	// Press Enter
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := result.(*Model)

	// Should show error message
	if !strings.Contains(m.commandPaletteText, "Unknown") {
		t.Error("expected error message for unknown command")
	}
}
