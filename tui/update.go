package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"math"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Initialize lists on first window size
		if len(m.effectsList.Items()) == 0 {
			m.InitLists(msg.Width, msg.Height)
			return m, nil
		}
		leftWidth, rightWidth, listHeight := panelDimensions(msg.Width, msg.Height)
		m.effectsList.SetSize(leftWidth, listHeight)
		m.parameterList.SetSize(rightWidth, listHeight)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Delegate non-key messages to both lists
	var cmd1, cmd2 tea.Cmd
	m.effectsList, cmd1 = m.effectsList.Update(msg)
	m.parameterList, cmd2 = m.parameterList.Update(msg)
	return m, tea.Batch(cmd1, cmd2)
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// App-level keys - handle first and return early
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "enter":
		return m.handleEnterKey()

	case "esc":
		return m.handleEscKey()

	case "S":
		m.showStatus = !m.showStatus
		m.effectsList.SetShowStatusBar(m.showStatus)
		m.parameterList.SetShowStatusBar(m.showStatus)
		return m, nil

	case "P":
		m.showPagination = !m.showPagination
		m.effectsList.SetShowPagination(m.showPagination)
		m.parameterList.SetShowPagination(m.showPagination)
		return m, nil

	case "T":
		m.showTitle = !m.showTitle
		m.effectsList.SetShowTitle(m.showTitle)
		m.parameterList.SetShowTitle(m.showTitle)
		return m, nil

	case "i":
		m.toggleGrainIntensity()
		return m, nil

	case "1":
		m.setBlendMode(0)
		return m, nil
	case "2":
		m.setBlendMode(1)
		return m, nil
	case "3":
		m.setBlendMode(2)
		return m, nil
	}

	// Parameter adjustment keys - only in parameter mode
	if m.navigationMode == modeParameterList {
		switch msg.String() {
		case "left", "h":
			m.adjustSelectedParameter(-0.05)
			m.refreshParameterList()
			return m, nil

		case "right", "l":
			m.adjustSelectedParameter(0.05)
			m.refreshParameterList()
			return m, nil
		}

		// Effects order keys - only in parameter mode with global section
		if m.currentSection == "global" {
			switch msg.String() {
			case "pgup", "pgdown":
				m.handleEffectsOrderKeys(msg)
				return m, nil
			case "r":
				m.handleEffectsOrderKeys(msg)
				m.refreshParameterList()
				return m, nil
			}
		}
	}

	// Delegate all other keys to the active list
	var cmd tea.Cmd
	if m.navigationMode == modeEffectsList {
		m.effectsList, cmd = m.effectsList.Update(msg)
		m.syncParameterPanel()
	} else if m.navigationMode == modeParameterList {
		m.parameterList, cmd = m.parameterList.Update(msg)
	}
	return m, cmd
}

func (m *Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.navigationMode {
	case modeEffectsList:
		// Switch focus to parameter panel
		m.navigationMode = modeParameterList
		m.syncParameterPanel()

	case modeParameterList:
		// Toggle parameter if it's a toggle item
		idx := m.parameterList.Index()
		items := m.parameterList.Items()
		if idx >= 0 && idx < len(items) {
			item := items[idx]
			if param, ok := item.(parameterItem); ok {
				if param.isToggle {
					m.toggleByControl(param.ctrl)
				}
				m.refreshParameterList()
			}
		}
	}
	return m, nil
}

func (m *Model) handleEscKey() (tea.Model, tea.Cmd) {
	switch m.navigationMode {
	case modeParameterList:
		// Switch focus back to effects panel
		m.navigationMode = modeEffectsList
	}
	return m, nil
}

func (m *Model) adjustSelectedParameter(delta float32) {
	idx := m.parameterList.Index()
	items := m.parameterList.Items()
	if idx >= 0 && idx < len(items) {
		item := items[idx]
		if param, ok := item.(parameterItem); ok {
			m.focused = param.ctrl
			m.adjustFocused(delta)
		}
	}
}

func (m *Model) toggleByControl(ctrl control) {
	m.focused = ctrl
	m.toggleFocused()
}

