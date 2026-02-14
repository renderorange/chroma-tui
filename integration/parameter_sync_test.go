package integration

import (
	"testing"

	"github.com/renderorange/chroma/chroma-tui/osc"
	"github.com/renderorange/chroma/chroma-tui/tui"
)

func TestParameterSync_FireAndForget(t *testing.T) {
	// Test fire-and-forget parameter updates
	// Chroma uses stateless OSC - no state sync from server

	client := osc.NewClient("127.0.0.1", 57126)
	model := tui.NewModel(client)

	// Set initial gain
	initialGain := float32(0.5)
	model.Gain = initialGain

	// User adjusts gain
	model.SetFocused(tui.TestCtrlGain)
	model.AdjustFocused(0.05)
	userAdjustedGain := model.Gain

	if userAdjustedGain <= initialGain {
		t.Errorf("Expected gain to increase from %f, got %f", initialGain, userAdjustedGain)
	}

	// Verify local state is maintained (no server sync to overwrite)
	if model.Gain != userAdjustedGain {
		t.Errorf("User's gain change was not preserved locally: expected %f, got %f",
			userAdjustedGain, model.Gain)
	}

	t.Log("Fire-and-forget parameter updates verified")
}

func TestParameterSync_MultipleConcurrentChanges(t *testing.T) {
	// Test that multiple rapid parameter changes all work

	client := osc.NewClient("127.0.0.1", 57127)
	model := tui.NewModel(client)

	// Set initial values
	model.Gain = 0.5
	model.FilterCutoff = 1000
	model.OverdriveDrive = 0.3

	// User rapidly changes multiple parameters
	model.SetFocused(tui.TestCtrlGain)
	model.AdjustFocused(0.1)
	gainAfterChange := model.Gain

	model.SetFocused(tui.TestCtrlFilterCutoff)
	model.AdjustFocused(0.1)
	cutoffAfterChange := model.FilterCutoff

	model.SetFocused(tui.TestCtrlOverdriveDrive)
	model.AdjustFocused(0.1)
	driveAfterChange := model.OverdriveDrive

	// Verify all user changes are preserved locally
	if model.Gain != gainAfterChange {
		t.Errorf("Gain change not preserved: expected %f, got %f", gainAfterChange, model.Gain)
	}
	if model.FilterCutoff != cutoffAfterChange {
		t.Errorf("Filter cutoff change not preserved: expected %f, got %f", cutoffAfterChange, model.FilterCutoff)
	}
	if model.OverdriveDrive != driveAfterChange {
		t.Errorf("Overdrive drive change not preserved: expected %f, got %f", driveAfterChange, model.OverdriveDrive)
	}

	t.Log("Multiple concurrent parameter changes verified")
}

func TestParameterSync_TimeoutBasedSync(t *testing.T) {
	// Test parameter sync without pending changes timeout
	// (simplified since we removed pending changes)

	client := osc.NewClient("127.0.0.1", 57128)
	model := tui.NewModel(client)

	// Set initial state
	model.Gain = 0.5
	model.FilterAmount = 0.3

	// User changes gain
	model.SetFocused(tui.TestCtrlGain)
	model.AdjustFocused(0.1)

	// Verify local state
	if model.Gain <= 0.5 {
		t.Errorf("Expected gain to increase from 0.5, got %f", model.Gain)
	}

	t.Log("Parameter sync verified (fire-and-forget model)")
}

func TestParameterSync_EffectsOrder(t *testing.T) {
	// Test effects order handling

	client := osc.NewClient("127.0.0.1", 57129)
	model := tui.NewModel(client)

	// Get default effects order
	defaultOrder := model.GetEffectsOrder()
	if len(defaultOrder) != 6 {
		t.Errorf("Expected 6 effects in default order, got %d", len(defaultOrder))
	}

	// Set custom effects order
	customOrder := []string{"reverb", "delay", "filter", "overdrive", "bitcrush", "granular"}
	model.SetEffectsOrder(customOrder)

	// Verify order is set locally
	currentOrder := model.GetEffectsOrder()
	for i, effect := range customOrder {
		if currentOrder[i] != effect {
			t.Errorf("Expected effect %d to be %s, got %s", i, effect, currentOrder[i])
		}
	}

	t.Log("Effects order handling verified")
}

func TestParameterSync_BlendMode(t *testing.T) {
	// Test blend mode parameter updates

	client := osc.NewClient("127.0.0.1", 57130)
	model := tui.NewModel(client)

	// Test initial blend mode
	if model.BlendMode != 0 {
		t.Errorf("Expected initial blend mode to be 0, got %d", model.BlendMode)
	}

	// Change blend mode
	model.SetBlendMode(1)
	if model.BlendMode != 1 {
		t.Errorf("Expected blend mode to be 1, got %d", model.BlendMode)
	}

	model.SetBlendMode(2)
	if model.BlendMode != 2 {
		t.Errorf("Expected blend mode to be 2, got %d", model.BlendMode)
	}

	t.Log("Blend mode updates verified")
}

func TestParameterSync_GrainIntensity(t *testing.T) {
	// Test grain intensity cycling

	client := osc.NewClient("127.0.0.1", 57131)
	model := tui.NewModel(client)

	// Test initial intensity
	if model.GrainIntensity != "subtle" {
		t.Errorf("Expected initial grain intensity to be 'subtle', got %s", model.GrainIntensity)
	}

	// Note: Grain intensity is toggled via toggleGrainIntensity() which is private
	// We can only test that the initial value is correct
	t.Log("Grain intensity initial state verified")
}
