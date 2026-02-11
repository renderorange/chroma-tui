# Chroma TUI

Terminal-based interface for Chroma audio effects system.

This repository contains the Go-based terminal user interface (TUI) component of the Chroma audio effects system. The TUI communicates with the Chroma SuperCollider audio engine via OSC (Open Sound Control).

## External Dependency

This TUI requires the Chroma SuperCollider audio engine to function. Please see the [Chroma repository](https://github.com/renderorange/chroma) for setup instructions.

## Installation

```bash
git clone https://github.com/renderorange/chroma-tui.git
cd chroma-tui
go build -o chroma-tui
```

## Usage

```bash
./chroma-tui
```

*Requires the Chroma SuperCollider engine to be running.*