func (m *Model) adjustFocused(delta float32) {
	switch m.focused {
	case ctrlGain:
		m.Gain = clamp(m.Gain+delta*2, 0, 2)
		if err := m.client.SetGain(m.Gain); err != nil {
			// Log error but continue - UDP is best-effort
		}
	case ctrlInputFreezeLen:
		m.InputFreezeLength = clamp(m.InputFreezeLength+delta*0.45, 0.05, 0.5)
		if err := m.client.SetInputFreezeLength(m.InputFreezeLength); err != nil {
		}
	case ctrlFilterAmount:
		m.FilterAmount = clamp(m.FilterAmount+delta, 0, 1)
		if err := m.client.SetFilterAmount(m.FilterAmount); err != nil {
		}
	case ctrlFilterCutoff:
		m.FilterCutoff = clamp(m.FilterCutoff+delta*7800, 200, 8000)
		if err := m.client.SetFilterCutoff(m.FilterCutoff); err != nil {
		}
	case ctrlFilterResonance:
		m.FilterResonance = clamp(m.FilterResonance+delta, 0, 1)
		if err := m.client.SetFilterResonance(m.FilterResonance); err != nil {
		}
	case ctrlOverdriveDrive:
		m.OverdriveDrive = clamp(m.OverdriveDrive+delta, 0, 1)
		if err := m.client.SetOverdriveDrive(m.OverdriveDrive); err != nil {
		}
	case ctrlOverdriveTone:
		m.OverdriveTone = clamp(m.OverdriveTone+delta, 0, 1)
		if err := m.client.SetOverdriveTone(m.OverdriveTone); err != nil {
		}
	case ctrlOverdriveBias:
		m.OverdriveBias = clamp(m.OverdriveBias+delta, -1, 1)
		if err := m.client.SetOverdriveBias(m.OverdriveBias); err != nil {
		}
	case ctrlOverdriveMix:
		m.OverdriveMix = clamp(m.OverdriveMix+delta, 0, 1)
		if err := m.client.SetOverdriveMix(m.OverdriveMix); err != nil {
		}
	case ctrlBitDepth:
		m.BitDepth = clamp(m.BitDepth+delta*12, 4, 16)
		if err := m.client.SetBitDepth(m.BitDepth); err != nil {
		}
	case ctrlBitcrushSampleRate:
		m.BitcrushSampleRate = clamp(m.BitcrushSampleRate+delta*43100, 1000, 44100)
		if err := m.client.SetBitcrushSampleRate(m.BitcrushSampleRate); err != nil {
		}
	case ctrlBitcrushDrive:
		m.BitcrushDrive = clamp(m.BitcrushDrive+delta, 0, 1)
		if err := m.client.SetBitcrushDrive(m.BitcrushDrive); err != nil {
		}
	case ctrlBitcrushMix:
		m.BitcrushMix = clamp(m.BitcrushMix+delta, 0, 1)
		if err := m.client.SetBitcrushMix(m.BitcrushMix); err != nil {
		}
	case ctrlGranularDensity:
		m.GranularDensity = adjustLogarithmic(m.GranularDensity, delta*0.8, 1, 50)
		if err := m.client.SetGranularDensity(m.GranularDensity); err != nil {
		}
	case ctrlGranularSize:
		m.GranularSize = adjustLogarithmic(m.GranularSize, delta*0.5, 0.01, 2.0)
		if err := m.client.SetGranularSize(m.GranularSize); err != nil {
		}
	case ctrlGranularPitchScatter:
		m.GranularPitchScatter = clamp(m.GranularPitchScatter+delta, 0, 1)
		if err := m.client.SetGranularPitchScatter(m.GranularPitchScatter); err != nil {
		}
	case ctrlGranularPosScatter:
		m.GranularPosScatter = clamp(m.GranularPosScatter+delta, 0, 1)
		if err := m.client.SetGranularPosScatter(m.GranularPosScatter); err != nil {
		}
	case ctrlGranularMix:
		m.GranularMix = clamp(m.GranularMix+delta, 0, 1)
		if err := m.client.SetGranularMix(m.GranularMix); err != nil {
		}
	case ctrlReverbDecayTime:
		m.ReverbDecayTime = clamp(m.ReverbDecayTime+delta*9.5, 0.5, 10)
		if err := m.client.SetReverbDecayTime(m.ReverbDecayTime); err != nil {
		}
	case ctrlReverbMix:
		m.ReverbMix = clamp(m.ReverbMix+delta, 0, 1)
		if err := m.client.SetReverbMix(m.ReverbMix); err != nil {
		}
	case ctrlDelayTime:
		m.DelayTime = clamp(m.DelayTime+delta*1.99, 0.01, 2.0)
		if err := m.client.SetDelayTime(m.DelayTime); err != nil {
		}
	case ctrlDelayDecayTime:
		m.DelayDecayTime = clamp(m.DelayDecayTime+delta*4.9, 0.1, 5.0)
		if err := m.client.SetDelayDecayTime(m.DelayDecayTime); err != nil {
		}
	case ctrlModRate:
		m.ModRate = clamp(m.ModRate+delta*9.9, 0.1, 10.0)
		if err := m.client.SetModRate(m.ModRate); err != nil {
		}
	case ctrlModDepth:
		m.ModDepth = adjustLogarithmic(m.ModDepth, delta*0.5, 0, 1)
		if err := m.client.SetModDepth(m.ModDepth); err != nil {
		}
	case ctrlDelayMix:
		m.DelayMix = clamp(m.DelayMix+delta, 0, 1)
		if err := m.client.SetDelayMix(m.DelayMix); err != nil {
		}
	case ctrlDryWet:
		m.DryWet = clamp(m.DryWet+delta, 0, 1)
		if err := m.client.SetDryWet(m.DryWet); err != nil {
		}
	}
}

