package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestLists_effectItem(t *testing.T) {
	item := effectItem{
		id:          "test",
		title:       "Test Title",
		description: "Test Description",
		enabled:     true,
	}

	if item.Title() != "Test Title" {
		t.Errorf("expected Title() to return 'Test Title', got %s", item.Title())
	}
	if item.Description() != "Test Description" {
		t.Errorf("expected Description() to return 'Test Description', got %s", item.Description())
	}
	if item.FilterValue() != "Test Title" {
		t.Errorf("expected FilterValue() to return 'Test Title', got %s", item.FilterValue())
	}
}

func TestLists_parameterItem(t *testing.T) {
	item := parameterItem{
		id:                "test",
		title:             "Test Title",
		customDescription: "Test Description",
		ctrl:              ctrlGain,
		isToggle:          true,
		isActive:          false,
	}

	if item.Title() != "Test Title" {
		t.Errorf("expected Title() to return 'Test Title', got %s", item.Title())
	}
	if item.Description() != "Test Description" {
		t.Errorf("expected Description() to return 'Test Description', got %s", item.Description())
	}
	if item.FilterValue() != "Test Title" {
		t.Errorf("expected FilterValue() to return 'Test Title', got %s", item.FilterValue())
	}
}

func TestLists_newEffectItem(t *testing.T) {
	item := newEffectItem("filter", "Filter", true)
	if item.id != "filter" {
		t.Errorf("expected id 'filter', got %s", item.id)
	}
	if item.title != "Filter" {
		t.Errorf("expected title 'Filter', got %s", item.title)
	}
	if !item.enabled {
		t.Error("expected enabled to be true")
	}

	itemDisabled := newEffectItem("delay", "Delay", false)
	if itemDisabled.enabled {
		t.Error("expected enabled to be false")
	}
}

func TestLists_newParameterItem(t *testing.T) {
	item := newParameterItem("gain", "Gain", 0.5, 0, 2, ctrlGain, false, false, 10)
	if item.id != "gain" {
		t.Errorf("expected id 'gain', got %s", item.id)
	}
	if item.ctrl != ctrlGain {
		t.Errorf("expected ctrl ctrlGain, got %v", item.ctrl)
	}
	if item.isToggle {
		t.Error("expected isToggle to be false")
	}

	toggleItem := newParameterItem("enabled", "Enabled", 0, 0, 0, ctrlFilterEnabled, true, true, 10)
	if !toggleItem.isToggle {
		t.Error("expected isToggle to be true")
	}
	if !toggleItem.isActive {
		t.Error("expected isActive to be true")
	}
}

func TestLists_buildEffectsList(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	items := model.buildEffectsList()
	if len(items) != 8 {
		t.Errorf("expected 8 effect items, got %d", len(items))
	}

	found := false
	for _, item := range items {
		if eff, ok := item.(effectItem); ok {
			if eff.id == "filter" {
				found = true
				if !eff.enabled {
					t.Error("expected filter to be enabled by default")
				}
			}
		}
	}
	if !found {
		t.Error("expected to find filter item")
	}
}

func TestLists_buildParameterList(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	tests := []struct {
		section   string
		wantCount int
	}{
		{"input", 3},
		{"filter", 4},
		{"overdrive", 5},
		{"bitcrush", 5},
		{"granular", 8},
		{"reverb", 3},
		{"delay", 6},
		{"global", 3},
	}

	for _, tt := range tests {
		t.Run(tt.section, func(t *testing.T) {
			items := model.buildParameterList(tt.section)
			if len(items) != tt.wantCount {
				t.Errorf("expected %d items for section %s, got %d", tt.wantCount, tt.section, len(items))
			}
		})
	}
}

func TestLists_formatSliderValue(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	result := model.formatSliderValue("Gain", 0.5, 0, 2)
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
}

func TestLists_formatEffectsOrder(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.EffectsOrder = []string{"filter", "overdrive", "granular"}

	result := model.formatEffectsOrder()
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestModel_NextControl(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	initial := model.focused
	model.NextControl()
	if model.focused != initial+1 {
		t.Errorf("expected focus to move to %d, got %d", initial+1, model.focused)
	}
}

func TestModel_PrevControl(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	model.focused = 5
	initial := model.focused
	model.PrevControl()
	if model.focused != initial-1 {
		t.Errorf("expected focus to move to %d, got %d", initial-1, model.focused)
	}
}

func TestModel_Focused(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.focused = ctrlGain

	if model.Focused() != ctrlGain {
		t.Errorf("expected Focused() to return ctrlGain, got %v", model.Focused())
	}
}

func TestModel_IsConnected(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	if model.IsConnected() != false {
		t.Error("expected IsConnected() to return false initially")
	}
}

func TestModel_SetEffectsOrder(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	newOrder := []string{"reverb", "delay", "filter"}
	model.SetEffectsOrder(newOrder)

	if len(model.EffectsOrder) != 3 {
		t.Errorf("expected EffectsOrder to have 3 items, got %d", len(model.EffectsOrder))
	}
}

func TestModel_GetEffectsOrder(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	order := model.GetEffectsOrder()
	if len(order) != 6 {
		t.Errorf("expected 6 effects by default, got %d", len(order))
	}
}

func TestUpdate_Init(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)

	cmd := model.Init()
	if cmd != nil {
		t.Error("expected Init() to return nil cmd")
	}
}

