package mode

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/klusekP/time_pr_tracking/internal/domain"
	"github.com/klusekP/time_pr_tracking/internal/timer"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

type Mode interface {
	Name() string

	State() ui.RunState

	HandleKey(msg tea.KeyMsg, task string) tea.Cmd

	OnTick() tea.Cmd

	RenderBody(width int, theme *ui.Theme, bigFont *ui.BigFont) string

	HelpKeys(theme *ui.Theme) []string

	UsesTaskInput() bool

	CapturingKeys() bool
}

type SessionsAware interface {
	SetSessions(recent []domain.Session, totalToday time.Duration)
}

func runStateFromTimer(s timer.State) ui.RunState {
	switch s {
	case timer.Running:
		return ui.StateRunning
	case timer.Paused:
		return ui.StatePaused
	default:
		return ui.StateIdle
	}
}
