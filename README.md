# Chroma TUI

Terminal-based interface for Chroma audio effects system.

This repository contains the Go-based terminal user interface (TUI) component of the Chroma audio effects system. The TUI communicates with an external Chroma SuperCollider audio engine via OSC (Open Sound Control).

## Features

### Audio Effects Chain
- **Filter**: Resonant low-pass filter with cutoff, resonance, and amount control
- **Overdrive**: Distortion with drive, tone, bias (-1.0 to 1.0), and mix control
- **Bitcrusher**: Bit depth reduction (4-16 bits), sample rate reduction, drive, and mix
- **Granular**: Granular synthesis with density, size, pitch/position scatter, and intensity modes
- **Reverb**: Algorithmic reverb with decay time and mix control
- **Delay**: Modulated delay with time, decay, rate/depth modulation, and mix

### Real-time Control System
- **Effects Reordering**: Visual chain reordering with PgUp/PgDn controls
- **Three-State Intensity**: Granular intensity (subtle/pronounced/extreme) and blend modes

### MIDI Integration
- **11 CC Mappings**: Hardware controller support
- **5 Note Mappings**: Freeze controls and blend mode switching
- **Customizable**: TOML-based configuration for custom mappings
- **Auto-Detect**: Automatic MIDI device discovery and connection

### Effects Reordering System

The TUI provides real-time effects chain reordering:

#### Controls
- **Navigation**: Tab/↑/↓/j/k to select effect or parameter
- **Panel Switching**: Enter to open parameter panel, Esc to return to effects list
- **Reordering**: PageUp to move effect up, PageDown to move down
- **Reset**: 'r' key to restore default order
- **Persistence**: Order saved and restored across sessions

#### OSC Integration
Effects order changes are transmitted via `/chroma/effectsOrder` endpoint.

#### Default Processing Order
```
Input → Filter → Overdrive → Bitcrusher → Granular → Reverb → Delay → Output
```

## External Dependency

