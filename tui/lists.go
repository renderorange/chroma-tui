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
	id          string
	title       string
	description string
	ctrl        control
	isToggle    bool
	isActive    bool
}

func (i parameterItem) Title() string       { return i.title }
func (i parameterItem) Description() string { return i.description }
func (i parameterItem) FilterValue() string { return i.title }

func newEffectItem(id, title string, enabled bool) effectItem {
	status := "[DISABLED]"
	if enabled {
		status = "[ENABLED]"
	}
	return effectItem{
		id:          id,
		title:       title,
		description: status,
		enabled:     enabled,
	}
}

func newParameterItem(id, title, description string, ctrl control, isToggle, isActive bool) parameterItem {
	prefix := "   "
	if isToggle {
		if isActive {
			prefix = "[X]"
		} else {
			prefix = "[ ]"
		}
	}
	return parameterItem{
		id:          id,
		title:       prefix + " " + title,
		description: description,
		ctrl:        ctrl,
		isToggle:    isToggle,
		isActive:    isActive,
	}
}

func (m *Model) buildEffectsList() []list.Item {
	items := []list.Item{
		newEffectItem("input", "Input", m.InputFrozen),
		newEffectItem("filter", "Filter", m.FilterEnabled),
		newEffectItem("overdrive", "Overdrive", m.OverdriveEnabled),
		newEffectItem("bitcrush", "Bitcrush", m.BitcrushEnabled),
		newEffectItem("granular", "Granular", m.GranularEnabled),
		newEffectItem("reverb", "Reverb", m.ReverbEnabled),
		newEffectItem("delay", "Delay", m.DelayEnabled),
		newEffectItem("global", "Global", true),
	}
	return items
}

