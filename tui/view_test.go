package tui

import (
	"fmt"
	"testing"

	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestView_RenderingWithDefaultModel(t *testing.T) {
	// Create a model with default values
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test that View() returns a non-empty string
	view := model.View()
	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}

	// Test that view contains expected section titles
	expectedSections := []string{
		"INPUT", "FILTER", "OVERDRIVE", "BITCRUSH",
		"GRANULAR", "REVERB", "DELAY", "EFFECTS ORDER", "GLOBAL",
	}

	for _, section := range expectedSections {
		if !contains(view, section) {
			t.Errorf("expected view to contain section '%s'", section)
		}
	}
}

func TestView_RenderingWithFocusedControl(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test that view renders correctly with different focused controls
	testCases := []struct {
		name  string
		focus control
	}{
		{"Gain focused", ctrlGain},
		{"Filter focused", ctrlFilterAmount},
		{"Overdrive focused", ctrlOverdriveDrive},
		{"Granular focused", ctrlGranularDensity},
		{"Reverb focused", ctrlReverbDecayTime},
		{"Delay focused", ctrlDelayTime},
		{"Effects Order focused", ctrlEffectsOrder},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set focus
			model.focused = tc.focus
			view := model.View()

			// View should be non-empty
			if len(view) == 0 {
				t.Errorf("expected non-empty view when %s is focused", tc.name)
			}

			// View should contain all section titles
			expectedSections := []string{
				"INPUT", "FILTER", "OVERDRIVE", "BITCRUSH",
				"GRANULAR", "REVERB", "DELAY", "EFFECTS ORDER", "GLOBAL",
			}

			for _, section := range expectedSections {
				if !contains(view, section) {
					t.Errorf("expected view to contain section '%s' when %s is focused", section, tc.name)
				}
			}
		})
	}
}

func TestView_RenderingWithDifferentWidths(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test with different widths
	testWidths := []int{40, 80, 120, 200}

	for _, width := range testWidths {
		t.Run(fmt.Sprintf("width_%d", width), func(t *testing.T) {
			model.width = width
			view := model.View()

			if len(view) == 0 {
				t.Errorf("expected View() to return non-empty string for width %d", width)
			}

			// View should adapt to width (no hard line breaks at inappropriate places)
			if width >= 80 && len(view) < 100 {
				t.Errorf("expected longer view for width %d", width)
			}
		})
	}
}

func TestView_RenderingWithMIDIPort(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Set MIDI port
	model.midiPort = "USB MIDI Device"
	view := model.View()

	// Should contain MIDI port information
	if !contains(view, "USB MIDI Device") {
		t.Error("expected view to contain MIDI port information when MIDI port is set")
	}
}

func TestView_RenderingWithEffectsOrder(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test with default effects order
	view := model.View()
	expectedEffects := []string{"filter", "overdrive", "bitcrush", "granular", "reverb", "delay"}

	for _, effect := range expectedEffects {
		if !contains(view, effect) {
			t.Errorf("expected view to contain effect '%s' in effects order", effect)
		}
	}

	// Test with custom effects order
	model.SetEffectsOrder([]string{"granular", "filter", "delay"})
	view = model.View()

	// Should contain the reordered effects
	if !contains(view, "granular") || !contains(view, "filter") || !contains(view, "delay") {
		t.Error("expected view to contain reordered effects")
	}
}

func TestView_RenderingWithDifferentBlendModes(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test that blend mode changes don't break view rendering
	model.BlendMode = 0 // mirror
	view := model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view with blend mode 0")
	}

	model.BlendMode = 1 // complement
	view = model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view with blend mode 1")
	}

	model.BlendMode = 2 // transform
	view = model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view with blend mode 2")
	}
}

func TestView_RenderingWithDifferentGrainIntensities(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test that grain intensity changes don't break view rendering
	testIntensities := []string{"subtle", "pronounced", "extreme"}

	for _, intensity := range testIntensities {
		model.GrainIntensity = intensity
		view := model.View()
		if len(view) == 0 {
			t.Errorf("expected non-empty view with grain intensity '%s'", intensity)
		}
	}
}

func TestView_RenderingWithEnabledDisabledEffects(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test with all effects enabled
	model.FilterEnabled = true
	model.OverdriveEnabled = true
	model.GranularEnabled = true
	model.BitcrushEnabled = true
	model.ReverbEnabled = true
	model.DelayEnabled = true

	view := model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view with effects enabled")
	}

	// Test with all effects disabled
	model.FilterEnabled = false
	model.OverdriveEnabled = false
	model.GranularEnabled = false
	model.BitcrushEnabled = false
	model.ReverbEnabled = false
	model.DelayEnabled = false

	view = model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view with effects disabled")
	}
}

func TestView_RenderingWithFrozenControls(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test with frozen controls
	model.InputFrozen = true
	model.GranularFrozen = true

	view := model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view with frozen controls")
	}
}

func TestView_RenderingWithParameterValues(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Set specific parameter values
	model.Gain = 0.75
	model.FilterCutoff = 1500
	model.OverdriveDrive = 0.5
	model.GranularDensity = 0.25
	model.ReverbDecayTime = 2.5
	model.DelayTime = 0.375

	view := model.View()

	// Should contain parameter values (this tests the formatting)
	if len(view) == 0 {
		t.Error("expected view to contain parameter values")
	}

	// Test that view changes when parameters change
	originalView := view
	model.Gain = 0.25
	newView := model.View()

	if newView == originalView {
		t.Error("expected view to change when parameter values change")
	}
}

func TestView_RenderingWithZeroWidth(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test with zero width (should default to 80)
	model.width = 0
	view := model.View()

	if len(view) == 0 {
		t.Error("expected View() to handle zero width gracefully")
	}
}

func TestView_RenderingWithNegativeWidth(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	// Test with negative width (should default to 80)
	model.width = -10
	view := model.View()

	if len(view) == 0 {
		t.Error("expected View() to handle negative width gracefully")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