func (m *Model) toggleFocused() {
	switch m.focused {
	case ctrlInputFreeze:
		m.InputFrozen = !m.InputFrozen
		if err := m.client.SetInputFreeze(m.InputFrozen); err != nil {
		}
	case ctrlFilterEnabled:
		m.FilterEnabled = !m.FilterEnabled
		if err := m.client.SetFilterEnabled(m.FilterEnabled); err != nil {
		}
	case ctrlOverdriveEnabled:
		m.OverdriveEnabled = !m.OverdriveEnabled
		if err := m.client.SetOverdriveEnabled(m.OverdriveEnabled); err != nil {
		}
	case ctrlBitcrushEnabled:
		m.BitcrushEnabled = !m.BitcrushEnabled
		if err := m.client.SetBitcrushEnabled(m.BitcrushEnabled); err != nil {
		}
	case ctrlGranularEnabled:
		m.GranularEnabled = !m.GranularEnabled
		if err := m.client.SetGranularEnabled(m.GranularEnabled); err != nil {
		}
	case ctrlGranularFreeze:
		m.GranularFrozen = !m.GranularFrozen
		if err := m.client.SetGranularFreeze(m.GranularFrozen); err != nil {
		}
	case ctrlReverbEnabled:
		m.ReverbEnabled = !m.ReverbEnabled
		if err := m.client.SetReverbEnabled(m.ReverbEnabled); err != nil {
		}
	case ctrlDelayEnabled:
		m.DelayEnabled = !m.DelayEnabled
		if err := m.client.SetDelayEnabled(m.DelayEnabled); err != nil {
		}
	}
}

func (m *Model) setBlendMode(mode int) {
	m.BlendMode = mode
	if err := m.client.SetBlendMode(mode); err != nil {
	}
}

func (m *Model) toggleGrainIntensity() {
	switch m.GrainIntensity {
	case "subtle":
		m.GrainIntensity = "pronounced"
	case "pronounced":
		m.GrainIntensity = "extreme"
	case "extreme":
		m.GrainIntensity = "subtle"
	default:
		m.GrainIntensity = "subtle"
	}
	if err := m.client.SetGrainIntensity(m.GrainIntensity); err != nil {
	}
}

func (m *Model) handleEffectsOrderKeys(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyPgUp:
		// Move selected effect up in order
		if m.selectedEffectIndex > 0 {
			order := m.GetEffectsOrder()
			// Swap with previous
			order[m.selectedEffectIndex], order[m.selectedEffectIndex-1] =
				order[m.selectedEffectIndex-1], order[m.selectedEffectIndex]
			m.SetEffectsOrder(order)
			m.selectedEffectIndex-- // Keep selection on moved effect
			// Trigger OSC update
			if err := m.client.SetEffectsOrder(order); err != nil {
			}
		}
	case tea.KeyPgDown:
		// Move selected effect down in order
		if m.selectedEffectIndex < len(m.EffectsOrder)-1 {
			order := m.GetEffectsOrder()
			// Swap with next
			order[m.selectedEffectIndex], order[m.selectedEffectIndex+1] =
				order[m.selectedEffectIndex+1], order[m.selectedEffectIndex]
			m.SetEffectsOrder(order)
			m.selectedEffectIndex++ // Keep selection on moved effect
			// Trigger OSC update
			if err := m.client.SetEffectsOrder(order); err != nil {
			}
		}
	case tea.KeyRunes:
		switch msg.Runes[0] {
		case 'r':
			// Reset to default order
			defaultOrder := []string{"filter", "overdrive", "bitcrush", "granular", "reverb", "delay"}
			m.SetEffectsOrder(defaultOrder)
			m.selectedEffectIndex = 0
			// Trigger OSC update
			if err := m.client.SetEffectsOrder(defaultOrder); err != nil {
			}
		}
	}
}

func adjustLogarithmic(current, delta, min, max float32) float32 {
	if current <= 0 || min <= 0 {
		return clamp(current+delta, min, max) // Fallback for edge cases
	}

	logCurrent := math.Log10(float64(current))
	logMin := math.Log10(float64(min))
	logMax := math.Log10(float64(max))

	newLog := logCurrent + float64(delta)*0.1*(logMax-logMin)
	newLog = math.Max(logMin, math.Min(logMax, newLog))

	return float32(math.Pow(10, newLog))
}

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
