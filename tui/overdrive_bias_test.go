package tui

import (
	"github.com/renderorange/chroma/chroma-tui/osc"
	"testing"
)

func TestOverdriveBiasImplementation(t *testing.T) {
	// Test overdrive bias implementation
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test default value
	if model.OverdriveBias != 0.5 {
		t.Errorf("Expected default OverdriveBias to be 0.5, got %.2f", model.OverdriveBias)
	}

	// Test parameter adjustment
	model.SetFocused(TestCtrlOverdriveBias)
	model.AdjustFocused(0.5) // Increase bias
	expected := float32(1.0)
	if model.OverdriveBias != expected {
		t.Errorf("After +0.5 adjustment expected %.2f, got %.2f", expected, model.OverdriveBias)
	}

	model.AdjustFocused(-1.5) // Decrease bias beyond range
	expected = float32(-0.5)  // Should be clamped at -1.0
	if model.OverdriveBias != expected {
		t.Errorf("After -1.5 adjustment expected clamped %.2f, got %.2f", expected, model.OverdriveBias)
	}

	// Test boundary conditions
	model.OverdriveBias = 0.0
	model.AdjustFocused(2.0) // Should clamp to 1.0
	if model.OverdriveBias != 1.0 {
		t.Errorf("After exceeding max boundary, expected 1.0, got %.2f", model.OverdriveBias)
	}

	model.OverdriveBias = 0.0
	model.AdjustFocused(-2.0) // Should clamp to -1.0
	if model.OverdriveBias != -1.0 {
		t.Errorf("After exceeding min boundary, expected -1.0, got %.2f", model.OverdriveBias)
	}
}
