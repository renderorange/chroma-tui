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

func TestUpdate_ListNavigationDelegation(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Effects list should start at index 0
	if model.effectsList.Index() != 0 {
		t.Fatalf("expected initial index 0, got %d", model.effectsList.Index())
	}

	// Press down arrow - should move to index 1
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*Model)

	if m.effectsList.Index() != 1 {
		t.Errorf("expected effects list index 1 after down arrow, got %d", m.effectsList.Index())
	}
}

func TestUpdate_ParameterPanelSyncsWithEffectSelection(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Start at index 0 (input), move down to index 1 (filter)
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*Model)

	// Current section should now be "filter"
	if m.currentSection != "filter" {
		t.Errorf("expected currentSection 'filter', got '%s'", m.currentSection)
	}

	// Parameter list should have filter parameters (4 items: enabled, amount, cutoff, resonance)
	items := m.parameterList.Items()
	if len(items) != 4 {
		t.Errorf("expected 4 filter parameters, got %d", len(items))
	}
}

func TestUpdate_FullSideBySideWorkflow(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(120, 40)

	// 1. Navigate down to "filter" in effects list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := model.Update(downMsg)
	m := result.(*Model)

	// Parameter panel should show filter params
	if m.currentSection != "filter" {
		t.Errorf("expected section 'filter', got '%s'", m.currentSection)
	}

	// 2. Press Enter to focus parameter panel
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(*Model)

	if m.navigationMode != modeParameterList {
		t.Errorf("expected parameter list mode")
	}

	// 3. Navigate down to Amount parameter (index 1, skipping the enabled toggle at index 0)
	result, _ = m.Update(downMsg)
	m = result.(*Model)

	// 4. Adjust parameter with right arrow
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	initialAmount := m.FilterAmount
	result, _ = m.Update(rightMsg)
	m = result.(*Model)

	if m.FilterAmount <= initialAmount {
		t.Errorf("expected FilterAmount to increase, was %f now %f", initialAmount, m.FilterAmount)
	}

	// 5. Press Esc to go back to effects panel
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	result, _ = m.Update(escMsg)
	m = result.(*Model)

	if m.navigationMode != modeEffectsList {
		t.Errorf("expected effects list mode")
	}

	// 6. Effects list cursor should still be on filter (index 1)
	if m.effectsList.Index() != 1 {
		t.Errorf("expected effects list index 1 (filter), got %d", m.effectsList.Index())
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
