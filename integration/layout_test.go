package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

func Test_Layout(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57130)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Test 1: Footer shows effects list keybindings
	view := model.View()
	if !strings.Contains(view, "enter: open params") {
		t.Error("expected footer to show 'enter: open params' in effects list mode")
	}

	// Test 2: Navigate and press Enter
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := result.(*Model)
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = result.(*Model)

	view = m.View()
	if !strings.Contains(view, "h/l: adjust value") {
		t.Error("expected footer to show 'h/l: adjust value' in parameter list mode")
	}

	// Test 3: Global section footer - navigate to global in effects list
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc}) // Go back to effects list
	m = result.(*Model)

	// Navigate to global (index 7 in the effects list)
	for i := 0; i < 7; i++ {
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = result.(*Model)
	}

	// Enter global section
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = result.(*Model)

	view = m.View()
	if !strings.Contains(view, "pgup/pgdn: reorder") {
		t.Errorf("expected footer to show 'pgup/pgdn: reorder' in global section")
	}

	// Test 4: Smoke test for rendering
	if len(view) == 0 {
		t.Error("expected non-empty view with panel spacing")
	}

	// Test 5: Wide terminal slider scaling
	result, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = result.(*Model)
	view = m.View()
	if !strings.Contains(view, "=") && !strings.Contains(view, "-") {
		t.Error("expected view to contain slider bars")
	}

	// Test 6: Minimum terminal size warning
	result, _ = m.Update(tea.WindowSizeMsg{Width: 50, Height: 15})
	m = result.(*Model)
	view = m.View()
	if !strings.Contains(view, "Terminal too small") {
		t.Errorf("expected warning for terminal below minimum size, got: %s", view)
	}
}
