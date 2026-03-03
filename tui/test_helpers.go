package tui

// Test helper methods for accessing private TUI functionality in tests

// SetFocused allows tests to set the focused control
func (m *Model) SetFocused(ctrl control) {
	m.focused = ctrl
}

// AdjustFocused allows tests to call the private adjustFocused method
func (m *Model) AdjustFocused(delta float32) {
	m.adjustFocused(delta)
}

// ToggleFocused allows tests to call the private toggleFocused method
func (m *Model) ToggleFocused() {
	m.toggleFocused()
}

// SetBlendMode allows tests to call the private setBlendMode method
func (m *Model) SetBlendMode(mode int) {
	m.setBlendMode(mode)
}

// SetNavigationMode allows tests to set the navigation mode
func (m *Model) SetNavigationMode(mode int) {
	m.navigationMode = navigationMode(mode)
}

// SetCurrentSection allows tests to set the current section
func (m *Model) SetCurrentSection(section string) {
	m.currentSection = section
}

// SetScreenForTesting allows tests to set the screen state directly
func (m *Model) SetScreenForTesting(screen int) {
	m.screen = screenState(screen)
}

// TestConstants provides access to control constants for tests
var (
	TestCtrlGain                 = ctrlGain
	TestCtrlInputFreeze          = ctrlInputFreeze
	TestCtrlInputFreezeLen       = ctrlInputFreezeLen
	TestCtrlFilterEnabled        = ctrlFilterEnabled
	TestCtrlFilterAmount         = ctrlFilterAmount
	TestCtrlFilterCutoff         = ctrlFilterCutoff
	TestCtrlFilterResonance      = ctrlFilterResonance
	TestCtrlOverdriveEnabled     = ctrlOverdriveEnabled
	TestCtrlOverdriveDrive       = ctrlOverdriveDrive
	TestCtrlOverdriveTone        = ctrlOverdriveTone
	TestCtrlOverdriveBias        = ctrlOverdriveBias
	TestCtrlOverdriveMix         = ctrlOverdriveMix
	TestCtrlGranularEnabled      = ctrlGranularEnabled
	TestCtrlGranularDensity      = ctrlGranularDensity
	TestCtrlGranularSize         = ctrlGranularSize
	TestCtrlGranularPitchScatter = ctrlGranularPitchScatter
	TestCtrlGranularPosScatter   = ctrlGranularPosScatter
	TestCtrlGranularMix          = ctrlGranularMix
	TestCtrlGranularFreeze       = ctrlGranularFreeze
	TestCtrlBitcrushEnabled      = ctrlBitcrushEnabled
	TestCtrlBitDepth             = ctrlBitDepth
	TestCtrlBitcrushSampleRate   = ctrlBitcrushSampleRate
	TestCtrlBitcrushDrive        = ctrlBitcrushDrive
	TestCtrlBitcrushMix          = ctrlBitcrushMix
	TestCtrlReverbEnabled        = ctrlReverbEnabled
	TestCtrlReverbDecayTime      = ctrlReverbDecayTime
	TestCtrlReverbMix            = ctrlReverbMix
	TestCtrlDelayEnabled         = ctrlDelayEnabled
	TestCtrlDelayTime            = ctrlDelayTime
	TestCtrlDelayDecayTime       = ctrlDelayDecayTime
	TestCtrlModRate              = ctrlModRate
	TestCtrlModDepth             = ctrlModDepth
	TestCtrlDelayMix             = ctrlDelayMix
	TestCtrlBlendMode            = ctrlBlendMode
	TestCtrlDryWet               = ctrlDryWet
)

// GetScreenForTesting returns the current screen state for testing
func (m *Model) GetScreenForTesting() int {
	return int(m.screen)
}

// GetEffectsListItems returns the effects list items for testing
func (m *Model) GetEffectsListItems() []interface{} {
	items := m.effectsList.Items()
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = item
	}
	return result
}

// GetWidthForTesting returns the model width for testing
func (m *Model) GetWidthForTesting() int {
	return m.width
}

// GetHeightForTesting returns the model height for testing
func (m *Model) GetHeightForTesting() int {
	return m.height
}

// ShowCommandPaletteForTesting sets the command palette visibility for testing
func (m *Model) ShowCommandPaletteForTesting(show bool) {
	m.showCommandPalette = show
}

// SetCommandPaletteTextForTesting sets the command palette text for testing
func (m *Model) SetCommandPaletteTextForTesting(text string) {
	m.commandPaletteText = text
}

// IsCommandPaletteVisibleForTesting returns whether palette is visible for testing
func (m *Model) IsCommandPaletteVisibleForTesting() bool {
	return m.showCommandPalette
}

// GetCommandPaletteTextForTesting returns the palette text for testing
func (m *Model) GetCommandPaletteTextForTesting() string {
	return m.commandPaletteText
}
