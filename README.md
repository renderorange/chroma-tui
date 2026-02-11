# Chroma TUI

Terminal-based interface for Chroma audio effects system.

This repository contains the Go-based terminal user interface (TUI) component of the Chroma audio effects system. The TUI communicates with an external Chroma SuperCollider audio engine via OSC (Open Sound Control).

## External Dependency

This TUI requires the Chroma SuperCollider audio engine to function. Please install the Chroma SuperCollider engine first:
- **Repository**: [Chroma](https://github.com/renderorange/chroma)
- **Setup**: Follow installation instructions in Chroma repository

The TUI will connect to the running Chroma engine via OSC protocol (localhost:57120/9000).

## Quick Start

```bash
# Clone this repository
git clone https://github.com/renderorange/chroma-tui.git
cd chroma-tui

# Build
go build -o chroma-tui

# Run (requires Chroma SuperCollider engine running)
./chroma-tui
```

## Requirements

### System Requirements
- **Go 1.24+**: For building the TUI
- **Terminal**: With proper Unicode and color support
- **Chroma SuperCollider**: External audio engine (required)

### Optional Hardware
- **MIDI Controller**: For hardware parameter control
- **Audio Interface**: Managed by external Chroma engine

## Controls

### Keyboard Navigation
- `↑/↓` - Navigate parameters
- `←/→` - Adjust parameter values
- `Tab` - Switch between effects
- `Space` - Toggle freeze/granular
- `q` - Quit application

### MIDI Mapping
Default MIDI mappings can be customized in `config/midi.toml`:
- CC 1: Master Gain
- CC 2: Filter Cutoff
- CC 3: Filter Resonance
- And more...

## Development

### Building from Source
```bash
git clone https://github.com/renderorange/chroma-tui.git
cd chroma-tui
go mod tidy
go build -o chroma-tui
```

### Running Tests
```bash
go test ./...
go test -v ./...
go test -cover ./...
```

### Cross-Platform Builds
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o chroma-tui-linux

# macOS  
GOOS=darwin GOARCH=amd64 go build -o chroma-tui-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o chroma-tui.exe
```

## Architecture

This repository contains only the TUI component:

```
┌─────────────────┐    OSC Protocol    ┌────────────────────┐
│   TUI Repo    │ ←───────────────→ │ Chroma Engine    │
│ (This Repo)   │                  │ (External)       │
│               │                  │                  │
│ • Interface    │                  │ • Audio Engine   │
│ • Controls    │                  │ • Effects        │
│ • Parameters  │                  │ • Spectrum       │
└─────────────────┘                  └────────────────────┘
```

### Package Structure
- `tui/` - Terminal UI components and state management
- `osc/` - OSC communication with external engine  
- `midi/` - Hardware MIDI controller support
- `config/` - Configuration management and user customization

## Platform Support

| Platform | MIDI Support | Build Status |
|----------|--------------|--------------|
| Linux    | Native ALSA  | Tested     |
| macOS    | CoreAudio     | Tested     |  
| Windows  | Win32 API    | Tested     |
| Other    | No-CGO Fallback | Tested     |

---

**Note**: This is the TUI component only. For the complete Chroma system including the audio engine, see the main [Chroma repository](https://github.com/renderorange/chroma).