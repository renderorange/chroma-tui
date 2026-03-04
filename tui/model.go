package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/renderorange/chroma/chroma-control/config"
	"github.com/renderorange/chroma/chroma-control/osc"
)

var (
	colorPrimary       = lipgloss.Color("#888888")
	colorSecondary     = lipgloss.Color("#666666")
	colorAccent        = lipgloss.Color("#AAAAAA")
	colorBackground    = lipgloss.Color("#000000")
	colorTextNormal    = lipgloss.Color("#CCCCCC")
	colorTextMuted     = lipgloss.Color("#444444")
	colorTextHighlight = lipgloss.Color("#888888")
	colorTextSuccess   = lipgloss.Color("#888888")
	colorTextError     = lipgloss.Color("#AAAAAA")
)

type control int

type navigationMode int

const (
	modeEffectsList navigationMode = iota
	modeParameterList
)

type screenState int

const (
	screenSplash screenState = iota
	screenMain
	screenSettings
	screenHelp
	screenPresetBrowser
)

type splashOption int

const (
	splashLast splashOption = iota
	splashNew
	splashLoad
	splashQuit
)

type settingsMenuItem int

const (
	settingsSave settingsMenuItem = iota
	settingsSaveAs
	settingsLoad
	settingsReset
	settingsBack
)

type presetBrowserMode int

const (
	browserModeList presetBrowserMode = iota
	browserModeConfirmDelete
	browserModeSaveAs
)

type presetBrowserState struct {
	mode        presetBrowserMode
	presets     []string
	selectedIdx int
	inputBuffer string
	confirmName string
	errorMsg    string
}

const (
	ctrlMasterEnabled control = iota
	ctrlGain
	ctrlInputFreezeLen
	ctrlInputFreeze
	ctrlFilterEnabled
	ctrlFilterAmount
	ctrlFilterCutoff
	ctrlFilterResonance
	ctrlOverdriveEnabled
	ctrlOverdriveDrive
	ctrlOverdriveTone
	ctrlOverdriveBias
	ctrlOverdriveMix
	ctrlBitcrushEnabled
	ctrlBitDepth
	ctrlBitcrushSampleRate
	ctrlBitcrushDrive
	ctrlBitcrushMix
	ctrlGranularEnabled
	ctrlGranularDensity
	ctrlGranularSize
	ctrlGranularPitchScatter
	ctrlGranularPosScatter
	ctrlGranularMix
	ctrlGranularFreeze
	ctrlGrainIntensity
	ctrlReverbEnabled
	ctrlReverbDecayTime
	ctrlReverbMix
	ctrlDelayEnabled
	ctrlDelayTime
	ctrlDelayDecayTime
	ctrlModRate
	ctrlModDepth
	ctrlDelayMix
	ctrlBlendMode
	ctrlDryWet
	ctrlEffectsOrder
	ctrlEffectsMoveUp
	ctrlEffectsMoveDown
	ctrlCount
)

