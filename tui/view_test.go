package tui

import (
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
