package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestUpdate_KeyboardNavigation(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Test Enter key - should enter filter section
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)

	updatedModelTyped := updatedModel.(*Model)
	if updatedModelTyped.navigationMode != modeParameterList {
		t.Errorf("expected to enter parameter list mode, got %d", updatedModelTyped.navigationMode)
	}

	// Test Esc key - should go back to effects list
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ = updatedModelTyped.Update(escMsg)

	updatedModelTyped = updatedModel.(*Model)
	if updatedModelTyped.navigationMode != modeEffectsList {
		t.Errorf("expected to return to effects list mode, got %d", updatedModelTyped.navigationMode)
	}
}

func TestUpdate_ArrowKeyNavigation(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Test navigation within effects list - list handles this internally
	// The list should be navigable
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)

	_ = updatedModel.(*Model)
	// List navigation doesn't change our state directly
}

func TestUpdate_ParameterAdjustment(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Enter filter section
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedEnter, _ := model.Update(enterMsg)
	model = *updatedEnter.(*Model)

	// Test Right arrow (should adjust parameter when in parameter mode)
	msg := tea.KeyMsg{Type: tea.KeyRight}
	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("expected no command when adjusting parameter")
	}

	_ = updatedModel.(*Model)
	// Parameter adjustment happens internally
}

func TestUpdate_ToggleControls(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Test toggle visibility keys
	testCases := []struct {
		name     string
		key      tea.KeyMsg
		getValue func(*Model) bool
	}{
		{
			name:     "Toggle Help",
			key:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("H")},
			getValue: func(m *Model) bool { return m.showHelp },
		},
		{
			name:     "Toggle Status",
			key:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("S")},
			getValue: func(m *Model) bool { return m.showStatus },
		},
		{
			name:     "Toggle Pagination",
			key:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("P")},
			getValue: func(m *Model) bool { return m.showPagination },
		},
		{
			name:     "Toggle Title",
			key:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("T")},
			getValue: func(m *Model) bool { return m.showTitle },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialValue := tc.getValue(&model)

			updatedModel, _ := model.Update(tc.key)

			updatedModelTyped := updatedModel.(*Model)
			newValue := tc.getValue(updatedModelTyped)

			if initialValue == newValue {
				t.Errorf("expected toggle to change value")
			}
		})
	}
}

func TestUpdate_Quit(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	updatedModel, cmd := model.Update(msg)

	updatedModelTyped := updatedModel.(*Model)
	if cmd == nil {
		t.Error("expected Quit command, got nil")
	}
	_ = updatedModelTyped
}