type Model struct {
	// State
	MasterEnabled        bool
	Gain                 float32
	InputFrozen          bool
	InputFreezeLength    float32
	FilterEnabled        bool
	FilterAmount         float32
	FilterCutoff         float32
	FilterResonance      float32
	OverdriveEnabled     bool
	OverdriveDrive       float32
	OverdriveTone        float32
	OverdriveBias        float32
	OverdriveMix         float32
	BitcrushEnabled      bool
	BitDepth             float32
	BitcrushSampleRate   float32
	BitcrushDrive        float32
	BitcrushMix          float32
	GranularEnabled      bool
	GranularDensity      float32
	GranularSize         float32
	GranularPitchScatter float32
	GranularPosScatter   float32
	GranularMix          float32
	GranularFrozen       bool
	GrainIntensity       string // "subtle", "pronounced", or "extreme"
	ReverbEnabled        bool
	ReverbDecayTime      float32
	ReverbMix            float32
	DelayEnabled         bool
	DelayTime            float32
	DelayDecayTime       float32
	ModRate              float32
	ModDepth             float32
	DelayMix             float32
	BlendMode            int
	DryWet               float32
	EffectsOrder         []string // Current effects processing order

	// UI state
	focused              control
	selectedEffectIndex  int  // Index of currently selected effect for reordering
	effectsOrderEditMode bool // Whether we're in effects reorder mode
	effectGrabbed        bool // Whether the selected effect is grabbed for moving
	connected            bool
	midiPort             string
	width                int
	height               int
	sliderWidth          int

	// Navigation state for list-based UI
	navigationMode navigationMode
	effectsList    list.Model
	parameterList  list.Model
	currentSection string // "filter", "overdrive", etc.
	showHelp       bool
	showStatus     bool
	showPagination bool
	showTitle      bool

	// OSC
	client *osc.Client

	// Version
	version string

	// Screen state
	screen     screenState
	prevScreen screenState

	settings config.Settings

	// UI overlays
	showCommandPalette bool
	commandPaletteText string
	showHelpPanel      bool

	// Splash screen
	splashSelection splashOption

	// Current preset tracking
	currentPresetName string
	loadedPresetHash  string // Hash at load time

	// Dirty tracking
	isDirty bool

	// Quit confirmation
	showQuitConfirm bool

	// Preset browser state
	presetBrowser presetBrowserState

	// Settings menu selection
	settingsSelection settingsMenuItem
}

func NewModel(client *osc.Client) Model {
	m := Model{
		// Defaults matching Chroma.sc
		MasterEnabled:        true,
		Gain:                 1.0,
		InputFreezeLength:    0.1,
		FilterEnabled:        true,
		FilterAmount:         0.5,
		FilterCutoff:         2000,
		FilterResonance:      0.3,
		OverdriveEnabled:     false,
		OverdriveDrive:       0.5,
		OverdriveTone:        0.7,
		OverdriveBias:        0.5,
		OverdriveMix:         0.0,
		BitcrushEnabled:      false,
		BitDepth:             8,
		BitcrushSampleRate:   11025,
		BitcrushDrive:        0.5,
		BitcrushMix:          0.3,
		GranularEnabled:      true,
		GranularDensity:      20,
		GranularSize:         0.15,
		GranularPitchScatter: 0.2,
		GranularPosScatter:   0.3,
		GranularMix:          0.5,
		GrainIntensity:       "subtle",
		ReverbEnabled:        false,
		ReverbDecayTime:      3,
		ReverbMix:            0.3,
		DelayEnabled:         false,
		DelayTime:            0.3,
		DelayDecayTime:       3,
		ModRate:              0.5,
		ModDepth:             0.3,
		DelayMix:             0.3,
		DryWet:               0.5,
		EffectsOrder: []string{
			"filter", "overdrive", "bitcrush",
			"granular", "reverb", "delay",
		},

		focused:             ctrlGain,
		selectedEffectIndex: 0,
		connected:           false,
		client:              client,
		navigationMode:      modeEffectsList,
		currentSection:      "master",
		showHelp:            false,
		showStatus:          true,
		showPagination:      true,
		showTitle:           true,
		sliderWidth:         10, // default, will be recalculated on resize
		splashSelection:     splashLast,
	}

	// Load UI settings only
	settings := config.LoadSettings()
	m.settings = settings
	m.showStatus = settings.ShowStatus
	m.showPagination = settings.ShowPagination
	m.showTitle = settings.ShowTitle

	// Start on splash screen
	m.screen = screenSplash
	m.prevScreen = screenSplash

	// If no last settings available, default to "new" selection
	if _, err := config.LoadAutosave(); err != nil {
		if config.LoadLastPresetName() == "" {
			m.splashSelection = splashNew
		}
	}

	return m
}

// SetVersion sets the application version for display.
func (m *Model) SetVersion(v string) {
	m.version = v
}

