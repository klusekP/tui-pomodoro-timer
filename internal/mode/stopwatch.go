package mode

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/klusekP/time_pr_tracking/internal/clock"
	"github.com/klusekP/time_pr_tracking/internal/domain"
	"github.com/klusekP/time_pr_tracking/internal/repository"
	"github.com/klusekP/time_pr_tracking/internal/timer"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

type Stopwatch struct {
	timer *timer.Timer
	repo  repository.Sessions
}

func NewStopwatch(clk clock.Clock, repo repository.Sessions) *Stopwatch {
	return &Stopwatch{
		timer: timer.New(clk),
		repo:  repo,
	}
}

func (s *Stopwatch) Name() string { return "Stopwatch" }

func (s *Stopwatch) State() ui.RunState { return runStateFromTimer(s.timer.State()) }

func (s *Stopwatch) Elapsed() time.Duration { return s.timer.Elapsed() }

func (s *Stopwatch) UsesTaskInput() bool { return true }

func (s *Stopwatch) CapturingKeys() bool { return false }

func (s *Stopwatch) OnTick() tea.Cmd { return nil }

func (s *Stopwatch) HandleKey(msg tea.KeyMsg, task string) tea.Cmd {
	switch msg.String() {
	case " ":
		s.timer.Toggle()
		return nil
	case "s":
		return s.stopAndSave(task)
	case "r":
		s.timer.Reset()
		return ui.Toast("Stopwatch reset")
	}
	return nil
}

func (s *Stopwatch) stopAndSave(task string) tea.Cmd {
	if s.timer.IsIdle() {
		return ui.Toast("Nothing to stop")
	}
	start, end, elapsed := s.timer.Stop()
	session := domain.Session{
		Kind:      domain.KindStopwatch,
		Task:      ui.DisplayTask(task),
		Elapsed:   elapsed,
		StartedAt: start,
		EndedAt:   end,
	}
	repo := s.repo
	return func() tea.Msg {
		saved, err := repo.Save(session)
		if err != nil {
			return ui.SavedMsg{Session: session, Err: err}
		}
		return ui.SavedMsg{Session: saved}
	}
}

func (s *Stopwatch) RenderBody(width int, theme *ui.Theme, bigFont *ui.BigFont) string {
	color := theme.StateColor(s.State())
	clockText := ui.FormatHMS(s.Elapsed())
	big := lipgloss.NewStyle().Foreground(color).Bold(true).Render(bigFont.Render(clockText))
	box := theme.ClockBox.Copy().BorderForeground(color).Render(big)
	inner := width - 4
	if inner < 0 {
		inner = 0
	}
	return lipgloss.PlaceHorizontal(inner, lipgloss.Center, box)
}

func (s *Stopwatch) HelpKeys(theme *ui.Theme) []string {
	return []string{
		theme.KeyCap.Render("space") + " start/pause",
		theme.KeyCap.Render("s") + " stop + save",
		theme.KeyCap.Render("r") + " reset",
	}
}
