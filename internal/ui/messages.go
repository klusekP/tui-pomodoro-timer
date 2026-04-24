package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/klusekP/time_pr_tracking/internal/domain"
)

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type SavedMsg struct {
	Session domain.Session
	Err     error
}

type SessionsLoadedMsg struct {
	Sessions []domain.Session
	Total    time.Duration
	Err      error
}

type AlarmDoneMsg struct{}

type ToastMsg string

func Toast(s string) tea.Cmd {
	return func() tea.Msg { return ToastMsg(s) }
}
