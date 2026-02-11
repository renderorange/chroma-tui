package osc

import (
	"fmt"
	"sync"

	"github.com/hypebeast/go-osc/osc"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
	Code    int
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// Error codes
const (
	ErrCodeInvalidType   = 1001
	ErrCodeOutOfRange    = 1002
	ErrCodeMissingField  = 1003
	ErrCodeProtocolError = 1004
)

type State struct {
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
	OverdriveMix         float32
	GranularEnabled      bool
	GranularDensity      float32
	GranularSize         float32
	GranularPitchScatter float32
	GranularPosScatter   float32
	GranularMix          float32
	GranularFrozen       bool
	GrainIntensity       string
	BitcrushEnabled      bool
	BitDepth             float32
	BitcrushSampleRate   float32
	BitcrushDrive        float32
	BitcrushMix          float32
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
	// Effects order
	EffectsOrder []string
}

type Server struct {
	server       *osc.Server
	stateChan    chan State
	stateMu      sync.RWMutex
	currentState State
}

func NewServer(port int) *Server {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	s := &Server{
		stateChan: make(chan State, 10), // Buffer for smooth 30fps updates
	}

	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/chroma/state", func(msg *osc.Message) {
		if len(msg.Arguments) >= 35 {
			state := State{
				Gain:                 toFloat32(msg.Arguments[0]),
				InputFrozen:          toInt(msg.Arguments[1]) == 1,
				InputFreezeLength:    toFloat32(msg.Arguments[2]),
				FilterEnabled:        toInt(msg.Arguments[3]) == 1,
				FilterAmount:         toFloat32(msg.Arguments[4]),
				FilterCutoff:         toFloat32(msg.Arguments[5]),
				FilterResonance:      toFloat32(msg.Arguments[6]),
				OverdriveEnabled:     toInt(msg.Arguments[7]) == 1,
				OverdriveDrive:       toFloat32(msg.Arguments[8]),
				OverdriveTone:        toFloat32(msg.Arguments[9]),
				OverdriveMix:         toFloat32(msg.Arguments[10]),
				GranularEnabled:      toInt(msg.Arguments[11]) == 1,
				GranularDensity:      toFloat32(msg.Arguments[12]),
				GranularSize:         toFloat32(msg.Arguments[13]),
				GranularPitchScatter: toFloat32(msg.Arguments[14]),
				GranularPosScatter:   toFloat32(msg.Arguments[15]),
				GranularMix:          toFloat32(msg.Arguments[16]),
				GranularFrozen:       toInt(msg.Arguments[17]) == 1,
				GrainIntensity:       toString(msg.Arguments[18]),
				BitcrushEnabled:      toInt(msg.Arguments[19]) == 1,
				BitDepth:             toFloat32(msg.Arguments[20]),
				BitcrushSampleRate:   toFloat32(msg.Arguments[21]),
				BitcrushDrive:        toFloat32(msg.Arguments[22]),
				BitcrushMix:          toFloat32(msg.Arguments[23]),
				ReverbEnabled:        toInt(msg.Arguments[24]) == 1,
				ReverbDecayTime:      toFloat32(msg.Arguments[25]),
				ReverbMix:            toFloat32(msg.Arguments[26]),
				DelayEnabled:         toInt(msg.Arguments[27]) == 1,
				DelayTime:            toFloat32(msg.Arguments[28]),
				DelayDecayTime:       toFloat32(msg.Arguments[29]),
				ModRate:              toFloat32(msg.Arguments[30]),
				ModDepth:             toFloat32(msg.Arguments[31]),
				DelayMix:             toFloat32(msg.Arguments[32]),
				BlendMode:            toInt(msg.Arguments[33]),
				DryWet:               toFloat32(msg.Arguments[34]),
			}
			s.stateMu.Lock()
			s.currentState = state
			s.stateMu.Unlock()
			// Non-blocking send
			select {
			case s.stateChan <- state:
			default:
			}
		}
	})

	// Effects order response handler
	d.AddMsgHandler("/chroma/effectsOrder", func(msg *osc.Message) {
		if len(msg.Arguments) > 0 {
			order := make([]string, len(msg.Arguments))
			for i, arg := range msg.Arguments {
				if str, ok := arg.(string); ok {
					order[i] = str
				}
			}
			s.stateMu.Lock()
			s.currentState.EffectsOrder = order
			state := s.currentState
			s.stateMu.Unlock()
			// Send to channel for real-time updates
			select {
			case s.stateChan <- state:
			default:
				// Channel full, skip update
			}
		}
	})

	s.server = &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) StateChan() <-chan State {
	return s.stateChan
}

func toFloat32(v interface{}) float32 {
	switch val := v.(type) {
	case float32:
		return val // Remove validation - trust SuperCollider
	case float64:
		return float32(val) // Remove validation
	case int32:
		return float32(val) // Remove validation
	case int:
		return float32(val) // Remove validation
	default:
		fmt.Printf("WARNING: toFloat32 invalid type %T, expected float32/float64/int32/int, using 0.0\n", v)
		return 0.0
	}
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int32:
		return int(val) // Remove validation
	case int:
		return val // Remove validation
	case float32:
		return int(val) // Remove validation
	case float64:
		return int(val) // Remove validation
	default:
		fmt.Printf("WARNING: toInt invalid type %T, expected int/int32/float32/float64, using 0\n", v)
		return 0
	}
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		if len(val) > 255 {
			fmt.Printf("WARNING: toString string length %d exceeds maximum 255, truncating\n", len(val))
			return val[:255]
		}
		return val
	default:
		fmt.Printf("WARNING: toString invalid type %T, expected string, using 'subtle'\n", v)
		return "subtle"
	}
}
