package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderProgressBar(ratio float64, width int, fill, empty lipgloss.Color) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	if width < 1 {
		width = 1
	}
	filled := int(float64(width)*ratio + 0.5)
	if filled > width {
		filled = width
	}
	filledStyle := lipgloss.NewStyle().Foreground(fill)
	emptyStyle := lipgloss.NewStyle().Foreground(empty)
	return filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", width-filled))
}
