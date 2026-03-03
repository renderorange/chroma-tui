package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

// paletteItem represents a command in the palette.
type paletteItem struct {
	command     Command
	matchScore  int
	matchRanges []int
}

// toggleCommandPalette shows/hides the command palette.
func (m *Model) toggleCommandPalette() {
	m.showCommandPalette = !m.showCommandPalette
	m.commandPaletteText = ""
}

// executeCommand executes a command from the palette input.
func (m *Model) executeCommand(input string) tea.Cmd {
	input = strings.TrimSpace(input)
	if input == "" {
		m.showCommandPalette = false
		return nil
	}

	// Strip leading ":" if present
	input = strings.TrimPrefix(input, ":")
	input = strings.TrimSpace(input)

	if input == "" {
		m.showCommandPalette = false
		return nil
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		m.showCommandPalette = false
		return nil
	}

	cmdName := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	cmd, ok := findCommand(cmdName)
	if !ok {
		m.commandPaletteText = "Unknown command: " + cmdName
		return nil
	}

	m.showCommandPalette = false
	m.commandPaletteText = ""
	return cmd.Handler(m, args)
}

// fuzzyMatchCommands finds commands matching the input using fuzzy search.
func fuzzyMatchCommands(input string) []paletteItem {
	if input == "" {
		// Return all commands when input is empty
		var items []paletteItem
		for _, cmd := range availableCommands() {
			items = append(items, paletteItem{command: cmd, matchScore: 0})
		}
		return items
	}

	var commands []string
	for _, cmd := range availableCommands() {
		commands = append(commands, cmd.Name)
	}

	matches := fuzzy.Find(input, commands)
	var items []paletteItem
	for _, match := range matches {
		cmd, ok := findCommand(match.Str)
		if ok {
			items = append(items, paletteItem{
				command:     cmd,
				matchScore:  match.Score,
				matchRanges: match.MatchedIndexes,
			})
		}
	}
	return items
}

func (m *Model) renderCommandPalette() string {
	if !m.showCommandPalette {
		return ""
	}

	background := m.renderMainBase()
	modalContent := m.renderCommandPaletteModal()
	return m.overlayModal(background, modalContent)
}

// renderCommandPaletteModal creates the command palette modal content.
func (m *Model) renderCommandPaletteModal() string {
	// Get palette input from commandPaletteText
	input := m.commandPaletteText

	// Ensure input starts with ":" for display
	displayInput := input
	if !strings.HasPrefix(displayInput, ":") {
		displayInput = ":" + displayInput
	}

	// For matching, strip the ":" if present
	matchInput := input
	if strings.HasPrefix(matchInput, ":") {
		matchInput = matchInput[1:]
	}

	// Match commands
	matches := fuzzyMatchCommands(matchInput)

	// Create styles
	primaryStyle := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	secondaryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	mutedStyle := lipgloss.NewStyle().Foreground(colorTextMuted)
	accentStyle := lipgloss.NewStyle().Foreground(colorAccent)
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 1)

	// Calculate dimensions (match help modal width)
	width := 60
	if m.width > 0 && m.width < width+4 {
		width = m.width - 4
	}
	if width < 40 {
		width = 40
	}

	var b strings.Builder

	// Input line - show what user typed with ":" prefix
	b.WriteString(primaryStyle.Render(displayInput))
	b.WriteString("\n")

	// Separator
	b.WriteString(mutedStyle.Render(strings.Repeat("─", width-2)))
	b.WriteString("\n")

	// Matches
	maxMatches := 5
	if len(matches) > maxMatches {
		matches = matches[:maxMatches]
	}

	for _, item := range matches {
		cmd := item.command

		// Highlight matched characters
		name := highlightMatches(cmd.Name, item.matchRanges, accentStyle)

		// Aliases
		aliases := ""
		if len(cmd.Aliases) > 0 {
			aliases = mutedStyle.Render(" (" + strings.Join(cmd.Aliases, ", ") + ")")
		}

		// Description
		desc := secondaryStyle.Render(" " + cmd.Description)

		b.WriteString(name)
		b.WriteString(aliases)
		b.WriteString(desc)
		b.WriteString("\n")
	}

	// Apply modal styling
	return modalStyle.Width(width).Render(b.String())
}

func highlightMatches(s string, ranges []int, highlightStyle lipgloss.Style) string {
	if len(ranges) == 0 {
		return s
	}

	var result strings.Builder
	lastEnd := 0

	for _, idx := range ranges {
		if idx >= len(s) {
			break
		}

		// Write unhighlighted portion
		if idx > lastEnd {
			result.WriteString(s[lastEnd:idx])
		}

		// Write highlighted character
		result.WriteString(highlightStyle.Render(string(s[idx])))
		lastEnd = idx + 1
	}

	// Write remaining unhighlighted portion
	if lastEnd < len(s) {
		result.WriteString(s[lastEnd:])
	}

	return result.String()
}
