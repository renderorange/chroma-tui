package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	spectrumBands      = 8
	minVisualizerWidth = 16
	waveformHeight     = 7 // Number of vertical levels for waveform (0-6)
)

var (
	spectrumBarStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	waveformStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	labelStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// renderSpectrum renders 8-band spectrum analyzer
func renderSpectrum(spectrum [spectrumBands]float32, width int) string {
	if width < minVisualizerWidth {
		return ""
	}

	barWidth := width / spectrumBands
	bars := make([]string, spectrumBands)
	levels := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	for i, val := range spectrum {
		// Clamp value 0-1
		if val < 0 {
			val = 0
		}
		if val > 1 {
			val = 1
		}

		// Map to 8 levels
		levelIdx := int(val * 7)
		bar := strings.Repeat(levels[levelIdx], barWidth)

		// Style with pink accent using pre-created style
		bars[i] = spectrumBarStyle.Render(bar)
	}

	return strings.Join(bars, "")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// renderMultiLineSpectrum renders spectrum with vertical height
func renderMultiLineSpectrum(spectrum [spectrumBands]float32, width, height int) string {
	if width < minVisualizerWidth || height < 1 {
		return ""
	}

	barWidth := width / spectrumBands
	levels := []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	// Pre-calculate styled characters for performance
	styledLevels := make([]string, len(levels))
	for i, level := range levels {
		styledLevels[i] = spectrumBarStyle.Render(strings.Repeat(level, barWidth))
	}

	var lines strings.Builder
	for line := height - 1; line >= 0; line-- {
		var lineBars strings.Builder

		for _, val := range spectrum {
			// Clamp value 0-1
			if val < 0 {
				val = 0
			}
			if val > 1 {
				val = 1
			}

			// Map to height levels
			threshold := float32(line) / float32(height)
			var levelIdx int

			if val >= threshold {
				// Fix threshold calculation bug: handle division by zero
				if threshold >= 1.0 {
					levelIdx = len(levels) - 1
				} else {
					normalizedVal := (val - threshold) / (1.0 - threshold)
					levelIdx = int(normalizedVal * float32(len(levels)-1))
					levelIdx = min(levelIdx, len(levels)-1)
				}
			} else {
				levelIdx = 0 // empty space
			}

			lineBars.WriteString(styledLevels[levelIdx])
		}

		if line > 0 {
			lines.WriteString(lineBars.String())
			lines.WriteString("\n")
		} else {
			lines.WriteString(lineBars.String())
		}
	}

	return lines.String()
}

// renderWaveform creates oscilloscope-style waveform display
func renderWaveform(waveform [64]float32, width int) string {
	if width < minVisualizerWidth {
		return ""
	}

	// Create vertical slices for display (reduce to terminal width)
	step := len(waveform) / width
	if step < 1 {
		step = 1
	}

	var chars []string
	prevY := 0
	prevI := 0 // Track previous actual waveform index for continuity

	for i := 0; i < len(waveform); i += step {
		// Normalize waveform value to -1 to 1
		val := waveform[i]
		if val < -1 {
			val = -1
		}
		if val > 1 {
			val = 1
		}

		// Map to character set - fix Y mapping to use proper 0..6 range
		var char string
		y := int((val + 1) * float32(waveformHeight) / 2.0) // Map -1..1 to 0..6

		if i == 0 {
			// First point
			char = "│"
		} else {
			// For continuity when step > 1, interpolate from previous actual point
			if step > 1 && prevI+step < len(waveform) {
				// Look at the actual previous waveform point for continuity
				prevVal := waveform[prevI]
				if prevVal < -1 {
					prevVal = -1
				}
				if prevVal > 1 {
					prevVal = 1
				}
				prevY = int((prevVal + 1) * float32(waveformHeight) / 2.0)
			}

			// Determine connection character based on direction
			// Add some randomness to get different corner types
			if y > prevY {
				// Going up - alternate between corner types
				if i%3 == 0 {
					char = "╭"
				} else {
					char = "╯"
				}
			} else if y < prevY {
				// Going down - alternate between corner types
				if i%3 == 0 {
					char = "╰"
				} else {
					char = "╮"
				}
			} else {
				char = "─"
			}
		}

		chars = append(chars, char)
		prevY = y
		prevI = i
	}

	return waveformStyle.Render(strings.Join(chars, ""))
}

// generateDynamicTitle creates a title that adapts to terminal width
func generateDynamicTitle(title string, width int) string {
	if width < len(title)+4 {
		return titleStyle.Render(title)
	}

	dashCount := width - len(title) - 4
	leftDashes := dashCount / 2
	rightDashes := dashCount - leftDashes

	return titleStyle.Render(
		strings.Repeat("─", leftDashes) + " " + title + " " + strings.Repeat("─", rightDashes),
	)
}

// renderFrequencyLabels creates simplified frequency label positioning
func renderFrequencyLabels(width int) string {
	if width < minVisualizerWidth {
		return ""
	}

	labels := []string{"30Hz", "120Hz", "375Hz", "1kHz", "3kHz", "5kHz", "9kHz", "16kHz"}
	barWidth := width / spectrumBands

	type labelPosition struct {
		pos   int
		label string
	}

	var labelPositions []labelPosition
	for i, label := range labels {
		if barWidth >= len(label) {
			// Center label in bar area
			startPos := i*barWidth + (barWidth-len(label))/2
			if startPos >= 0 && startPos+len(label) <= width {
				labelPositions = append(labelPositions, labelPosition{startPos, label})
			}
		}
	}

	// Build label line
	labelLine := make([]rune, width)
	for i := range labelLine {
		labelLine[i] = ' '
	}

	for _, lp := range labelPositions {
		for j, char := range lp.label {
			if lp.pos+j < width {
				labelLine[lp.pos+j] = char
			}
		}
	}

	return labelStyle.Render(string(labelLine))
}

// validateSpectrumData checks if spectrum data contains valid values
func validateSpectrumData(spectrum [spectrumBands]float32) bool {
	for _, val := range spectrum {
		if val > 0 {
			return true // Found at least one non-zero value
		}
	}
	return false // All values are zero or negative
}

// renderVisualizer creates full visualizer section with spectrum and waveform
func (m Model) renderVisualizer(width int) string {
	if width < minVisualizerWidth {
		return ""
	}

	// Input validation
	if !validateSpectrumData(m.Spectrum) {
		return generateDynamicTitle("SPECTRUM", width) + "\n" +
			labelStyle.Render("No spectrum data available")
	}

	var sections []string

	// Spectrum title with dynamic width
	sections = append(sections, generateDynamicTitle("SPECTRUM", width))

	// Multi-line spectrum (6 lines tall)
	spectrumLines := renderMultiLineSpectrum(m.Spectrum, width, 6)
	sections = append(sections, spectrumLines)

	// Frequency labels with simplified positioning
	sections = append(sections, renderFrequencyLabels(width))

	// Waveform section
	sections = append(sections, generateDynamicTitle("WAVEFORM", width))
	waveformLine := renderWaveform(m.Waveform, width)
	sections = append(sections, waveformLine)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