func TestUpdate_handleEnterKey_parameterMode(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Enter parameter mode
	model.navigationMode = modeParameterList
	model.currentSection = "filter"

	// Select a toggle item and press enter
	model.parameterList.SetItems([]list.Item{
		newParameterItem("enabled", "Filter", 0, 0, 0, ctrlFilterEnabled, true, false, 10),
	})
	model.parameterList.Select(0)

	updatedModel, _ := model.handleEnterKey()
	_ = updatedModel
	// Note: the returned model may be different from the original
}

func TestView_Loading(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.width = 80 // Simulate window size received

	view := model.View()
	// When width is set but lists aren't initialized, view will show list items (empty)
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestUpdate_toggleByControl(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	model.FilterEnabled = false
	model.toggleByControl(ctrlFilterEnabled)

	if !model.FilterEnabled {
		t.Error("expected FilterEnabled to be true after toggleByControl")
	}
}

func TestUpdate_toggleGrainIntensity(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	initial := model.GrainIntensity
	model.toggleGrainIntensity()

	if model.GrainIntensity == initial {
		t.Error("expected GrainIntensity to change")
	}
}

func TestUpdate_adjustLogarithmic(t *testing.T) {
	tests := []struct {
		name        string
		current     float32
		delta       float32
		min         float32
		max         float32
		wantInRange bool
	}{
		{"normal", 0.5, 0.1, 0.01, 0.5, true},
		{"at_min", 0.01, -0.1, 0.01, 0.5, true},
		{"at_max", 0.5, 0.1, 0.01, 0.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adjustLogarithmic(tt.current, tt.delta, tt.min, tt.max)
			if tt.wantInRange && (result < tt.min || result > tt.max) {
				t.Errorf("expected result in range [%f, %f], got %f", tt.min, tt.max, result)
			}
		})
	}
}

func TestUpdate_setBlendMode(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	model.setBlendMode(1)
	if model.BlendMode != 1 {
		t.Errorf("expected BlendMode=1, got %d", model.BlendMode)
	}

	model.setBlendMode(2)
	if model.BlendMode != 2 {
		t.Errorf("expected BlendMode=2, got %d", model.BlendMode)
	}
}

func TestUpdate_toggleGrainIntensity_all(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	model.GrainIntensity = "subtle"
	model.toggleGrainIntensity()
	if model.GrainIntensity != "pronounced" {
		t.Errorf("expected 'pronounced', got %s", model.GrainIntensity)
	}

	model.toggleGrainIntensity()
	if model.GrainIntensity != "extreme" {
		t.Errorf("expected 'extreme', got %s", model.GrainIntensity)
	}

	model.toggleGrainIntensity()
	if model.GrainIntensity != "subtle" {
		t.Errorf("expected 'subtle', got %s", model.GrainIntensity)
	}
}

func TestHelper_ToggleFocused(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	m := NewModel(client)
	m.InitLists(80, 40)

	m.SetFocused(ctrlFilterEnabled)
	if m.focused != ctrlFilterEnabled {
		t.Errorf("expected focused to be ctrlFilterEnabled")
	}

	m.ToggleFocused()
	// Toggle should work without panic
}

func TestHelper_SetBlendMode(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	m := NewModel(client)
	m.InitLists(80, 40)

	m.SetBlendMode(2)
	if m.BlendMode != 2 {
		t.Errorf("expected BlendMode=2, got %d", m.BlendMode)
	}
}

func TestView_ModeEffectsList(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	view := model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view")
	}
}

func TestView_ModeParameterList(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)
	model.navigationMode = modeParameterList
	model.currentSection = "filter"
	model.refreshParameterList()

	view := model.View()
	if len(view) == 0 {
		t.Error("expected non-empty view in parameter mode")
	}
}

func TestUpdate_ParameterListNavigation(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	// Enter parameter mode
	model.navigationMode = modeParameterList
	model.currentSection = "filter"
	model.refreshParameterList()

	// Parameter list should start at index 0
	if model.parameterList.Index() != 0 {
		t.Fatalf("expected initial parameter index 0, got %d", model.parameterList.Index())
	}

	// Press down arrow to navigate parameters
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*Model)

	if m.parameterList.Index() != 1 {
		t.Errorf("expected parameter list index 1 after down arrow, got %d", m.parameterList.Index())
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	model := NewModel(client)
	model.InitLists(80, 40)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)

	m := updatedModel.(*Model)
	if m.width != 100 || m.height != 50 {
		t.Errorf("expected width=100, height=50, got width=%d, height=%d", m.width, m.height)
	}
}

func TestFormatSliderValue_DynamicWidth(t *testing.T) {
	tests := []struct {
		name        string
		value       float32
		min         float32
		max         float32
		sliderWidth int
		wantFilled  int
	}{
		{
			name:        "10 char slider at 50%",
			value:       0.5,
			min:         0,
			max:         1,
			sliderWidth: 10,
			wantFilled:  5,
		},
		{
			name:        "20 char slider at 50%",
			value:       0.5,
			min:         0,
			max:         1,
			sliderWidth: 20,
			wantFilled:  10,
		},
		{
			name:        "30 char slider at 75%",
			value:       0.75,
			min:         0,
			max:         1,
			sliderWidth: 30,
			wantFilled:  22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSliderValue(tt.value, tt.min, tt.max, tt.sliderWidth)

			filled := strings.Count(result, "=")
			if filled != tt.wantFilled {
				t.Errorf("expected %d filled chars, got %d in: %s", tt.wantFilled, filled, result)
			}

			bar := strings.Split(result, " ")[0]
			if len(bar) != tt.sliderWidth {
				t.Errorf("expected slider width %d, got %d", tt.sliderWidth, len(bar))
			}
		})
	}
}