// switchScreen changes the current screen, saving the previous one.
func (m *Model) switchScreen(s screenState) {
	m.prevScreen = m.screen
	m.screen = s
}

// goBack returns to the previous screen.
func (m *Model) goBack() {
	m.screen = m.prevScreen
}

// panelDimensions calculates the width and height for the effects
// and parameter list panels, accounting for app padding.
func panelDimensions(width, height int) (leftWidth, rightWidth, listHeight int) {
	if width < 60 || height < 20 {
		return 10, 10, 10
	}

	availableWidth := width - 6 // app padding: Padding(2,3) = 4 chars + divider (1 char)

	// Left panel fixed at 20 chars, right panel takes remaining width
	leftWidth = 20
	rightWidth = availableWidth - leftWidth

	// Height calculation: app padding (4) + footer (1) + status bar (1)
	listHeight = height - 6

	// Ensure non-negative dimensions
	if leftWidth < 5 {
		leftWidth = 5
	}
	if rightWidth < 10 {
		rightWidth = 10
	}
	if listHeight < 5 {
		listHeight = 5
	}

	return
}

func (m *Model) InitLists(width, height int) {
	leftWidth, rightWidth, listHeight := panelDimensions(width, height)

	effectsDelegate := list.NewDefaultDelegate()
	m.effectsList = list.New(m.buildEffectsList(), effectsDelegate, leftWidth, listHeight)
	m.effectsList.SetShowTitle(false)
	m.effectsList.SetShowHelp(m.showHelp)
	m.effectsList.SetShowStatusBar(false)
	m.effectsList.SetShowPagination(m.showPagination)

	parameterDelegate := list.NewDefaultDelegate()
	m.parameterList = list.New(nil, parameterDelegate, rightWidth, listHeight)
	m.parameterList.SetShowTitle(false)
	m.parameterList.SetShowHelp(m.showHelp)
	m.parameterList.SetShowStatusBar(false)
	m.parameterList.SetShowPagination(m.showPagination)

	// Initialize parameter list with first effect's parameters
	m.syncParameterPanel()
}

func (m *Model) refreshParameterList() {
	_, rightWidth, _ := panelDimensions(m.width, m.height)

	m.sliderWidth = rightWidth - 24 - 9 - 4
	if m.sliderWidth < 10 {
		m.sliderWidth = 10
	}

	m.parameterList.SetItems(m.buildParameterList(m.currentSection))
}

func (m *Model) refreshEffectsList() {
	if m.effectsList.Items() == nil {
		return
	}
	m.effectsList.SetItems(m.buildEffectsList())
}

func (m *Model) NextControl() {
	m.focused = (m.focused + 1) % ctrlCount
}

func (m *Model) PrevControl() {
	m.focused = (m.focused - 1 + ctrlCount) % ctrlCount
}

func (m *Model) Focused() control {
	return m.focused
}

func (m *Model) IsConnected() bool {
	return m.connected
}

func (m *Model) SetMidiPort(name string) {
	m.midiPort = name
}

func (m *Model) SetConnected(connected bool) {
	m.connected = connected
}

func (m *Model) SetEffectsOrder(order []string) {
	m.EffectsOrder = order
	// Refresh the effects list to reflect new order
	m.refreshEffectsList()
}

func (m *Model) GetEffectsOrder() []string {
	if len(m.EffectsOrder) == 0 {
		// Set default order (master always comes first, not included in reorderable list)
		m.EffectsOrder = []string{
			"filter", "overdrive", "bitcrush",
			"granular", "reverb", "delay",
		}
	}
	return m.EffectsOrder
}

func (m *Model) syncParameterPanel() {
	idx := m.effectsList.Index()
	items := m.effectsList.Items()
	if idx >= 0 && idx < len(items) {
		if eff, ok := items[idx].(effectItem); ok {
			if m.currentSection != eff.id || len(m.parameterList.Items()) == 0 {
				m.currentSection = eff.id
				m.refreshParameterList()
			}
		}
	}
}