This TUI requires the Chroma SuperCollider audio engine to function. Please install the Chroma SuperCollider engine first:
- **Repository**: [Chroma](https://github.com/renderorange/chroma)
- **Setup**: Follow installation instructions in Chroma repository

The TUI will connect to the running Chroma engine via OSC protocol (localhost:57120).

## Architecture

Chroma-TUI communicates with Chroma via **stateless OSC over UDP**. This means:

- The TUI sends control messages to Chroma (fire-and-forget)
- Chroma processes audio independently and maintains its own state
- There is no two-way synchronization of parameters
- The TUI's displayed values reflect what was last sent, not necessarily Chroma's current state
- Multiple TUI instances can control the same Chroma instance
- If Chroma restarts, it will use default values until the TUI sends new commands

This design provides low latency and simplicity but means the TUI does not receive confirmation that commands were received.

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
- **Terminal**: With proper Unicode and color support (minimum 60x20 size)
- **Chroma SuperCollider**: External audio engine (required)

## Controls

### Keyboard Navigation
| Navigation | Action |
|------------|--------|
| `Tab`, `↓`, `j` | Next parameter/effect |
| `Shift+Tab`, `↑`, `k` | Previous parameter/effect |
| `Enter` | Open parameter panel (from effects list) |
| `Esc` | Return to effects list (from parameters) |
| `←`, `h` | Decrease parameter value |
| `→`, `l` | Increase parameter value |
| `Space` | Toggle boolean parameters |

### Special Controls
| Key | Action |
|-----|--------|
| `i` | Cycle grain intensity (subtle → pronounced → extreme) |
| `1`, `2`, `3` | Set blend mode (Mirror, Complement, Transform) |
| `PageUp` | Move selected effect up in chain (global section) |
| `PageDown` | Move selected effect down in chain (global section) |
| `r` | Reset effects order to default (global section) |
| `q`, `Ctrl+C` | Quit application |

### Parameter Reference

#### Input & Gain
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Gain | 0.0 - 2.0 | 1.0 | Linear |
| Input Freeze Length | 0.05 - 0.5s | 0.1s | Linear |
| Input Freeze | Boolean | false | Toggle |

#### Filter
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Enabled | Boolean | true | Toggle |
| Amount | 0.0 - 1.0 | 0.5 | Linear |
| Cutoff | 200 - 8000 Hz | 2000 | Linear |
| Resonance | 0.0 - 1.0 | 0.3 | Linear |

#### Overdrive Effect
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Enabled | Boolean | false | Toggle |
| Drive | 0.0 - 1.0 | 0.5 | Linear |
| Tone | 0.0 - 1.0 | 0.7 | Linear |
| Bias | -1.0 - 1.0 | 0.5 | Linear |
| Mix | 0.0 - 1.0 | 0.0 | Linear |

#### Bitcrusher Effect
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Enabled | Boolean | false | Toggle |
| Bit Depth | 4 - 16 bits | 8 | Linear |
| Sample Rate | 1000 - 44100 Hz | 11025 | Linear |
| Drive | 0.0 - 1.0 | 0.5 | Linear |
| Mix | 0.0 - 1.0 | 0.3 | Linear |

#### Granular Effect
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Enabled | Boolean | true | Toggle |
| Density | 1 - 50 grains | 20 | Logarithmic |
| Size | 0.01 - 0.5s | 0.15s | Logarithmic |
| Pitch Scatter | 0.0 - 1.0 | 0.2 | Linear |
| Position Scatter | 0.0 - 1.0 | 0.3 | Linear |
| Mix | 0.0 - 1.0 | 0.5 | Linear |
| Freeze | Boolean | false | Toggle |
| Grain Intensity | Subtle/Pronounced/Extreme | Subtle | Three-state |

#### Reverb Effect
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Enabled | Boolean | false | Toggle |
| Decay Time | 0.5 - 10.0s | 3.0s | Linear |
| Mix | 0.0 - 1.0 | 0.3 | Linear |

#### Delay Effect
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Enabled | Boolean | false | Toggle |
| Time | 0.1 - 1.0s | 0.3s | Linear |
| Decay Time | 0.5 - 10.0s | 3.0s | Linear |
| Modulation Rate | 0.1 - 5.0 Hz | 0.5 Hz | Logarithmic |
| Modulation Depth | 0.0 - 1.0 | 0.3 | Logarithmic |
| Mix | 0.0 - 1.0 | 0.3 | Linear |

#### Master Controls
| Parameter | Range | Default | Control |
|-----------|-------|---------|---------|
| Blend Mode | Mirror/Complement/Transform | Mirror | Three-state |
| Dry/Wet | 0.0 - 1.0 | 0.5 | Linear |
| Effects Order | Customizable | filter→overdrive→bitcrush→granular→reverb→delay | Reorderable (PageUp/PageDown) |

### Special Parameter Behaviors

#### Logarithmic Scaling
- **Granular Density**: 1-50 grains (logarithmic)
- **Granular Size**: 0.01-0.5s (logarithmic)
- **Modulation Rate**: 0.1-5.0 Hz (logarithmic)
- **Modulation Depth**: 0.0-1.0 (logarithmic)

#### Three-State Values
- **GrainIntensity**: Cycle with `i` key → subtle → pronounced → extreme → subtle
- **BlendMode**: `1`/`2`/`3` keys or MIDI notes → Mirror (0) → Complement (1) → Transform (2)

#### Dynamic Interface
- **Slider Width**: Automatically scales based on terminal width (minimum 10 characters)
- **Panel Ratio**: 25/75 split between effects list and parameters
- **Context Footer**: Shows relevant keybindings for current mode (effects list, parameters, or global section)

### MIDI Integration

#### Default Control Change (CC) Mappings
| CC # | Parameter | Range | Value Mapping |
|------|-----------|-------|---------------|
| 1 | Gain | 0.0-2.0 | 0-127 → 0.0-2.0 |
| 2 | Input Freeze Length | 0.05-0.5s | 0-127 → 0.05-0.5s |
| 3 | Filter Amount | 0.0-1.0 | 0-127 → 0.0-1.0 |
| 4 | Filter Cutoff | 200-8000 Hz | 0-127 → 200-8000 |
| 5 | Filter Resonance | 0.0-1.0 | 0-127 → 0.0-1.0 |
| 6 | Granular Density | 1-50 | 0-127 → 1-50 (log) |
| 7 | Granular Size | 0.01-0.5s | 0-127 → 0.01-0.5s (log) |
| 8 | Granular Mix | 0.0-1.0 | 0-127 → 0.0-1.0 |
| 9 | Reverb Mix & Delay Mix | 0.0-1.0 | 0-127 → 0.0-1.0 |
| 10 | Decay Time | 0.5-10.0s | 0-127 → 0.5-10.0s |
| 11 | Dry/Wet | 0.0-1.0 | 0-127 → 0.0-1.0 |

#### Note Mappings
| Note | Action |
|------|--------|
| C4 (60) | Toggle Input Freeze |
| D4 (62) | Toggle Granular Freeze |
| E4 (64) | Set Blend Mode to Mirror |
| F4 (65) | Set Blend Mode to Complement |
| G4 (67) | Set Blend Mode to Transform |

## Configuration

### MIDI Configuration

Create `~/.config/chroma/midi.toml` for custom mappings:

```toml
[cc_mappings]
1 = "gain"                    # Master Gain
2 = "input_freeze_length"     # Input Freeze Length
3 = "filter_amount"           # Filter Amount
4 = "filter_cutoff"           # Filter Cutoff
5 = "filter_resonance"        # Filter Resonance
6 = "granular_density"        # Granular Density
7 = "granular_size"           # Granular Size
8 = "granular_mix"            # Granular Mix
9 = "reverb_mix_delay_mix"    # Reverb & Delay Mix
10 = "decay_time"             # Reverb & Delay Decay
11 = "dry_wet"                # Master Dry/Wet

[note_mappings]
60 = "input_freeze"           # C4: Toggle Input Freeze
62 = "granular_freeze"        # D4: Toggle Granular Freeze
64 = "blend_mode_mirror"      # E4: Blend Mode 0
65 = "blend_mode_complement"  # F4: Blend Mode 1
67 = "blend_mode_transform"   # G4: Blend Mode 2
```

### OSC Protocol Reference

#### Parameter Control
All parameters are available via individual OSC endpoints. Examples:

```
/chroma/gain f 0.75                    # Set master gain
/chroma/overdriveBias f -0.25           # Set overdrive bias
/chroma/grainIntensity s "pronounced"   # Set grain intensity
/chroma/filterCutoff f 1200             # Set filter cutoff
```

#### Special OSC Messages
```
/chroma/effectsOrder s "granular" s "filter" s "delay"  # Reorder effects
/chroma/sync                            # Request state sync
```

#### OSC State Reception
```
/chroma/state f 0.75 i 1 f 0.5 ...     # Receive complete state (35 args)
/chroma/effectsOrder s "filter" s "delay" ...  # Receive effects order
```

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

```
┌─────────────────┐    OSC Protocol    ┌────────────────────┐
│   TUI Repo    │ ←───────────────→ │ Chroma Engine    │
│ (This Repo)   │                  │ (External)       │
│               │                  │                  │
│ • UI Layout   │                  │ • Audio Engine   │
│ • 38 Parameters │                │ • Effects        │
│ • Effects Order │                │ • State Sync     │
│ • MIDI Mapping  │                │ • OSC Server     │
│ • OSC Client    │                │                  │
│ • OSC Server    │                │                  │
│ • Config System │                │                  │
└─────────────────┘                  └────────────────────┘
```

### Component Architecture
- **TUI/**: Terminal interface with side-by-side layout, dynamic sliders, context-sensitive footer
- **OSC/**: Bidirectional communication, parameter synchronization
- **MIDI/**: Hardware controller integration, mapping system
- **Config/**: User preferences, MIDI configuration, persistence

## Troubleshooting

### Connection Issues

**OSC Communication Failed**
```
Error: OSC connection timeout
Solution: 
1. Verify Chroma SuperCollider engine is running
2. Check localhost ports 57120/9000 availability  
3. Confirm both applications use same OSC protocol version
4. Test with: telnet localhost 57120
```

**MIDI Device Not Detected**
```
Error: No MIDI devices found
Solution:
1. Verify device connection and system recognition
2. Check permissions for MIDI device access
3. Verify ~/.config/chroma/midi.toml exists and is valid
4. Test system MIDI: aconnect -l (Linux) or Audio MIDI Setup (macOS)
```

### Parameter Issues

**Parameter Values Not Updating**
- Check status bar for OSC connection indicator
- Look for "pending changes" indicator in UI
- Verify parameter range compliance
- Test OSC messages manually with oscsend utility

**Effects Reordering Not Working**
- Ensure you're in the global section (navigate to Effects Order in parameter list)
- Use PageUp/PageDown for reordering (arrow keys navigate the list)
- Press Enter from effects list to access parameter panel
- Check OSC connection for order synchronization
- Reset with 'r' key if order becomes corrupted

### Performance Issues

**High CPU Usage**
- Reduce granular density or size parameters
- Disable unused effects in chain
- Check for MIDI controller spam (stuck CC)
- Monitor with: top -p $(pgrep chroma-tui)

**UI Responsiveness**
- Pending changes system prevents OSC overwrites during active adjustment
- Parameters return to normal sync after 500ms of inactivity
- Check network latency for OSC communication
- Focused panel shows green border, unfocused shows gray border
- Footer updates dynamically based on current navigation mode

### System Verification

Test TUI functionality:

```bash
# Test build
go build -o chroma-tui

# Test OSC communication (requires engine running)
echo "Sending test OSC message..." && \
echo "/chroma/gain f 0.5" | oscsend localhost 57120

# Test configuration parsing
go test ./config -v

# Test all components
go test ./... -v
```

### Debug Information

Enable debug logging:
```bash
CHROMA_DEBUG=1 ./chroma-tui
```

### Compatibility

**Platforms:**
- Linux (Ubuntu 22.04+): Full MIDI support via ALSA
- macOS (12+): CoreAudio MIDI integration
- Windows (10+): Win32 API MIDI support

**Go Version:**
- Go 1.24+

**Chroma Engine:**
- Minimum: Chroma SuperCollider v0.3.0+
- Protocol: OSC messages on localhost:57120/9000

## Platform Support

| Platform | MIDI Support | Build Status |
|----------|--------------|--------------|
| Linux    | Native ALSA  | Tested     |
| macOS    | CoreAudio     | Tested     |  
| Windows  | Win32 API    | Tested     |
| Other    | No-CGO Fallback | Tested     |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

For the Chroma audio engine, see the main [Chroma repository](https://github.com/renderorange/chroma).
