package integration

import (
	"github.com/renderorange/chroma/chroma-control/tui"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-control/osc"
)

func Test_Layout(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57130)
	model := tui.NewModel(client)
	model.InitLists(80, 40)
	model.SetScreenForTesting(1) // Set to main screen for testing

	// Test 1: Footer shows effects list keybindings
	view := model.View()
	if !strings.Contains(view, "enter:params") {
		t.Error("expected footer to show 'enter:params' in effects list mode")
	}

	// Test 2: Navigate and press Enter
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := result.(*tui.Model)
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = result.(*tui.Model)

	view = m.View()
	if !strings.Contains(view, "h/l:adjust") {
		t.Error("expected footer to show 'h/l:adjust' in parameter list mode")
	}

	// Test 3: Master section footer - navigate to master in effects list
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc}) // Go back to effects list
	m = result.(*tui.Model)

	// Navigate to master (index 0 in the effects list) - need to go up from filter
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = result.(*tui.Model)

	// Enter master section
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = result.(*tui.Model)

	view = m.View()
	if !strings.Contains(view, "h/l:adjust") {
		t.Errorf("expected footer to show 'h/l:adjust' in master section")
	}

	// Test 3b: Navigate to Effects Order and check for reorder controls
	// Effects Order is the last item in Master section (index 6)
	for i := 0; i < 6; i++ {
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = result.(*tui.Model)
	}

	view = m.View()
	if !strings.Contains(view, "pgup/pgdn:reorder") {
		t.Errorf("expected footer to show 'pgup/pgdn:reorder' when Effects Order is selected")
	}

	// Test 4: Smoke test for rendering
	if len(view) == 0 {
		t.Error("expected non-empty view with panel spacing")
	}

	// Test 5: Wide terminal slider scaling
	result, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = result.(*tui.Model)
	view = m.View()
	if !strings.Contains(view, "=") && !strings.Contains(view, "-") {
		t.Error("expected view to contain slider bars")
	}

	// Test 6: Minimum terminal size warning
	result, _ = m.Update(tea.WindowSizeMsg{Width: 50, Height: 15})
	m = result.(*tui.Model)
	view = m.View()
	if !strings.Contains(view, "Terminal too small") {
		t.Errorf("expected warning for terminal below minimum size, got: %s", view)
	}
}
