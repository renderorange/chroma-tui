package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

type control int

type navigationMode int

const (
	modeEffectsList navigationMode = iota
	modeParameterList
)

const (
	ctrlGain control = iota
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
	focused             control
	selectedEffectIndex int // Index of currently selected effect for reordering
	connected           bool
	midiPort            string
	width               int
	height              int

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
}

func NewModel(client *osc.Client) Model {
	return Model{
		// Defaults matching Chroma.sc
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
		currentSection:      "input",
		showHelp:            true,
		showStatus:          true,
		showPagination:      true,
		showTitle:           true,
	}
}

func (m *Model) InitLists(width, height int) {
	effectsDelegate := list.NewDefaultDelegate()

	m.effectsList = list.New(m.buildEffectsList(), effectsDelegate, width/2, height-4)
	m.effectsList.Title = "Effects"
	m.effectsList.SetShowHelp(m.showHelp)
	m.effectsList.SetShowStatusBar(m.showStatus)
	m.effectsList.SetShowPagination(m.showPagination)
	m.effectsList.SetShowTitle(m.showTitle)

	parameterDelegate := list.NewDefaultDelegate()
	m.parameterList = list.New(nil, parameterDelegate, width/2, height-4)
	m.parameterList.SetShowHelp(m.showHelp)
	m.parameterList.SetShowStatusBar(m.showStatus)
	m.parameterList.SetShowPagination(m.showPagination)
	m.parameterList.SetShowTitle(m.showTitle)

	m.refreshParameterList()
}

func (m *Model) refreshParameterList() {
	m.parameterList.SetItems(m.buildParameterList(m.currentSection))
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
	// Trigger OSC update
	// Will be connected to OSC client in next task
}

func (m *Model) GetEffectsOrder() []string {
	if len(m.EffectsOrder) == 0 {
		// Set default order
		m.EffectsOrder = []string{
			"filter", "overdrive", "bitcrush",
			"granular", "reverb", "delay",
		}
	}
	return m.EffectsOrder
}
