package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type effectItem struct {
	id          string
	title       string
	description string
	enabled     bool
}

func (i effectItem) Title() string       { return i.title }
func (i effectItem) Description() string { return i.description }
func (i effectItem) FilterValue() string { return i.title }

type parameterItem struct {
	id                string
	title             string
	value             float32
	min               float32
	max               float32
	ctrl              control
	isToggle          bool
	isActive          bool
	sliderWidth       int
	customDescription string // For special items like intensity, blendMode, effectsOrder
}

func (i parameterItem) Title() string { return i.title }
func (i parameterItem) Description() string {
	if i.customDescription != "" {
		return i.customDescription
	}
	if i.isToggle {
		return ""
	}
	return formatSliderValue(i.value, i.min, i.max, i.sliderWidth)
}
func (i parameterItem) FilterValue() string { return i.title }

func newEffectItem(id, title string, enabled bool, showStatus bool) effectItem {
	description := ""
	if showStatus {
		if enabled {
			description = "[ENABLED]"
		} else {
			description = "[DISABLED]"
		}
	}
	return effectItem{
		id:          id,
		title:       title,
		description: description,
		enabled:     enabled,
	}
}

func newParameterItem(id, title string, value, min, max float32, ctrl control, isToggle, isActive bool, sliderWidth int) parameterItem {
	prefix := ""
	if isToggle {
		if isActive {
			prefix = "[X] "
		} else {
			prefix = "[ ] "
		}
	}
	return parameterItem{
		id:          id,
		title:       prefix + title,
		value:       value,
		min:         min,
		max:         max,
		ctrl:        ctrl,
		isToggle:    isToggle,
		isActive:    isActive,
		sliderWidth: sliderWidth,
	}
}

func newParameterItemWithDesc(id, title string, customDesc string, ctrl control, sliderWidth int) parameterItem {
	return parameterItem{
		id:                id,
		title:             title,
		customDescription: customDesc,
		ctrl:              ctrl,
		sliderWidth:       sliderWidth,
	}
}

func (m *Model) buildEffectsList() []list.Item {
	// Master always comes first
	items := []list.Item{
		newEffectItem("master", "Master", m.MasterEnabled, true),
	}

	// Build rest of list from EffectsOrder
	for _, effectID := range m.EffectsOrder {
		switch effectID {
		case "filter":
			items = append(items, newEffectItem("filter", "Filter", m.FilterEnabled, true))
		case "overdrive":
			items = append(items, newEffectItem("overdrive", "Overdrive", m.OverdriveEnabled, true))
		case "bitcrush":
			items = append(items, newEffectItem("bitcrush", "Bitcrush", m.BitcrushEnabled, true))
		case "granular":
			items = append(items, newEffectItem("granular", "Granular", m.GranularEnabled, true))
		case "reverb":
			items = append(items, newEffectItem("reverb", "Reverb", m.ReverbEnabled, true))
		case "delay":
			items = append(items, newEffectItem("delay", "Delay", m.DelayEnabled, true))
		}
	}

	return items
}

