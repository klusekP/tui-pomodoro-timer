package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	ColorPrimary  lipgloss.Color
	ColorAccent   lipgloss.Color
	ColorRunning  lipgloss.Color
	ColorPaused   lipgloss.Color
	ColorIdle     lipgloss.Color
	ColorAlarm    lipgloss.Color
	ColorText     lipgloss.Color
	ColorTextDim  lipgloss.Color
	ColorBgPanel  lipgloss.Color
	ColorBarFill  lipgloss.Color
	ColorBarEmpty lipgloss.Color

	App         lipgloss.Style
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	ClockBox    lipgloss.Style
	StatusRun   lipgloss.Style
	StatusPause lipgloss.Style
	StatusIdle  lipgloss.Style
	StatusAlarm lipgloss.Style
	PanelTitle  lipgloss.Style
	Panel       lipgloss.Style
	InputLabel  lipgloss.Style
	Help        lipgloss.Style
	KeyCap      lipgloss.Style
	Toast       lipgloss.Style
	Error       lipgloss.Style
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
}

func NewDefaultTheme() *Theme {
	t := &Theme{
		ColorPrimary:  lipgloss.Color("#7D56F4"),
		ColorAccent:   lipgloss.Color("#F25D94"),
		ColorRunning:  lipgloss.Color("#43BF6D"),
		ColorPaused:   lipgloss.Color("#F5A623"),
		ColorIdle:     lipgloss.Color("#6B7280"),
		ColorAlarm:    lipgloss.Color("#FF5F5F"),
		ColorText:     lipgloss.Color("#E6E6E6"),
		ColorTextDim:  lipgloss.Color("#8B8B8B"),
		ColorBgPanel:  lipgloss.Color("#1E1E2E"),
		ColorBarFill:  lipgloss.Color("#7D56F4"),
		ColorBarEmpty: lipgloss.Color("#3A3A4A"),
	}

	t.App = lipgloss.NewStyle().Padding(1, 2)

	t.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.ColorPrimary).
		Padding(0, 1)

	t.Subtitle = lipgloss.NewStyle().
		Foreground(t.ColorTextDim).
		Italic(true)

	t.ClockBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.ColorPrimary).
		Padding(1, 4).
		Align(lipgloss.Center)

	t.StatusRun = lipgloss.NewStyle().Foreground(t.ColorRunning).Bold(true)
	t.StatusPause = lipgloss.NewStyle().Foreground(t.ColorPaused).Bold(true)
	t.StatusIdle = lipgloss.NewStyle().Foreground(t.ColorIdle).Bold(true)
	t.StatusAlarm = lipgloss.NewStyle().Foreground(t.ColorAlarm).Bold(true)

	t.PanelTitle = lipgloss.NewStyle().
		Foreground(t.ColorAccent).
		Bold(true).
		Padding(0, 1)

	t.Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.ColorIdle).
		Padding(0, 1)

	t.InputLabel = lipgloss.NewStyle().
		Foreground(t.ColorAccent).
		Bold(true).
		MarginRight(1)

	t.Help = lipgloss.NewStyle().
		Foreground(t.ColorTextDim).
		Padding(0, 1)

	t.KeyCap = lipgloss.NewStyle().
		Foreground(t.ColorText).
		Background(t.ColorBgPanel).
		Padding(0, 1).
		Bold(true)

	t.Toast = lipgloss.NewStyle().Foreground(t.ColorRunning).Bold(true)
	t.Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5555")).
		Bold(true)

	t.TabActive = lipgloss.NewStyle().
		Foreground(t.ColorText).
		Background(t.ColorPrimary).
		Padding(0, 2).
		Bold(true)

	t.TabInactive = lipgloss.NewStyle().
		Foreground(t.ColorTextDim).
		Padding(0, 2)

	return t
}

func (t *Theme) StateColor(s RunState) lipgloss.Color {
	switch s {
	case StateRunning:
		return t.ColorRunning
	case StatePaused:
		return t.ColorPaused
	case StateFinished:
		return t.ColorAlarm
	default:
		return t.ColorIdle
	}
}
