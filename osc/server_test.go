package osc

import (
	"testing"

	"github.com/hypebeast/go-osc/osc"
)

func TestNewServer(t *testing.T) {
	s := NewServer(57121)
	if s == nil {
		t.Fatal("expected non-nil server")
	}
	if s.server == nil {
		t.Fatal("expected non-nil internal server")
	}
	if s.stateChan == nil {
		t.Fatal("expected non-nil state channel")
	}
	if cap(s.stateChan) != 10 {
		t.Errorf("expected state channel capacity 10, got %d", cap(s.stateChan))
	}
}

func TestStateChan(t *testing.T) {
	s := NewServer(57125)
	if s.StateChan() == nil {
		t.Fatal("expected non-nil state channel")
	}
}

func TestStateMessageHandler(t *testing.T) {
	s := NewServer(57126)

	// Create a state message with all 35 arguments matching current protocol
	msg := osc.NewMessage("/chroma/state")

	// Arguments in order (matching sendState in Chroma.sc):
	msg.Append(float32(1.0))   // Gain
	msg.Append(int32(0))       // InputFrozen
	msg.Append(float32(0.1))   // InputFreezeLength
	msg.Append(int32(1))       // FilterEnabled
	msg.Append(float32(0.5))   // FilterAmount
	msg.Append(float32(2000))  // FilterCutoff
	msg.Append(float32(0.3))   // FilterResonance
	msg.Append(int32(0))       // OverdriveEnabled
	msg.Append(float32(0.5))   // OverdriveDrive
	msg.Append(float32(0.7))   // OverdriveTone
	msg.Append(float32(0.0))   // OverdriveMix
	msg.Append(int32(1))       // GranularEnabled
	msg.Append(float32(20))    // GranularDensity
	msg.Append(float32(0.15))  // GranularSize
	msg.Append(float32(0.2))   // GranularPitchScatter
	msg.Append(float32(0.3))   // GranularPosScatter
	msg.Append(float32(0.5))   // GranularMix
	msg.Append(int32(0))       // GranularFrozen
	msg.Append("subtle")       // GrainIntensity
	msg.Append(int32(0))       // BitcrushEnabled
	msg.Append(float32(8))     // BitDepth
	msg.Append(float32(11025)) // BitcrushSampleRate
	msg.Append(float32(0.5))   // BitcrushDrive
	msg.Append(float32(0.3))   // BitcrushMix
	msg.Append(int32(0))       // ReverbEnabled
	msg.Append(float32(3))     // ReverbDecayTime
	msg.Append(float32(0.3))   // ReverbMix
	msg.Append(int32(0))       // DelayEnabled
	msg.Append(float32(0.3))   // DelayTime
	msg.Append(float32(3))     // DelayDecayTime
	msg.Append(float32(0.5))   // ModRate
	msg.Append(float32(0.3))   // ModDepth
	msg.Append(float32(0.3))   // DelayMix
	msg.Append(int32(0))       // BlendMode
	msg.Append(float32(0.5))   // DryWet

	// Dispatch the message to trigger the handler
	d := s.server.Dispatcher
	d.Dispatch(msg)

	// Check that the state was updated
	select {
	case state := <-s.stateChan:
		if state.Gain != 1.0 {
			t.Errorf("expected Gain 1.0, got %f", state.Gain)
		}
		if state.FilterEnabled != true {
			t.Error("expected FilterEnabled true")
		}
		if state.OverdriveEnabled != false {
			t.Error("expected OverdriveEnabled false")
		}
		if state.GranularEnabled != true {
			t.Error("expected GranularEnabled true")
		}
		if state.GrainIntensity != "subtle" {
			t.Errorf("expected GrainIntensity 'subtle', got %s", state.GrainIntensity)
		}
		if state.DryWet != 0.5 {
			t.Errorf("expected DryWet 0.5, got %f", state.DryWet)
		}
	default:
		t.Error("expected state update on state channel")
	}
}

func TestStateMessageHandlerTooFewArguments(t *testing.T) {
	s := NewServer(57127)

	// Create a state message with only 30 arguments (should be ignored, need 35)
	msg := osc.NewMessage("/chroma/state")
	for i := 0; i < 30; i++ {
		msg.Append(float32(0.5))
	}

	// Dispatch the message to trigger the handler
	d := s.server.Dispatcher
	d.Dispatch(msg)

	// Check that no state was sent (channel should be empty)
	select {
	case <-s.stateChan:
		t.Error("expected no state update with too few arguments")
	default:
		// This is expected - no state should be sent
	}
}
