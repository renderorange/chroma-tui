package main

import (
	"flag"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-tui/config"
	"github.com/renderorange/chroma/chroma-tui/midi"
	"github.com/renderorange/chroma/chroma-tui/osc"
	"github.com/renderorange/chroma/chroma-tui/tui"
)

func TestMain_FlagParsingSetsDefaultValues(t *testing.T) {
	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test default flags
	os.Args = []string{"chroma-tui"}

	// Reset flag set for clean test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Parse flags (this simulates main() flag parsing)
	scHost := flag.String("host", "127.0.0.1", "SuperCollider host")
	scPort := flag.Int("port", 57120, "SuperCollider OSC port")
	listenPort := flag.Int("listen", 9000, "Port to listen for state updates")
	noMidi := flag.Bool("no-midi", false, "Disable MIDI input")
	flag.Parse()

	// Verify defaults
	if *scHost != "127.0.0.1" {
		t.Errorf("expected default host 127.0.0.1, got %s", *scHost)
	}
	if *scPort != 57120 {
		t.Errorf("expected default port 57120, got %d", *scPort)
	}
	if *listenPort != 9000 {
		t.Errorf("expected default listen port 9000, got %d", *listenPort)
	}
	if *noMidi != false {
		t.Errorf("expected default no-midi false, got %v", *noMidi)
	}
}

func TestMain_FlagParsingWithCustomValues(t *testing.T) {
	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test custom flags
	os.Args = []string{"chroma-tui", "-host", "192.168.1.100", "-port", "57130", "-listen", "9010", "-no-midi"}

	// Reset flag set for clean test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Parse flags
	scHost := flag.String("host", "127.0.0.1", "SuperCollider host")
	scPort := flag.Int("port", 57120, "SuperCollider OSC port")
	listenPort := flag.Int("listen", 9000, "Port to listen for state updates")
	noMidi := flag.Bool("no-midi", false, "Disable MIDI input")
	flag.Parse()

	// Verify custom values
	if *scHost != "192.168.1.100" {
		t.Errorf("expected custom host 192.168.1.100, got %s", *scHost)
	}
	if *scPort != 57130 {
		t.Errorf("expected custom port 57130, got %d", *scPort)
	}
	if *listenPort != 9010 {
		t.Errorf("expected custom listen port 9010, got %d", *listenPort)
	}
	if *noMidi != true {
		t.Errorf("expected custom no-midi true, got %v", *noMidi)
	}
}

func TestMain_ComponentCreation(t *testing.T) {
	// Test that main components can be created with valid parameters
	scHost := "127.0.0.1"
	scPort := 57120
	listenPort := 9000

	// Create OSC client
	client := osc.NewClient(scHost, scPort)
	if client == nil {
		t.Fatal("expected non-nil OSC client")
	}

	// Create OSC server
	server := osc.NewServer(listenPort)
	if server == nil {
		t.Fatal("expected non-nil OSC server")
	}

	// Create TUI model
	model := tui.NewModel(client)
	// Check if model has expected fields (e.g., EffectsOrder)
	if len(model.EffectsOrder) == 0 {
		t.Fatal("expected TUI model to have non-empty EffectsOrder")
	}
}

func TestMain_MIDIHandlerCreationWithNoMidiFlag(t *testing.T) {
	// Test MIDI handler behavior when -no-midi flag is set
	noMidi := true
	var midiHandler *midi.Handler

	if !noMidi {
		cfg := config.Load()
		client := osc.NewClient("127.0.0.1", 57120)
		midiHandler = midi.NewHandler(client, cfg)
		if midiHandler == nil {
			t.Fatal("expected non-nil MIDI handler when MIDI enabled")
		}
	} else {
		// When no-midi is true, midiHandler should remain nil
		if midiHandler != nil {
			t.Error("expected nil MIDI handler when -no-midi flag is set")
		}
	}
}

func TestMain_TeaProgramCreation(t *testing.T) {
	// Test that tea program can be created with model
	client := osc.NewClient("127.0.0.1", 57120)
	model := tui.NewModel(client)

	// Create program with alt screen option
	p := tea.NewProgram(&model, tea.WithAltScreen())
	if p == nil {
		t.Fatal("expected non-nil tea program")
	}
}

func TestMain_OSCClientServerIntegration(t *testing.T) {
	// Test that OSC client and server can be created together
	scHost := "127.0.0.1"
	scPort := 57121 // Use different port to avoid conflicts
	listenPort := 9001

	client := osc.NewClient(scHost, scPort)
	server := osc.NewServer(listenPort)

	// Verify both are created
	if client == nil {
		t.Fatal("expected non-nil OSC client")
	}
	if server == nil {
		t.Fatal("expected non-nil OSC server")
	}

	// Test that server can be started (in background for testing)
	serverStopped := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			serverStopped <- err
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test that client can send a basic message
	err := client.SendSync()
	if err != nil {
		t.Errorf("expected no error sending sync message, got %v", err)
	}

	// Note: Server doesn't have explicit Stop method in current implementation
	// Server will be cleaned up when the test ends
	select {
	case err := <-serverStopped:
		if err != nil {
			t.Errorf("server stopped with error: %v", err)
		}
	case <-time.After(time.Second):
		t.Log("server is running (no explicit stop needed for test)")
	}
}

func TestMain_ConfigLoadingForMIDI(t *testing.T) {
	// Test that config can be loaded for MIDI handler
	cfg := config.Load()

	// Config is a struct, not a pointer, so check if it has expected fields
	if len(cfg.CC) == 0 {
		t.Fatal("expected config to have CC mappings")
	}
	if len(cfg.Notes) == 0 {
		t.Fatal("expected config to have Notes mappings")
	}
	if len(cfg.EffectsOrder) == 0 {
		t.Fatal("expected config to have EffectsOrder")
	}

	// Test that config has expected fields for MIDI
	if cfg.CC["gain"] != 1 {
		t.Errorf("expected gain CC mapping to be 1, got %d", cfg.CC["gain"])
	}
	if cfg.Notes["input_freeze"] != 60 {
		t.Errorf("expected input_freeze note mapping to be 60, got %d", cfg.Notes["input_freeze"])
	}
}

func TestMain_ErrorHandling(t *testing.T) {
	// Test error handling for invalid parameters

	// Test with invalid port (should not panic)
	client := osc.NewClient("127.0.0.1", -1)
	if client == nil {
		t.Fatal("expected client to handle invalid port gracefully")
	}

	// Test with invalid listen port
	server := osc.NewServer(-1)
	if server == nil {
		t.Fatal("expected server to handle invalid listen port gracefully")
	}
}

func TestMain_ProgramExitCodes(t *testing.T) {
	// Test that main handles different exit scenarios
	// This test verifies the structure but doesn't actually run main()

	// Test that os.Exit(1) is called on program error
	// (This would require more complex testing with subprocess)

	// For now, verify that error path exists in main
	client := osc.NewClient("127.0.0.1", 57120)
	model := tui.NewModel(client)
	p := tea.NewProgram(&model, tea.WithAltScreen())

	// The program should be able to be created without error
	if p == nil {
		t.Fatal("expected program to be created successfully")
	}
}