func (m *Model) applyPreset(preset config.Preset) {
	// Master
	m.MasterEnabled = preset.MasterEnabled
	m.Gain = preset.Gain
	m.InputFrozen = preset.InputFrozen
	m.InputFreezeLength = preset.InputFreezeLength
	m.DryWet = preset.DryWet
	m.BlendMode = preset.BlendMode
	m.EffectsOrder = preset.EffectsOrder

	// Filter
	m.FilterEnabled = preset.FilterEnabled
	m.FilterAmount = preset.FilterAmount
	m.FilterCutoff = preset.FilterCutoff
	m.FilterResonance = preset.FilterResonance

	// Overdrive
	m.OverdriveEnabled = preset.OverdriveEnabled
	m.OverdriveDrive = preset.OverdriveDrive
	m.OverdriveTone = preset.OverdriveTone
	m.OverdriveBias = preset.OverdriveBias
	m.OverdriveMix = preset.OverdriveMix

	// Bitcrush
	m.BitcrushEnabled = preset.BitcrushEnabled
	m.BitDepth = preset.BitDepth
	m.BitcrushSampleRate = preset.BitcrushSampleRate
	m.BitcrushDrive = preset.BitcrushDrive
	m.BitcrushMix = preset.BitcrushMix

	// Granular
	m.GranularEnabled = preset.GranularEnabled
	m.GranularDensity = preset.GranularDensity
	m.GranularSize = preset.GranularSize
	m.GranularPitchScatter = preset.GranularPitchScatter
	m.GranularPosScatter = preset.GranularPosScatter
	m.GranularMix = preset.GranularMix
	m.GranularFrozen = preset.GranularFrozen
	m.GrainIntensity = preset.GrainIntensity

	// Reverb
	m.ReverbEnabled = preset.ReverbEnabled
	m.ReverbDecayTime = preset.ReverbDecayTime
	m.ReverbMix = preset.ReverbMix

	// Delay
	m.DelayEnabled = preset.DelayEnabled
	m.DelayTime = preset.DelayTime
	m.DelayDecayTime = preset.DelayDecayTime
	m.ModRate = preset.ModRate
	m.ModDepth = preset.ModDepth
	m.DelayMix = preset.DelayMix

	// Refresh UI
	m.refreshEffectsList()
	m.refreshParameterList()

	// Send to SuperCollider
	m.syncAllToOSC()

	// Store hash for dirty detection
	m.loadedPresetHash = preset.Hash()
	m.checkDirty()
}

func (m *Model) buildCurrentPreset() config.Preset {
	return config.Preset{
		MasterEnabled:        m.MasterEnabled,
		Gain:                 m.Gain,
		InputFrozen:          m.InputFrozen,
		InputFreezeLength:    m.InputFreezeLength,
		DryWet:               m.DryWet,
		BlendMode:            m.BlendMode,
		EffectsOrder:         m.EffectsOrder,
		FilterEnabled:        m.FilterEnabled,
		FilterAmount:         m.FilterAmount,
		FilterCutoff:         m.FilterCutoff,
		FilterResonance:      m.FilterResonance,
		OverdriveEnabled:     m.OverdriveEnabled,
		OverdriveDrive:       m.OverdriveDrive,
		OverdriveTone:        m.OverdriveTone,
		OverdriveBias:        m.OverdriveBias,
		OverdriveMix:         m.OverdriveMix,
		BitcrushEnabled:      m.BitcrushEnabled,
		BitDepth:             m.BitDepth,
		BitcrushSampleRate:   m.BitcrushSampleRate,
		BitcrushDrive:        m.BitcrushDrive,
		BitcrushMix:          m.BitcrushMix,
		GranularEnabled:      m.GranularEnabled,
		GranularDensity:      m.GranularDensity,
		GranularSize:         m.GranularSize,
		GranularPitchScatter: m.GranularPitchScatter,
		GranularPosScatter:   m.GranularPosScatter,
		GranularMix:          m.GranularMix,
		GranularFrozen:       m.GranularFrozen,
		GrainIntensity:       m.GrainIntensity,
		ReverbEnabled:        m.ReverbEnabled,
		ReverbDecayTime:      m.ReverbDecayTime,
		ReverbMix:            m.ReverbMix,
		DelayEnabled:         m.DelayEnabled,
		DelayTime:            m.DelayTime,
		DelayDecayTime:       m.DelayDecayTime,
		ModRate:              m.ModRate,
		ModDepth:             m.ModDepth,
		DelayMix:             m.DelayMix,
	}
}

