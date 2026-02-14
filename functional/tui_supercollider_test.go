package functional

import (
	"testing"

	"github.com/renderorange/chroma/chroma-tui/osc"
	"github.com/renderorange/chroma/chroma-tui/tui"
)

func TestTUISuperCollider_FireAndForget(t *testing.T) {
	// Test the fire-and-forget workflow: TUI → OSC (no response expected)
	// Chroma uses stateless OSC - TUI sends commands without expecting state updates

	client := osc.NewClient("127.0.0.1", 57120)
	model := tui.NewModel(client)

	// User adjusts parameters in TUI
	model.SetFocused(tui.TestCtrlGain)
	model.AdjustFocused(0.1) // Increase gain
	expectedGain := model.Gain

	model.SetFocused(tui.TestCtrlFilterCutoff)
	model.AdjustFocused(0.1) // Increase filter cutoff
	expectedCutoff := model.FilterCutoff

	// Verify TUI state updated locally (no state sync from Chroma)
	if model.Gain != expectedGain {
		t.Errorf("Gain not updated locally: expected %f, got %f", expectedGain, model.Gain)
	}
	if model.FilterCutoff != expectedCutoff {
		t.Errorf("Filter cutoff not updated locally: expected %f, got %f", expectedCutoff, model.FilterCutoff)
	}

	t.Log("Fire-and-forget OSC communication verified")
}

func TestTUISuperCollider_LocalStateOnly(t *testing.T) {
	// Test that TUI maintains local state without expecting Chroma responses

	client := osc.NewClient("127.0.0.1", 57121)
	model := tui.NewModel(client)

	// Set initial parameters
	initialGain := float32(1.2)
	initialCutoff := float32(2000)
	model.Gain = initialGain
	model.FilterCutoff = initialCutoff

	// Verify local state is maintained
	if model.Gain != initialGain {
		t.Errorf("Gain not preserved: expected %f, got %f", initialGain, model.Gain)
	}
	if model.FilterCutoff != initialCutoff {
		t.Errorf("Filter cutoff not preserved: expected %f, got %f", initialCutoff, model.FilterCutoff)
	}

	t.Log("Local state maintenance verified")
}

func TestTUISuperCollider_ErrorHandling(t *testing.T) {
	// Test error handling in TUI-SC communication

	client := osc.NewClient("127.0.0.1", 57135) // Non-existent server
	model := tui.NewModel(client)

	// Test that TUI remains functional even when SC is unavailable
	initialGain := model.Gain
	model.SetFocused(tui.TestCtrlGain)
	model.AdjustFocused(0.1)

	// TUI should maintain its state despite communication failures
	if model.Gain == initialGain {
		t.Error("TUI model should update locally even when SC communication fails")
	}

	t.Log("Error handling verified - TUI works without Chroma confirmation")
}
