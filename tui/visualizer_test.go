package tui

import (
	"math"
	"strings"
	"testing"

	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestRenderMultiLineSpectrum(t *testing.T) {
	spectrum := [8]float32{1.0, 0.8, 0.6, 0.4, 0.2, 0.4, 0.6, 0.8}
	result := renderMultiLineSpectrum(spectrum, 64, 6)

	// Should have 6 lines (height parameter)
	lines := strings.Split(result, "\n")
	if len(lines) != 6 {
		t.Errorf("Expected 6 lines, got %d", len(lines))
	}

	// First line should be tallest (top of bars)
	if !strings.Contains(lines[0], "█") {
		t.Errorf("First line should contain █, got: %s", lines[0])
	}

	// Last line should be shortest
	if !strings.Contains(lines[5], "▁") {
		t.Errorf("Last line should contain ▁, got: %s", lines[5])
	}
}

func TestRenderMultiLineSpectrumBoundaryValues(t *testing.T) {
	tests := []struct {
		name     string
		spectrum [8]float32
		width    int
		height   int
		wantErr  bool
	}{
		{
			name:     "all zeros",
			spectrum: [8]float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
			width:    64,
			height:   4,
		},
		{
			name:     "all ones",
			spectrum: [8]float32{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
			width:    64,
			height:   4,
		},
		{
			name:     "mixed values",
			spectrum: [8]float32{0.0, 0.25, 0.5, 0.75, 1.0, 0.75, 0.5, 0.25},
			width:    64,
			height:   4,
		},
		{
			name:     "threshold edge case",
			spectrum: [8]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			width:    64,
			height:   2, // threshold will be 0.5
		},
		{
			name:     "too narrow",
			spectrum: [8]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			width:    15, // below minVisualizerWidth
			height:   4,
		},
		{
			name:     "zero height",
			spectrum: [8]float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			width:    64,
			height:   0,
		},
		{
			name:     "negative values",
			spectrum: [8]float32{-0.5, 0.0, 0.5, 1.0, 1.5, 0.5, 0.0, -0.5},
			width:    64,
			height:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderMultiLineSpectrum(tt.spectrum, tt.width, tt.height)

			// Check edge cases return empty string
			if tt.width < minVisualizerWidth || tt.height < 1 {
				if result != "" {
					t.Errorf("Expected empty string for invalid dimensions, got: %s", result)
				}
				return
			}

			// Valid cases should produce expected output
			lines := strings.Split(result, "\n")
			if len(lines) != tt.height {
				t.Errorf("Expected %d lines, got %d", tt.height, len(lines))
			}

			// Each line should have exactly 8 bands
			for i, line := range lines {
				if line == "" {
					t.Errorf("Line %d should not be empty", i)
				}
			}
		})
	}
}

func TestRenderMultiLineSpectrumOutputPatterns(t *testing.T) {
	// Test specific output patterns for known inputs
	spectrum := [8]float32{1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0}
	result := renderMultiLineSpectrum(spectrum, 16, 2) // Minimum width for 8 bands

	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(lines))
	}

	// Top line (line 1) should have full bars for 1.0 values
	// Since threshold is 0.5, 1.0 >= 0.5 should show some character
	topLine := lines[0]
	if len(topLine) == 0 {
		t.Error("Top line should not be empty")
	}

	// Bottom line (line 0) should have full bars for 1.0 values
	// Since threshold is 0.0, all 1.0 >= 0.0 should show full blocks
	bottomLine := lines[1]
	if len(bottomLine) == 0 {
		t.Error("Bottom line should not be empty")
	}

	// Test threshold calculation fix - when threshold approaches 1.0
	spectrumHigh := [8]float32{0.99, 0.99, 0.99, 0.99, 0.99, 0.99, 0.99, 0.99}
	resultHigh := renderMultiLineSpectrum(spectrumHigh, 16, 1)

	// Should not panic and should produce valid output
	if resultHigh == "" {
		t.Error("High threshold case should not return empty string")
	}
}

func TestRenderWaveform(t *testing.T) {
	waveform := [64]float32{}
	// Create simple sine wave pattern
	for i := range waveform {
		waveform[i] = float32(math.Sin(float64(i) * 0.2))
	}

	result := renderWaveform(waveform, 64)

	// Should contain waveform characters
	if !strings.Contains(result, "─") {
		t.Error("Result should contain '─' character")
	}
	if !strings.Contains(result, "╭") {
		t.Error("Result should contain '╭' character")
	}
	if !strings.Contains(result, "╮") {
		t.Error("Result should contain '╮' character")
	}
	if result == "" {
		t.Error("Result should not be empty")
	}
}

func TestRenderWaveformEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		waveform [64]float32
		width    int
		wantErr  bool
	}{
		{
			name:     "empty waveform",
			waveform: [64]float32{},
			width:    64,
		},
		{
			name:     "all zeros",
			waveform: [64]float32{},
			width:    64,
		},
		{
			name: "extreme positive values",
			waveform: func() [64]float32 {
				var w [64]float32
				for i := range w {
					w[i] = 2.0 // Above normal range
				}
				return w
			}(),
			width: 64,
		},
		{
			name: "extreme negative values",
			waveform: func() [64]float32 {
				var w [64]float32
				for i := range w {
					w[i] = -2.0 // Below normal range
				}
				return w
			}(),
			width: 64,
		},
		{
			name: "alternating extremes",
			waveform: func() [64]float32 {
				var w [64]float32
				for i := range w {
					if i%2 == 0 {
						w[i] = 1.0
					} else {
						w[i] = -1.0
					}
				}
				return w
			}(),
			width: 64,
		},
		{
			name:     "minimum width",
			waveform: [64]float32{},
			width:    16, // minVisualizerWidth
		},
		{
			name:     "too narrow",
			waveform: [64]float32{},
			width:    15, // below minVisualizerWidth
			wantErr:  true,
		},
		{
			name: "large width with step > 1",
			waveform: func() [64]float32 {
				var w [64]float32
				for i := range w {
					w[i] = float32(math.Sin(float64(i) * 0.5))
				}
				return w
			}(),
			width: 8, // Will cause step = 8, testing continuity fix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the waveform for tests that use empty arrays
			if tt.name == "empty waveform" || tt.name == "all zeros" {
				// Already zeros by default
			} else if tt.name == "minimum width" {
				// Create a simple pattern
				for i := range tt.waveform {
					tt.waveform[i] = float32(math.Sin(float64(i) * 0.1))
				}
			}

			result := renderWaveform(tt.waveform, tt.width)

			// Check edge cases return empty string
			if tt.width < minVisualizerWidth {
				if result != "" {
					t.Errorf("Expected empty string for width %d, got: %s", tt.width, result)
				}
				return
			}

			// Valid cases should produce expected output
			if result == "" {
				t.Error("Result should not be empty for valid inputs")
			}

			// Should contain some waveform characters
			if tt.width >= minVisualizerWidth {
				hasWaveformChars := strings.Contains(result, "─") ||
					strings.Contains(result, "╭") ||
					strings.Contains(result, "╮") ||
					strings.Contains(result, "╯") ||
					strings.Contains(result, "╰") ||
					strings.Contains(result, "│")
				if !hasWaveformChars {
					t.Errorf("Result should contain waveform characters, got: %s", result)
				}
			}
		})
	}
}

func TestRenderWaveformYMapping(t *testing.T) {
	// Test that Y mapping correctly handles the full -1..1 range
	waveform := [64]float32{}

	// Test extreme values
	for i := range waveform {
		if i < 16 {
			waveform[i] = -1.0 // Bottom
		} else if i < 32 {
			waveform[i] = 0.0 // Middle
		} else if i < 48 {
			waveform[i] = 1.0 // Top
		} else {
			waveform[i] = -0.5 // Lower middle
		}
	}

	result := renderWaveform(waveform, 64)

	// Should not be empty
	if result == "" {
		t.Error("Result should not be empty with extreme values")
	}

	// Should contain direction changes due to the extreme value transitions
	hasCorners := strings.Contains(result, "╭") ||
		strings.Contains(result, "╮") ||
		strings.Contains(result, "╯") ||
		strings.Contains(result, "╰")
	if !hasCorners {
		t.Errorf("Result should contain corner characters for extreme value transitions, got: %s", result)
	}
}

func TestRenderVisualizerCombination(t *testing.T) {
	model := Model{
		width:    80,
		Spectrum: [8]float32{0.8, 0.6, 0.4, 0.2, 0.4, 0.6, 0.8, 1.0},
		Waveform: [64]float32{},
	}

	result := model.renderVisualizer(80)

	// Should contain both spectrum and waveform sections
	if !strings.Contains(result, "SPECTRUM") {
		t.Error("Result should contain SPECTRUM section")
	}
	if !strings.Contains(result, "WAVEFORM") {
		t.Error("Result should contain WAVEFORM section")
	}
	if !strings.Contains(result, "30Hz") {
		t.Error("Result should contain frequency labels")
	}
}

func TestEffectsOrderRendering(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.EffectsOrder = []string{"filter", "granular", "delay"}
	model.focused = ctrlEffectsOrder

	view := model.View()
	if !strings.Contains(view, "EFFECTS ORDER") {
		t.Error("View should contain EFFECTS ORDER section")
	}
	if !strings.Contains(view, "filter") {
		t.Error("View should contain filter effect")
	}
	if !strings.Contains(view, "granular") {
		t.Error("View should contain granular effect")
	}
	if !strings.Contains(view, "delay") {
		t.Error("View should contain delay effect")
	}
}