func (m *Model) buildParameterList(section string) []list.Item {
	var items []list.Item

	switch section {
	case "master":
		items = []list.Item{
			newParameterItem("masterEnabled", "Master", 0, 0, 0, ctrlMasterEnabled, true, m.MasterEnabled, m.sliderWidth),
			newParameterItem("gain", "Gain", m.Gain, 0, 2, ctrlGain, false, false, m.sliderWidth),
			newParameterItem("freezeLength", "Loop Length", m.InputFreezeLength, 0.05, 0.5, ctrlInputFreezeLen, false, false, m.sliderWidth),
			newParameterItem("freeze", "Input Freeze", 0, 0, 0, ctrlInputFreeze, true, m.InputFrozen, m.sliderWidth),
			newParameterItemWithDesc("blendMode", "Blend Mode", m.formatBlendMode(), ctrlBlendMode, m.sliderWidth),
			newParameterItem("dryWet", "Dry/Wet", m.DryWet, 0, 1, ctrlDryWet, false, false, m.sliderWidth),
			newParameterItemWithDesc("effectsOrder", "Effects Order", m.formatEffectsOrder(), ctrlEffectsOrder, m.sliderWidth),
		}
	case "filter":
		items = []list.Item{
			newParameterItem("enabled", "Filter", 0, 0, 0, ctrlFilterEnabled, true, m.FilterEnabled, m.sliderWidth),
			newParameterItem("amount", "Amount", m.FilterAmount, 0, 1, ctrlFilterAmount, false, false, m.sliderWidth),
			newParameterItem("cutoff", "Cutoff", m.FilterCutoff, 200, 8000, ctrlFilterCutoff, false, false, m.sliderWidth),
			newParameterItem("resonance", "Resonance", m.FilterResonance, 0, 1, ctrlFilterResonance, false, false, m.sliderWidth),
		}
	case "overdrive":
		items = []list.Item{
			newParameterItem("enabled", "Overdrive", 0, 0, 0, ctrlOverdriveEnabled, true, m.OverdriveEnabled, m.sliderWidth),
			newParameterItem("drive", "Drive", m.OverdriveDrive, 0, 1, ctrlOverdriveDrive, false, false, m.sliderWidth),
			newParameterItem("tone", "Tone", m.OverdriveTone, 0, 1, ctrlOverdriveTone, false, false, m.sliderWidth),
			newParameterItem("bias", "Bias", m.OverdriveBias, -1, 1, ctrlOverdriveBias, false, false, m.sliderWidth),
			newParameterItem("mix", "Mix", m.OverdriveMix, 0, 1, ctrlOverdriveMix, false, false, m.sliderWidth),
		}
	case "bitcrush":
		items = []list.Item{
			newParameterItem("enabled", "Bitcrush", 0, 0, 0, ctrlBitcrushEnabled, true, m.BitcrushEnabled, m.sliderWidth),
			newParameterItem("bitDepth", "Bit Depth", m.BitDepth, 4, 16, ctrlBitDepth, false, false, m.sliderWidth),
			newParameterItem("sampleRate", "Sample Rate", m.BitcrushSampleRate, 1000, 44100, ctrlBitcrushSampleRate, false, false, m.sliderWidth),
			newParameterItem("drive", "Drive", m.BitcrushDrive, 0, 1, ctrlBitcrushDrive, false, false, m.sliderWidth),
			newParameterItem("mix", "Mix", m.BitcrushMix, 0, 1, ctrlBitcrushMix, false, false, m.sliderWidth),
		}
	case "granular":
		items = []list.Item{
			newParameterItem("enabled", "Granular", 0, 0, 0, ctrlGranularEnabled, true, m.GranularEnabled, m.sliderWidth),
			newParameterItem("density", "Density", m.GranularDensity, 1, 50, ctrlGranularDensity, false, false, m.sliderWidth),
			newParameterItem("size", "Grain Size", m.GranularSize, 0.01, 0.5, ctrlGranularSize, false, false, m.sliderWidth),
			newParameterItem("pitchScat", "Pitch Scatter", m.GranularPitchScatter, 0, 1, ctrlGranularPitchScatter, false, false, m.sliderWidth),
			newParameterItem("posScat", "Position Scatter", m.GranularPosScatter, 0, 1, ctrlGranularPosScatter, false, false, m.sliderWidth),
			newParameterItem("mix", "Mix", m.GranularMix, 0, 1, ctrlGranularMix, false, false, m.sliderWidth),
			newParameterItem("freeze", "Grain Freeze", 0, 0, 0, ctrlGranularFreeze, true, m.GranularFrozen, m.sliderWidth),
			newParameterItemWithDesc("intensity", "Intensity", m.formatGrainIntensity(), ctrlGrainIntensity, m.sliderWidth),
		}
	case "reverb":
		items = []list.Item{
			newParameterItem("enabled", "Reverb", 0, 0, 0, ctrlReverbEnabled, true, m.ReverbEnabled, m.sliderWidth),
			newParameterItem("decay", "Decay", m.ReverbDecayTime, 0.5, 10, ctrlReverbDecayTime, false, false, m.sliderWidth),
			newParameterItem("mix", "Mix", m.ReverbMix, 0, 1, ctrlReverbMix, false, false, m.sliderWidth),
		}
	case "delay":
		items = []list.Item{
			newParameterItem("enabled", "Delay", 0, 0, 0, ctrlDelayEnabled, true, m.DelayEnabled, m.sliderWidth),
			newParameterItem("time", "Time", m.DelayTime, 0.1, 1, ctrlDelayTime, false, false, m.sliderWidth),
			newParameterItem("decay", "Decay", m.DelayDecayTime, 0.5, 10, ctrlDelayDecayTime, false, false, m.sliderWidth),
			newParameterItem("modRate", "Mod Rate", m.ModRate, 0.1, 10, ctrlModRate, false, false, m.sliderWidth),
			newParameterItem("modDepth", "Mod Depth", m.ModDepth, 0, 1, ctrlModDepth, false, false, m.sliderWidth),
			newParameterItem("mix", "Mix", m.DelayMix, 0, 1, ctrlDelayMix, false, false, m.sliderWidth),
		}
	}

	return items
}

