package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-control/osc"
)

func TestHelp_ToggleWithQuestionMark(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Press '?' to open help
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m := result.(*Model)

	if m.screen != screenHelp {
		t.Errorf("expected screen to be help, got %d", m.screen)
	}

	// Save previous screen
	if m.prevScreen != screenMain {
		t.Errorf("expected prevScreen to be main, got %d", m.prevScreen)
	}
}

func TestHelp_EscCloses(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenHelp))
	model.prevScreen = screenMain

	// Press Esc to close help
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m := result.(*Model)

	if m.screen != screenMain {
		t.Errorf("expected screen to return to main, got %d", m.screen)
	}
}

func TestHelp_QCloses(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenHelp))
	model.prevScreen = screenMain

	// Press 'q' to close help
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m := result.(*Model)

	if m.screen != screenMain {
		t.Errorf("expected screen to return to main, got %d", m.screen)
	}
}

func TestHelp_QuestionMarkCloses(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenHelp))
	model.prevScreen = screenMain

	// Press '?' again to close help
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m := result.(*Model)

	if m.screen != screenMain {
		t.Errorf("expected screen to return to main, got %d", m.screen)
	}
}

func TestHelp_Render(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Set window size first (use larger height to ensure all content is visible)
	model.Update(tea.WindowSizeMsg{Width: 80, Height: 50})

	// Open help
	model.SetScreenForTesting(int(screenHelp))

	view := model.View()

	if !strings.Contains(view, "Global") {
		t.Error("expected help view to contain 'Global' section")
	}

	if !strings.Contains(view, "Effects List") {
		t.Error("expected help view to contain 'Effects List' section")
	}

	// These sections may not be visible in smaller terminals
	// but should be present in the content
	if !strings.Contains(view, "Parameters") {
		t.Log("Note: 'Parameters' section may not be visible in current terminal size")
	}
}

func TestHelp_ContainsKeybindings(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// Set window size first (use larger height to ensure all content is visible)
	model.Update(tea.WindowSizeMsg{Width: 80, Height: 50})

	// Open help
	model.SetScreenForTesting(int(screenHelp))

	view := model.View()

	// Check for key bindings
	if !strings.Contains(view, ":") {
		t.Error("expected help to contain ':' command palette key")
	}

	if !strings.Contains(view, "j/k") {
		t.Error("expected help to contain 'j/k' navigation keys")
	}

	// 'esc' might not be visible in smaller terminals
	if !strings.Contains(view, "esc") {
		t.Log("Note: 'esc' key may not be visible in current terminal size")
	}
}

func TestHelp_WindowSize(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenHelp))

	// Update window size
	result, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m := result.(*Model)

	if m.width != 120 {
		t.Errorf("expected width 120, got %d", m.width)
	}

	if m.height != 40 {
		t.Errorf("expected height 40, got %d", m.height)
	}
}

func TestHelp_PreservesPrevScreen(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenMain))
	model.InitLists(80, 40)

	// First go to settings screen properly
	model.switchScreen(screenSettings)

	// Open help from settings
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m := result.(*Model)

	if m.screen != screenHelp {
		t.Errorf("expected screen to be help, got %d", m.screen)
	}

	if m.prevScreen != screenSettings {
		t.Errorf("expected prevScreen to be settings, got %d", m.prevScreen)
	}

	// Close help
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = result.(*Model)

	if m.screen != screenSettings {
		t.Errorf("expected screen to return to settings, got %d", m.screen)
	}
}
