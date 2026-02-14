package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-tui/config"
	"github.com/renderorange/chroma/chroma-tui/midi"
	"github.com/renderorange/chroma/chroma-tui/osc"
	"github.com/renderorange/chroma/chroma-tui/tui"
)

func main() {
	scHost := flag.String("host", "127.0.0.1", "SuperCollider host")
	scPort := flag.Int("port", 57120, "SuperCollider OSC port")
	noMidi := flag.Bool("no-midi", false, "Disable MIDI input")
	flag.Parse()

	// Create OSC client
	client := osc.NewClient(*scHost, *scPort)

	// Create TUI model
	model := tui.NewModel(client)

	// Start MIDI handler
	var midiHandler *midi.Handler
	if !*noMidi {
		cfg := config.Load()
		midiHandler = midi.NewHandler(client, cfg)
		if err := midiHandler.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "MIDI warning: %v\n", err)
		} else {
			model.SetMidiPort(midiHandler.PortName())
			defer midiHandler.Stop()
		}
	}

	// Create program
	p := tea.NewProgram(&model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
