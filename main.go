package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderorange/chroma/chroma-control/config"
	"github.com/renderorange/chroma/chroma-control/midi"
	"github.com/renderorange/chroma/chroma-control/osc"
	"github.com/renderorange/chroma/chroma-control/tui"
)

var version = "dev"

func main() {
	scHost := flag.String("host", "127.0.0.1", "SuperCollider host")
	scPort := flag.Int("port", 57120, "SuperCollider OSC port")
	noMidi := flag.Bool("no-midi", false, "Disable MIDI input")
	flag.Parse()

	// Create OSC client
	client := osc.NewClient(*scHost, *scPort)

	// Create TUI model
	model := tui.NewModel(client)
	model.SetVersion(version)

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
