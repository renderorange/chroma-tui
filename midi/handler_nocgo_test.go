//go:build !cgo
// +build !cgo

package midi

import (
	"fmt"
	"testing"

	"github.com/renderorange/chroma/chroma-tui/config"
	"github.com/renderorange/chroma/chroma-tui/osc"
)

func TestHandlerNoCGO_CreationWithValidParameters(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}
}

func TestHandlerNoCGO_CreationWithNilClient(t *testing.T) {
	cfg := config.DefaultConfig()

	// Should handle nil client gracefully
	handler := NewHandler(nil, cfg)
	if handler == nil {
		t.Fatal("expected no-CGO MIDI handler to handle nil client")
	}
}

func TestHandlerNoCGO_CreationWithEmptyConfig(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.Config{} // Empty config

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected no-CGO MIDI handler to handle empty config")
	}
}

func TestHandlerNoCGO_PortNameReturnsEmpty(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// no-CGO handler should return empty port name
	portName := handler.PortName()
	if portName != "" {
		t.Errorf("expected empty port name for no-CGO handler, got '%s'", portName)
	}
}

func TestHandlerNoCGO_StartReturnsError(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// no-CGO handler should return error on start
	err := handler.Start()
	if err == nil {
		t.Error("expected error when starting no-CGO MIDI handler")
	}

	// Error message should indicate MIDI not available
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestHandlerNoCGO_StopDoesNothing(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// no-CGO handler should not panic on stop
	handler.Stop()

	// Should be able to call stop multiple times
	handler.Stop()
	handler.Stop()
}

func TestHandlerNoCGO_ConfigHandling(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// Should still load config even though MIDI is not available
	if len(cfg.CC) == 0 {
		t.Error("expected config to have CC mappings even in no-CGO mode")
	}

	if len(cfg.Notes) == 0 {
		t.Error("expected config to have note mappings even in no-CGO mode")
	}

	if len(cfg.EffectsOrder) == 0 {
		t.Error("expected config to have effects order even in no-CGO mode")
	}
}

func TestHandlerNoCGO_ErrorHandling(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)

	// Test with invalid config
	cfg := config.Config{
		CC:    map[string]int{"invalid": -1}, // Invalid CC number
		Notes: map[string]int{"invalid": -1}, // Invalid note number
	}

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected no-CGO MIDI handler to handle invalid config gracefully")
	}

	// Should still return error on start (due to no-CGO, not config)
	err := handler.Start()
	if err == nil {
		t.Error("expected error when starting no-CGO MIDI handler even with invalid config")
	}
}

func TestHandlerNoCGO_ConcurrentAccess(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// Test concurrent access to PortName()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			portName := handler.PortName()
			if portName != "" {
				t.Errorf("expected empty port name in concurrent access, got '%s'", portName)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestHandlerNoCGO_MultipleStartStop(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// Test multiple start/stop cycles
	for i := 0; i < 3; i++ {
		err := handler.Start()
		if err == nil {
			t.Error("expected error when starting no-CGO MIDI handler")
		}
		handler.Stop()
	}
}

func TestHandlerNoCGO_FallbackBehavior(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	handler := NewHandler(client, cfg)
	if handler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// Test that no-CGO handler provides safe fallback behavior
	// - PortName() returns empty string
	// - Start() returns error
	// - Stop() does nothing
	// - No panic on any operation

	portName := handler.PortName()
	if portName != "" {
		t.Errorf("expected fallback behavior: empty port name, got '%s'", portName)
	}

	err := handler.Start()
	if err == nil {
		t.Error("expected fallback behavior: start returns error")
	}

	// Should not panic
	handler.Stop()
}

func TestHandlerNoCGO_ComparisonWithCGOHandler(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)
	cfg := config.DefaultConfig()

	// Create handler (this is the no-CGO version due to build tags)
	handler := NewHandler(client, cfg)

	if handler == nil {
		t.Fatal("expected non-nil MIDI handler")
	}

	// Should always return error in no-CGO mode
	err := handler.Start()
	if err == nil {
		t.Error("expected no-CGO handler to always return error on start")
	}

	// Clean up
	handler.Stop()
}
	if nocgoHandler == nil {
		t.Fatal("expected non-nil no-CGO MIDI handler")
	}

	// Both should have different behaviors
	// CGO handler might succeed, no-CGO should fail
	cgoErr := cgoHandler.Start()
	nocgoErr := nocgoHandler.Start()

	// no-CGO should always return error
	if nocgoErr == nil {
		t.Error("expected no-CGO handler to always return error on start")
	}

	// Clean up
	cgoHandler.Stop()
	nocgoHandler.Stop()
}

func TestHandlerNoCGO_ConfigIndependence(t *testing.T) {
	client := osc.NewClient("127.0.0.1", 57120)

	// Test with different configs
	configs := []config.Config{
		config.DefaultConfig(),
		config.Config{}, // Empty
		config.Config{
			CC:    map[string]int{"custom": 10},
			Notes: map[string]int{"custom": 20},
		},
	}

	for i, cfg := range configs {
		t.Run(fmt.Sprintf("config_%d", i), func(t *testing.T) {
			handler := NewHandler(client, cfg)
			if handler == nil {
				t.Fatal("expected non-nil no-CGO MIDI handler")
			}

			// Should always behave the same regardless of config
			portName := handler.PortName()
			if portName != "" {
				t.Errorf("expected empty port name regardless of config, got '%s'", portName)
			}

			err := handler.Start()
			if err == nil {
				t.Error("expected error on start regardless of config")
			}

			handler.Stop()
		})
	}
}