func (m *Model) buildParameterList(section string) []list.Item {
	var items []list.Item

	switch section {
	case "input":
		items = []list.Item{
			newParameterItem("gain", "Gain", m.formatSliderValue("Gain", m.Gain, 0, 2), ctrlGain, false, false),
			newParameterItem("freezeLength", "Loop Length", m.formatSliderValue("Loop", m.InputFreezeLength, 0.05, 0.5), ctrlInputFreezeLen, false, false),
			newParameterItem("freeze", "Input Freeze", "Toggle loop freeze", ctrlInputFreeze, true, m.InputFrozen),
		}
	case "filter":
		items = []list.Item{
			newParameterItem("enabled", "Filter", "Enable/disable filter", ctrlFilterEnabled, true, m.FilterEnabled),
			newParameterItem("amount", "Amount", m.formatSliderValue("Amount", m.FilterAmount, 0, 1), ctrlFilterAmount, false, false),
			newParameterItem("cutoff", "Cutoff", m.formatSliderValue("Cutoff", m.FilterCutoff, 200, 8000), ctrlFilterCutoff, false, false),
			newParameterItem("resonance", "Resonance", m.formatSliderValue("Resonance", m.FilterResonance, 0, 1), ctrlFilterResonance, false, false),
		}
	case "overdrive":
		items = []list.Item{
			newParameterItem("enabled", "Overdrive", "Enable/disable overdrive", ctrlOverdriveEnabled, true, m.OverdriveEnabled),
			newParameterItem("drive", "Drive", m.formatSliderValue("Drive", m.OverdriveDrive, 0, 1), ctrlOverdriveDrive, false, false),
			newParameterItem("tone", "Tone", m.formatSliderValue("Tone", m.OverdriveTone, 0, 1), ctrlOverdriveTone, false, false),
			newParameterItem("bias", "Bias", m.formatSliderValue("Bias", m.OverdriveBias, -1, 1), ctrlOverdriveBias, false, false),
			newParameterItem("mix", "Mix", m.formatSliderValue("Mix", m.OverdriveMix, 0, 1), ctrlOverdriveMix, false, false),
		}
	case "bitcrush":
		items = []list.Item{
			newParameterItem("enabled", "Bitcrush", "Enable/disable bitcrush", ctrlBitcrushEnabled, true, m.BitcrushEnabled),
			newParameterItem("bitDepth", "Bit Depth", m.formatSliderValue("BitDepth", m.BitDepth, 4, 16), ctrlBitDepth, false, false),
			newParameterItem("sampleRate", "Sample Rate", m.formatSliderValue("SampleRate", m.BitcrushSampleRate, 1000, 44100), ctrlBitcrushSampleRate, false, false),
			newParameterItem("drive", "Drive", m.formatSliderValue("Drive", m.BitcrushDrive, 0, 1), ctrlBitcrushDrive, false, false),
			newParameterItem("mix", "Mix", m.formatSliderValue("Mix", m.BitcrushMix, 0, 1), ctrlBitcrushMix, false, false),
		}
	case "granular":
		items = []list.Item{
			newParameterItem("enabled", "Granular", "Enable/disable granular", ctrlGranularEnabled, true, m.GranularEnabled),
			newParameterItem("density", "Density", m.formatSliderValue("Density", m.GranularDensity, 1, 50), ctrlGranularDensity, false, false),
			newParameterItem("size", "Grain Size", m.formatSliderValue("Size", m.GranularSize, 0.01, 0.5), ctrlGranularSize, false, false),
			newParameterItem("pitchScat", "Pitch Scatter", m.formatSliderValue("PitchScat", m.GranularPitchScatter, 0, 1), ctrlGranularPitchScatter, false, false),
			newParameterItem("posScat", "Position Scatter", m.formatSliderValue("PosScat", m.GranularPosScatter, 0, 1), ctrlGranularPosScatter, false, false),
			newParameterItem("mix", "Mix", m.formatSliderValue("Mix", m.GranularMix, 0, 1), ctrlGranularMix, false, false),
			newParameterItem("freeze", "Grain Freeze", "Toggle grain freeze", ctrlGranularFreeze, true, m.GranularFrozen),
			newParameterItem("intensity", "Intensity", "Current: "+m.GrainIntensity, ctrlGrainIntensity, false, false),
		}
	case "reverb":
		items = []list.Item{
			newParameterItem("enabled", "Reverb", "Enable/disable reverb", ctrlReverbEnabled, true, m.ReverbEnabled),
			newParameterItem("decay", "Decay", m.formatSliderValue("Decay", m.ReverbDecayTime, 0.5, 10), ctrlReverbDecayTime, false, false),
			newParameterItem("mix", "Mix", m.formatSliderValue("Mix", m.ReverbMix, 0, 1), ctrlReverbMix, false, false),
		}
	case "delay":
		items = []list.Item{
			newParameterItem("enabled", "Delay", "Enable/disable delay", ctrlDelayEnabled, true, m.DelayEnabled),
			newParameterItem("time", "Time", m.formatSliderValue("Time", m.DelayTime, 0.1, 1), ctrlDelayTime, false, false),
			newParameterItem("decay", "Decay", m.formatSliderValue("Decay", m.DelayDecayTime, 0.5, 10), ctrlDelayDecayTime, false, false),
			newParameterItem("modRate", "Mod Rate", m.formatSliderValue("ModRate", m.ModRate, 0.1, 10), ctrlModRate, false, false),
			newParameterItem("modDepth", "Mod Depth", m.formatSliderValue("ModDepth", m.ModDepth, 0, 1), ctrlModDepth, false, false),
			newParameterItem("mix", "Mix", m.formatSliderValue("Mix", m.DelayMix, 0, 1), ctrlDelayMix, false, false),
		}
	case "global":
		modeNames := []string{"MIRROR", "COMPLEMENT", "TRANSFORM"}
		modeDesc := "Current: " + modeNames[m.BlendMode]
		items = []list.Item{
			newParameterItem("blendMode", "Blend Mode", modeDesc, ctrlBlendMode, false, false),
			newParameterItem("dryWet", "Dry/Wet", m.formatSliderValue("Dry/Wet", m.DryWet, 0, 1), ctrlDryWet, false, false),
			newParameterItem("effectsOrder", "Effects Order", m.formatEffectsOrder(), ctrlEffectsOrder, false, false),
		}
	}

	return items
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
		parts = append(parts, fmt.Sprintf("%d.%s", i+1, effect))
	}
	return strings.Join(parts, " ")
}
