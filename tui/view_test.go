package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestView_RenderingWithDefaultModel(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	view := model.View()
	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}

	found := false
	for _, c := range view {
		if c == 'E' || c == 'e' {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected view to contain some content")
	}
}

func TestView_RenderingWithFocusedControl(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	view := model.View()
	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}
}

func TestView_RenderingWithDifferentWidths(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)

	testWidths := []int{80, 120, 200}
	for _, width := range testWidths {
		model := NewModel(client)
		model.InitLists(width, 40)

		view := model.View()
		if len(view) == 0 {
			t.Errorf("expected non-empty view for width %d", width)
		}
	}
}

func TestView_RenderingWithMIDIPort(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.SetMidiPort("MIDI In")
	model.InitLists(80, 40)

	view := model.View()
	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}
}

func TestView_RenderingWithEffectsOrder(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	view := model.View()
	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}
}

func TestView_RenderingWithParameterValues(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.Gain = 0.75
	model.FilterCutoff = 4000
	model.InitLists(80, 40)

	view := model.View()
	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}
}

func TestView_ContainsBothPanels(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(100, 40)

	view := model.View()

	// View should contain effects list title
	if !strings.Contains(view, "Effects") {
		t.Error("expected view to contain 'Effects' title")
	}

	// View should contain parameter content (Input is default section)
	if !strings.Contains(view, "Gain") {
		t.Error("expected view to contain 'Gain' parameter")
	}
}

func TestView_FooterRendering(t *testing.T) {
	tests := []struct {
		name            string
		navigationMode  int
		currentSection  string
		expectedContent string
	}{
		{
			name:            "effects list mode",
			navigationMode:  0, // modeEffectsList
			currentSection:  "input",
			expectedContent: "enter: open params",
		},
		{
			name:            "parameter list mode normal section",
			navigationMode:  1, // modeParameterList
			currentSection:  "filter",
			expectedContent: "h/l: adjust value",
		},
		{
			name:            "parameter list mode global section",
			navigationMode:  1, // modeParameterList
			currentSection:  "global",
			expectedContent: "pgup/pgdn: reorder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := osc.NewClient("127.0.0.1", 57120)
			model := NewModel(client)
			model.InitLists(80, 40)
			model.SetNavigationMode(tt.navigationMode)
			model.SetCurrentSection(tt.currentSection)

			footer := model.renderFooter(76) // 80 - 4 for app padding

			if !strings.Contains(footer, tt.expectedContent) {
				t.Errorf("expected footer to contain '%s', got: %s", tt.expectedContent, footer)
			}
		})
	}
}

func TestView_StatusBarRendering(t *testing.T) {
	tests := []struct {
		name            string
		connected       bool
		midiPort        string
		expectedContent []string
	}{
		{
			name:            "connected with midi",
			connected:       true,
			midiPort:        "USB MIDI Device",
			expectedContent: []string{"Connected", "MIDI: USB MIDI Device"},
		},
		{
			name:            "disconnected with no midi",
			connected:       false,
			midiPort:        "",
			expectedContent: []string{"Disconnected", "No MIDI"},
		},
		{
			name:            "connected with no midi",
			connected:       true,
			midiPort:        "",
			expectedContent: []string{"Connected", "No MIDI"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := osc.NewClient("127.0.0.1", 57120)
			model := NewModel(client)
			model.InitLists(80, 40)
			model.SetConnected(tt.connected)
			model.SetMidiPort(tt.midiPort)

			statusBar := model.renderStatusBar(76) // 80 - 4 for app padding

			for _, expected := range tt.expectedContent {
				if !strings.Contains(statusBar, expected) {
					t.Errorf("expected status bar to contain '%s', got: %s", expected, statusBar)
				}
			}
		})
	}
}

func TestView_VerticalStackingAndAlignment(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57122)
	model := NewModel(client)
	model.InitLists(80, 40)

	view := model.View()

	// Verify non-empty
	if len(view) == 0 {
		t.Error("expected View() to return non-empty string")
	}

	// Verify footer present
	if !strings.Contains(view, "j/k: navigate") {
		t.Error("expected view to contain footer with navigation hints")
	}

	// Verify status bar present
	if !strings.Contains(view, "Connected") && !strings.Contains(view, "Disconnected") {
		t.Error("expected view to contain status bar with connection info")
	}
}

func TestView_MinimumTerminalSize(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57123)
	model := NewModel(client)

	// Simulate window size message to set width and height
	msg := tea.WindowSizeMsg{Width: 50, Height: 15}
	model.Update(msg)

	view := model.View()

	if !strings.Contains(view, "Terminal too small") {
		t.Errorf("expected warning for small terminal, got: %s", view)
	}
}
