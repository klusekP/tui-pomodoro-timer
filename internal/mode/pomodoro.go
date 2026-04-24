package mode

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/klusekP/time_pr_tracking/internal/alarm"
	"github.com/klusekP/time_pr_tracking/internal/clock"
	"github.com/klusekP/time_pr_tracking/internal/timer"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

const (
	minPomodoroTarget = 1 * time.Minute
	maxPomodoroTarget = 8 * time.Hour
)

type Pomodoro struct {
	timer    *timer.Timer
	alarm    alarm.Alarm
	target   time.Duration
	finished bool
}

func NewPomodoro(clk clock.Clock, a alarm.Alarm, defaultTarget time.Duration) *Pomodoro {
	if a == nil {
		a = alarm.Noop{}
	}
	return &Pomodoro{
		timer:  timer.New(clk),
		alarm:  a,
		target: ui.ClampDuration(defaultTarget, minPomodoroTarget, maxPomodoroTarget),
	}
}

func (p *Pomodoro) Name() string { return "Pomodoro" }

func (p *Pomodoro) UsesTaskInput() bool { return false }

func (p *Pomodoro) CapturingKeys() bool { return false }

func (p *Pomodoro) State() ui.RunState {
	if p.finished {
		return ui.StateFinished
	}
	return runStateFromTimer(p.timer.State())
}

func (p *Pomodoro) Elapsed() time.Duration {
	if p.finished {
		return p.target
	}
	return p.timer.Elapsed()
}

func (p *Pomodoro) Target() time.Duration { return p.target }

func (p *Pomodoro) Remaining() time.Duration {
	if r := p.target - p.Elapsed(); r > 0 {
		return r
	}
	return 0
}

func (p *Pomodoro) Ratio() float64 {
	if p.target <= 0 {
		return 0
	}
	r := float64(p.Elapsed()) / float64(p.target)
	if r > 1 {
		return 1
	}
	if r < 0 {
		return 0
	}
	return r
}

func (p *Pomodoro) OnTick() tea.Cmd {
	if p.finished || p.timer.State() != timer.Running {
		return nil
	}
	if p.timer.Elapsed() < p.target {
		return nil
	}
	p.finish()
	return tea.Batch(
		p.playAlarm(),
		ui.Toast(fmt.Sprintf("Pomodoro %s finished!", ui.FormatHMS(p.target))),
	)
}

func (p *Pomodoro) finish() {
	p.timer.Pause()
	p.finished = true
}

func (p *Pomodoro) playAlarm() tea.Cmd {
	a := p.alarm
	return func() tea.Msg {
		_ = a.Play()
		return ui.AlarmDoneMsg{}
	}
}

func (p *Pomodoro) HandleKey(msg tea.KeyMsg, _ string) tea.Cmd {
	switch msg.String() {
	case " ":

		if p.finished {
			p.reset()
			p.timer.Start()
			return nil
		}
		p.timer.Toggle()
		return nil

	case "s":
		if p.timer.IsIdle() && !p.finished {
			return ui.Toast("Nothing to stop")
		}
		p.reset()
		return ui.Toast("Pomodoro stopped")

	case "r":
		p.reset()
		return ui.Toast("Pomodoro reset")

	case "w":
		return p.setTarget(25 * time.Minute)
	case "b":
		return p.setTarget(5 * time.Minute)
	case "l":
		return p.setTarget(15 * time.Minute)
	case "+", "=":
		return p.setTarget(p.target + 5*time.Minute)
	case "-", "_":
		return p.setTarget(p.target - 5*time.Minute)
	}
	return nil
}

func (p *Pomodoro) reset() {
	p.timer.Reset()
	p.finished = false
}

func (p *Pomodoro) setTarget(d time.Duration) tea.Cmd {
	p.target = ui.ClampDuration(d, minPomodoroTarget, maxPomodoroTarget)
	if p.State() == ui.StateIdle || p.State() == ui.StateFinished {
		return ui.Toast(fmt.Sprintf("Pomodoro duration set: %s", ui.FormatHMS(p.target)))
	}
	return ui.Toast(fmt.Sprintf("New target: %s (session in progress)", ui.FormatHMS(p.target)))
}

func (p *Pomodoro) RenderBody(width int, theme *ui.Theme, bigFont *ui.BigFont) string {
	color := theme.StateColor(p.State())

	clockText := ui.FormatHMS(p.Remaining())
	big := lipgloss.NewStyle().Foreground(color).Bold(true).Render(bigFont.Render(clockText))

	bar := ui.RenderProgressBar(p.Ratio(), 28, color, theme.ColorBarEmpty)

	meta := fmt.Sprintf("Target: %s    %s / %s    %3.0f%%",
		ui.FormatHMS(p.target),
		ui.FormatHMS(p.Elapsed()),
		ui.FormatHMS(p.target),
		p.Ratio()*100,
	)
	metaStyled := lipgloss.NewStyle().Foreground(theme.ColorTextDim).Render(meta)

	rows := []string{big, "", bar, metaStyled}
	if p.finished {
		banner := lipgloss.NewStyle().
			Foreground(theme.ColorAlarm).
			Bold(true).
			Blink(true).
			Render("⏰  POMODORO FINISHED!  ⏰")
		rows = append(rows, "", banner)
	}
	inner := lipgloss.JoinVertical(lipgloss.Center, rows...)
	box := theme.ClockBox.Copy().BorderForeground(color).Render(inner)
	avail := width - 4
	if avail < 0 {
		avail = 0
	}
	return lipgloss.PlaceHorizontal(avail, lipgloss.Center, box)
}

func (p *Pomodoro) HelpKeys(theme *ui.Theme) []string {
	return []string{
		theme.KeyCap.Render("space") + " start/pause",
		theme.KeyCap.Render("s") + " stop",
		theme.KeyCap.Render("r") + " reset",
		theme.KeyCap.Render("w/b/l") + " 25/5/15 min",
		theme.KeyCap.Render("+/-") + " ±5 min",
	}
}
