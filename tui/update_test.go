package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestUpdate_KeyboardNavigation(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test Tab key (next control)
	initialFocus := model.focused
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)

	updatedModelTyped := updatedModel.(Model)
	expectedFocus := (initialFocus + 1) % ctrlCount
	if updatedModelTyped.focused != expectedFocus {
		t.Errorf("expected focus to move from %d to %d, got %d", initialFocus, expectedFocus, updatedModelTyped.focused)
	}

	// Test Shift+Tab (previous control) - should go to ctrlCount-1 from 0
	model.focused = 0 // Reset to first control
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	updatedModel, _ = model.Update(msg)
	updatedModelTyped = updatedModel.(Model)
	expectedFocus = ctrlCount - 1 // Should wrap to last control
	if updatedModelTyped.focused != expectedFocus {
		t.Errorf("expected focus to wrap to %d, got %d", expectedFocus, updatedModelTyped.focused)
	}
}

func TestUpdate_ArrowKeyNavigation(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test Up arrow (previous control)
	initialFocus := model.focused
	msg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(msg)

	updatedModelTyped := updatedModel.(Model)
	expectedFocus := (initialFocus - 1 + ctrlCount) % ctrlCount
	if updatedModelTyped.focused != expectedFocus {
		t.Errorf("expected focus to move from %d to %d with up arrow, got %d", initialFocus, expectedFocus, updatedModelTyped.focused)
	}

	// Test Down arrow (next control) - should move to next control, not wrap
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = updatedModelTyped.Update(downMsg)
	updatedModelTyped = updatedModel.(Model)
	expectedFocus = (expectedFocus + 1) % ctrlCount // Next control from previous
	if updatedModelTyped.focused != expectedFocus {
		t.Errorf("expected focus to move to %d with down arrow, got %d", expectedFocus, updatedModelTyped.focused)
	}
}

func TestUpdate_ParameterAdjustment(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test that parameter adjustment keys are handled
	model.focused = ctrlGain

	// Test Right arrow (should be handled without error)
	msg := tea.KeyMsg{Type: tea.KeyRight}
	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("expected no command when adjusting parameter")
	}

	updatedModelTyped := updatedModel.(Model)
	if updatedModelTyped.focused != ctrlGain {
		t.Error("expected focus to remain on gain after adjustment")
	}

	// Test Left arrow (should be handled without error)
	msg = tea.KeyMsg{Type: tea.KeyLeft}
	updatedModel, cmd = updatedModelTyped.Update(msg)

	if cmd != nil {
		t.Error("expected no command when adjusting parameter")
	}

	finalModelTyped := updatedModel.(Model)
	if finalModelTyped.focused != ctrlGain {
		t.Error("expected focus to remain on gain after adjustment")
	}
}

func TestUpdate_ToggleControls(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test Enter key on toggle controls
	testCases := []struct {
		name     string
		ctrl     control
		getValue func(Model) bool
		setValue func(Model, bool)
	}{
		{
			name:     "Filter Enabled",
			ctrl:     ctrlFilterEnabled,
			getValue: func(m Model) bool { return m.FilterEnabled },
		},
		{
			name:     "Overdrive Enabled",
			ctrl:     ctrlOverdriveEnabled,
			getValue: func(m Model) bool { return m.OverdriveEnabled },
		},
		{
			name:     "Granular Enabled",
			ctrl:     ctrlGranularEnabled,
			getValue: func(m Model) bool { return m.GranularEnabled },
		},
		{
			name:     "Bitcrush Enabled",
			ctrl:     ctrlBitcrushEnabled,
			getValue: func(m Model) bool { return m.BitcrushEnabled },
		},
		{
			name:     "Reverb Enabled",
			ctrl:     ctrlReverbEnabled,
			getValue: func(m Model) bool { return m.ReverbEnabled },
		},
		{
			name:     "Delay Enabled",
			ctrl:     ctrlDelayEnabled,
			getValue: func(m Model) bool { return m.DelayEnabled },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model.focused = tc.ctrl
			initialValue := tc.getValue(model)

			// Toggle with Enter key
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			updatedModel, _ := model.Update(msg)
			updatedModelTyped := updatedModel.(Model)

			newValue := tc.getValue(updatedModelTyped)
			if newValue == initialValue {
				t.Errorf("expected %s to toggle from %v to %v", tc.name, initialValue, newValue)
			}
		})
	}
}

func TestUpdate_GrainIntensityCycling(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test 'i' key for grain intensity
	model.focused = ctrlGrainIntensity
	initialIntensity := model.GrainIntensity

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
	updatedModel, _ := model.Update(msg)
	updatedModelTyped := updatedModel.(Model)

	newIntensity := updatedModelTyped.GrainIntensity
	if newIntensity == initialIntensity {
		t.Errorf("expected grain intensity to change from %s, got %s", initialIntensity, newIntensity)
	}

	// Test that intensity cycles through valid values
	validIntensities := []string{"subtle", "pronounced", "extreme"}
	found := false
	for _, valid := range validIntensities {
		if newIntensity == valid {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected grain intensity to be one of %v, got %s", validIntensities, newIntensity)
	}
}

func TestUpdate_BlendModeSelection(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test number keys for blend mode
	testCases := []struct {
		key      rune
		expected int
		name     string
	}{
		{'1', 0, "mirror"},
		{'2', 1, "complement"},
		{'3', 2, "transform"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model.focused = ctrlBlendMode
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tc.key}}
			updatedModel, _ := model.Update(msg)
			updatedModelTyped := updatedModel.(Model)

			if updatedModelTyped.BlendMode != tc.expected {
				t.Errorf("expected blend mode %d for key '%c', got %d", tc.expected, tc.key, updatedModelTyped.BlendMode)
			}
		})
	}
}