func formatSliderValue(value, min, max float32, sliderWidth int) string {
	normalized := (value - min) / (max - min)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	filled := int(normalized * float32(sliderWidth))
	if filled > sliderWidth {
		filled = sliderWidth
	}

	bar := strings.Repeat("=", filled) + strings.Repeat("-", sliderWidth-filled)
	return fmt.Sprintf("%s %s", bar, formatValue(value, min, max))
}

func (m *Model) formatSliderValue(label string, value, min, max float32) string {
	norm := (value - min) / (max - min)
	if norm < 0 {
		norm = 0
	}
	if norm > 1 {
		norm = 1
	}

	barWidth := 10
	filled := int(norm * float32(barWidth))
	bar := strings.Repeat("=", filled) + strings.Repeat("-", barWidth-filled)

	valueStr := formatValue(value, min, max)
	return bar + " " + valueStr
}

func (m *Model) formatEffectsOrder() string {
	var parts []string
	for i, effect := range m.EffectsOrder {
		if m.effectsOrderEditMode && i == m.selectedEffectIndex {
			if m.effectGrabbed {
				// Use double brackets when grabbed
				parts = append(parts, fmt.Sprintf("[[%s]]", effect))
			} else {
				// Use single brackets when just selected
				parts = append(parts, fmt.Sprintf("[%s]", effect))
			}
		} else {
			parts = append(parts, effect)
		}
	}

	result := strings.Join(parts, " → ")

	if m.effectsOrderEditMode {
		if m.effectGrabbed {
			result += "\nh/l: move · esc: ungrab"
		} else {
			result += "\nh/l: select · enter: grab · esc: done"
		}
	}

	return result
}

func (m *Model) formatGrainIntensity() string {
	options := []string{"subtle", "pronounced", "extreme"}
	var parts []string
	for _, opt := range options {
		if opt == m.GrainIntensity {
			parts = append(parts, fmt.Sprintf("[%s]", opt))
		} else {
			parts = append(parts, opt)
		}
	}
	return strings.Join(parts, " ")
}

func (m *Model) formatBlendMode() string {
	options := []string{"mirror", "complement", "transform"}
	var parts []string
	for i, opt := range options {
		if i == m.BlendMode {
			parts = append(parts, fmt.Sprintf("[%s]", opt))
		} else {
			parts = append(parts, opt)
		}
	}
	return strings.Join(parts, " ")
}