func (m *Model) checkDirty() {
	current := m.buildCurrentPreset()
	m.isDirty = m.loadedPresetHash != current.Hash()
}

func (m *Model) saveUISettings() {
	settings := config.Settings{
		ShowStatus:     m.showStatus,
		ShowPagination: m.showPagination,
		ShowTitle:      m.showTitle,
	}
	config.SaveSettings(settings)
}

func (m *Model) syncAllToOSC() {
	if m.client == nil {
		return
	}
	// Send all values to SuperCollider
	m.client.SetMasterEnabled(m.MasterEnabled)
	m.client.SetGain(m.Gain)
	m.client.SetInputFreeze(m.InputFrozen)
	m.client.SetInputFreezeLength(m.InputFreezeLength)
	m.client.SetDryWet(m.DryWet)
	m.client.SetBlendMode(m.BlendMode)
	m.client.SetEffectsOrder(m.EffectsOrder)

	m.client.SetFilterEnabled(m.FilterEnabled)
	m.client.SetFilterAmount(m.FilterAmount)
	m.client.SetFilterCutoff(m.FilterCutoff)
	m.client.SetFilterResonance(m.FilterResonance)

	m.client.SetOverdriveEnabled(m.OverdriveEnabled)
	m.client.SetOverdriveDrive(m.OverdriveDrive)
	m.client.SetOverdriveTone(m.OverdriveTone)
	m.client.SetOverdriveBias(m.OverdriveBias)
	m.client.SetOverdriveMix(m.OverdriveMix)

	m.client.SetBitcrushEnabled(m.BitcrushEnabled)
	m.client.SetBitDepth(m.BitDepth)
	m.client.SetBitcrushSampleRate(m.BitcrushSampleRate)
	m.client.SetBitcrushDrive(m.BitcrushDrive)
	m.client.SetBitcrushMix(m.BitcrushMix)

	m.client.SetGranularEnabled(m.GranularEnabled)
	m.client.SetGranularDensity(m.GranularDensity)
	m.client.SetGranularSize(m.GranularSize)
	m.client.SetGranularPitchScatter(m.GranularPitchScatter)
	m.client.SetGranularPosScatter(m.GranularPosScatter)
	m.client.SetGranularMix(m.GranularMix)
	m.client.SetGranularFreeze(m.GranularFrozen)
	m.client.SetGrainIntensity(m.GrainIntensity)

	m.client.SetReverbEnabled(m.ReverbEnabled)
	m.client.SetReverbDecayTime(m.ReverbDecayTime)
	m.client.SetReverbMix(m.ReverbMix)

	m.client.SetDelayEnabled(m.DelayEnabled)
	m.client.SetDelayTime(m.DelayTime)
	m.client.SetDelayDecayTime(m.DelayDecayTime)
	m.client.SetModRate(m.ModRate)
	m.client.SetModDepth(m.ModDepth)
	m.client.SetDelayMix(m.DelayMix)
}

func (m *Model) updateFocusedFromSelection() {
	idx := m.parameterList.Index()
	items := m.parameterList.Items()
	if idx >= 0 && idx < len(items) {
		if param, ok := items[idx].(parameterItem); ok {
			m.focused = param.ctrl
		}
	}
}
