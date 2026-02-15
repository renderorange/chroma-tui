package tui

import (
	"strings"
	"testing"

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
