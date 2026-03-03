package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-control/osc"
)

func TestSplash_Render(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetVersion("1.0.0")
	model.SetScreenForTesting(int(screenSplash))

	// Set window size
	model.width = 80
	model.height = 24

	view := model.View()

	// Verify splash screen content
	if !strings.Contains(view, "chroma [control]") {
		t.Error("expected splash to contain app name")
	}

	if !strings.Contains(view, "1.0.0") {
		t.Error("expected splash to contain version")
	}

	if !strings.Contains(view, "Press any key") {
		t.Error("expected splash to contain 'Press any key' prompt")
	}
}

func TestSplash_KeyPressTransitionsToMain(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenSplash))
	model.InitLists(80, 40)

	// Press any key (using 'enter' as example)
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := result.(*Model)

	if m.screen != screenMain {
		t.Errorf("expected screen to transition to main, got %d", m.screen)
	}
}

func TestSplash_EscQuits(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenSplash))

	result, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = result.(*Model)

	// Should return Quit command
	if cmd == nil {
		t.Error("expected Quit command on Esc, got nil")
		return
	}

	// Verify it's a quit command by checking if it returns tea.Quit()
	// We can't easily test the command itself, but we can verify it's not nil
}

func TestSplash_CtrlCQuits(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenSplash))

	result, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	_ = result.(*Model)

	if cmd == nil {
		t.Error("expected Quit command on Ctrl+C, got nil")
	}
}

func TestSplash_WindowSize(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetScreenForTesting(int(screenSplash))

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

func TestSplash_DefaultVersion(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	// Don't set version - should default to empty or "dev"
	model.SetScreenForTesting(int(screenSplash))
	model.width = 80
	model.height = 24

	view := model.View()

	// Should still render without panicking
	if !strings.Contains(view, "chroma [control]") {
		t.Error("expected splash to render even without version set")
	}
}