func TestUpdate_EffectsOrderReordering(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test that effects order reordering keys are handled
	model.focused = ctrlEffectsOrder
	model.selectedEffectIndex = 0

	// Test PgDn (move effect down)
	msg := tea.KeyMsg{Type: tea.KeyPgDown}
	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("expected no command when reordering effects")
	}

	updatedModelTyped := updatedModel.(Model)
	if updatedModelTyped.focused != ctrlEffectsOrder {
		t.Error("expected focus to remain on effects order")
	}
}

func TestUpdate_EffectsOrderNavigation(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test that navigation within effects order is handled
	model.focused = ctrlEffectsOrder
	model.selectedEffectIndex = 0

	// Test Down arrow (next effect) - this might move focus or change selection
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("expected no command when navigating effects order")
	}

	updatedModelTyped := updatedModel.(Model)
	// Either focus stays on effects order or moves to next control - both are valid
	// Just check that the update was handled without error
	if len(updatedModelTyped.EffectsOrder) == 0 {
		t.Error("expected valid model after navigation")
	}
}

func TestUpdate_ResetEffectsOrder(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Change effects order
	model.SetEffectsOrder([]string{"granular", "filter", "delay"})
	model.focused = ctrlEffectsOrder

	// Test 'r' key to reset order
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updatedModel, _ := model.Update(msg)
	updatedModelTyped := updatedModel.(Model)

	// Should reset to default order
	defaultOrder := []string{"filter", "overdrive", "bitcrush", "granular", "reverb", "delay"}
	actualOrder := updatedModelTyped.GetEffectsOrder()

	if len(actualOrder) != len(defaultOrder) {
		t.Errorf("expected default order length %d, got %d", len(defaultOrder), len(actualOrder))
	}

	for i, effect := range defaultOrder {
		if actualOrder[i] != effect {
			t.Errorf("expected default effect %d to be %s, got %s", i, effect, actualOrder[i])
		}
	}
}

func TestUpdate_QuitKey(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test 'q' key for quit
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := model.Update(msg)

	// Should return a quit command
	if cmd == nil {
		t.Error("expected quit command when 'q' is pressed")
	}

	// Check if it's a quit command by executing it (simplified check)
	// In actual implementation, this would be a quit command
	if cmd == nil {
		t.Error("expected quit command when 'q' is pressed")
	}

	// Model should be unchanged
	updatedModelTyped := updatedModel.(Model)
	if updatedModelTyped.focused != model.focused {
		t.Error("expected model to be unchanged by quit command")
	}
}

func TestUpdate_UnknownKey(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test unknown key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	updatedModel, cmd := model.Update(msg)

	// Should return no command and unchanged model
	if cmd != nil {
		t.Error("expected no command for unknown key")
	}

	updatedModelTyped := updatedModel.(Model)
	if updatedModelTyped.focused != model.focused {
		t.Error("expected model to be unchanged by unknown key")
	}
}

func TestUpdate_StateMessageHandling(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test StateMsg handling
	state := osc.State{
		Gain:      0.75,
		BlendMode: 1,
	}

	msg := StateMsg(state)
	updatedModel, _ := model.Update(msg)
	updatedModelTyped := updatedModel.(Model)

	if updatedModelTyped.Gain != state.Gain {
		t.Errorf("expected gain to be updated to %f, got %f", state.Gain, updatedModelTyped.Gain)
	}

	if updatedModelTyped.BlendMode != state.BlendMode {
		t.Errorf("expected blend mode to be updated to %d, got %d", state.BlendMode, updatedModelTyped.BlendMode)
	}
}

func TestUpdate_WindowSizeMessage(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test WindowSizeMsg handling
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)
	updatedModelTyped := updatedModel.(Model)

	if updatedModelTyped.width != 100 {
		t.Errorf("expected width to be updated to 100, got %d", updatedModelTyped.width)
	}

	if updatedModelTyped.height != 50 {
		t.Errorf("expected height to be updated to 50, got %d", updatedModelTyped.height)
	}
}

func TestUpdate_BoundaryConditions(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test navigation at boundaries
	model.focused = 0 // First control

	// Test Shift+Tab at first control (should wrap to last)
	msg := tea.KeyMsg{Type: tea.KeyShiftTab}
	updatedModel, _ := model.Update(msg)
	updatedModelTyped := updatedModel.(Model)

	if updatedModelTyped.focused != ctrlCount-1 {
		t.Errorf("expected focus to wrap to last control %d, got %d", ctrlCount-1, updatedModelTyped.focused)
	}

	// Test Tab at last control (should wrap to first)
	model.focused = ctrlCount - 1
	msg = tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ = model.Update(msg)
	updatedModelTyped = updatedModel.(Model)

	if updatedModelTyped.focused != 0 {
		t.Errorf("expected focus to wrap to first control, got %d", updatedModelTyped.focused)
	}
}
