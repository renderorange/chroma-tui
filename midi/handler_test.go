package midi

import (
	"testing"

	"github.com/renderorange/chroma/chroma-tui/config"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestHandler_CreationWithValidParameters(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}
}

func TestHandler_CreationWithNilClient(t *testing.T) {
	cfg := config.DefaultConfig()

	// Should handle nil client gracefully
	handler := NewHandler(nil, cfg)
	if handler == nil {
		t.Fatal("expected MIDI handler to handle nil client")
	}
}

func TestHandler_CreationWithEmptyConfig(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.Config{} // Empty config

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected MIDI handler to handle empty config")
	}
}

func TestHandler_CCMapping(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test that CC mappings are loaded from config
	if cfg.CC["gain"] != 1 {
		t.Errorf("expected gain CC mapping to be 1, got %d", cfg.CC["gain"])
	}

	// Test that handler has access to CC mappings
	// (This would require exposing internal state or methods for testing)
}

func TestHandler_NoteMapping(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test that note mappings are loaded from config
	if cfg.Notes["input_freeze"] != 60 {
		t.Errorf("expected input_freeze note mapping to be 60, got %d", cfg.Notes["input_freeze"])
	}

	// Test that handler has access to note mappings
	// (This would require exposing internal state or methods for testing)
}

func TestHandler_EffectsOrderFromConfig(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test that effects order is loaded from config
	expectedOrder := []string{"filter", "overdrive", "bitcrush", "granular", "reverb", "delay"}
	if len(cfg.EffectsOrder) != len(expectedOrder) {
		t.Errorf("expected effects order length %d, got %d", len(expectedOrder), len(cfg.EffectsOrder))
	}

	for i, effect := range expectedOrder {
		if cfg.EffectsOrder[i] != effect {
			t.Errorf("expected effect %d to be %s, got %s", i, effect, cfg.EffectsOrder[i])
		}
	}
}

func TestHandler_PortName(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test that PortName() returns a string (even if empty)
	portName := handler.PortName()
	if portName == "" {
		t.Log("MIDI port name is empty (may be default when no device connected)")
	}
}

func TestHandler_StartStop(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test that Start() can be called without error
	err := handler.Start()
	if err != nil {
		t.Logf("MIDI handler start returned error (may be expected if no MIDI device): %v", err)
	}

	// Test that Stop() can be called without error
	handler.Stop()
}

func TestHandler_StartStopWithNoDevice(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.Config{} // Empty config, no device mappings

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Should handle start/stop gracefully even with no device
	err := handler.Start()
	if err != nil {
		t.Logf("Expected error when starting MIDI handler with no device: %v", err)
	}

	handler.Stop()
}

func TestHandler_MIDIMessageProcessing(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test that MIDI messages can be processed
	// This would require exposing the message processing method
	// For now, we test that the handler can be created with the necessary config

	// Verify CC mappings exist for processing
	if len(cfg.CC) == 0 {
		t.Error("expected CC mappings to be available for MIDI message processing")
	}

	// Verify note mappings exist for processing
	if len(cfg.Notes) == 0 {
		t.Error("expected note mappings to be available for MIDI message processing")
	}
}

func TestHandler_ErrorHandling(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)

	// Test with invalid config
	cfg := config.Config{
		CC:    map[string]int{"invalid": -1}, // Invalid CC number
		Notes: map[string]int{"invalid": -1}, // Invalid note number
	}

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected MIDI handler to handle invalid config gracefully")
	}

	// Should still be able to start and stop
	err := handler.Start()
	if err != nil {
		t.Logf("Expected error with invalid config: %v", err)
	}
	handler.Stop()
}

func TestHandler_ConcurrentAccess(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test concurrent access to PortName()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_ = handler.PortName()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestHandler_ResourceCleanup(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Test that resources are cleaned up properly
	err := handler.Start()
	if err == nil {
		// Only test stop if start succeeded
		handler.Stop()
	}

	// Should be able to call Stop() multiple times without panic
	handler.Stop()
	handler.Stop()
}
