package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/klusekP/time_pr_tracking/internal/mode"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

type ViewState struct {
	Width, Height int

	Active    mode.Mode
	Modes     []mode.Mode
	ActiveIdx int

	Task        textinput.Model
	EditingTask bool

	TotalToday time.Duration

	Toast  string
	ErrMsg string
}

type Renderer struct {
	theme   *ui.Theme
	bigFont *ui.BigFont
}

func NewRenderer(theme *ui.Theme, bigFont *ui.BigFont) *Renderer {
	return &Renderer{theme: theme, bigFont: bigFont}
}

func (r *Renderer) Render(vs ViewState) string {
	if vs.Width == 0 || vs.Height == 0 {
		return "Loading..."
	}

	header := r.renderHeader(vs)
	tabs := r.renderTabs(vs)
	body := vs.Active.RenderBody(vs.Width, r.theme, r.bigFont)
	footer := r.renderFooter(vs)

	parts := []string{header, "", tabs, ""}
	if vs.Active.UsesTaskInput() {
		parts = append(parts, r.renderTaskRow(vs), "")
	}
	parts = append(parts, body)
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	contentHeight := lipgloss.Height(content)
	footerHeight := lipgloss.Height(footer)
	avail := vs.Height - contentHeight - footerHeight - 2
	if avail < 0 {
		avail = 0
	}
	spacer := strings.Repeat("\n", avail)

	full := lipgloss.JoinVertical(lipgloss.Left, content, spacer, footer)
	return r.theme.App.Render(full)
}

func (r *Renderer) renderHeader(vs ViewState) string {
	title := r.theme.Title.Render("⏱  Time Tracker")
	status := r.renderStatus(vs.Active.State())

	today := r.theme.Subtitle.Render(
		fmt.Sprintf("Saved today: %s", ui.FormatHMS(vs.TotalToday)),
	)

	left := lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", status)
	gap := vs.Width - lipgloss.Width(left) - lipgloss.Width(today) - 4
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + today
}

func (r *Renderer) renderStatus(s ui.RunState) string {
	switch s {
	case ui.StateRunning:
		return r.theme.StatusRun.Render("● RUNNING")
	case ui.StatePaused:
		return r.theme.StatusPause.Render("❚❚ PAUSED")
	case ui.StateFinished:
		return r.theme.StatusAlarm.Render("⏰ DONE!")
	default:
		return r.theme.StatusIdle.Render("○ IDLE")
	}
}

func (r *Renderer) renderTabs(vs ViewState) string {
	parts := make([]string, 0, len(vs.Modes)*2+1)
	for i, m := range vs.Modes {
		label := fmt.Sprintf(" %d  %s ", i+1, m.Name())
		if i == vs.ActiveIdx {
			parts = append(parts, r.theme.TabActive.Render(label))
		} else {
			parts = append(parts, r.theme.TabInactive.Render(label))
		}
		parts = append(parts, " ")
	}
	parts = append(parts, r.theme.Subtitle.Render(" (tab — switch mode)"))
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (r *Renderer) renderTaskRow(vs ViewState) string {
	label := r.theme.InputLabel.Render("Task:")
	var box string
	if vs.EditingTask {
		box = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(r.theme.ColorAccent).
			Padding(0, 1).
			Render(vs.Task.View())
	} else {
		value := strings.TrimSpace(vs.Task.Value())
		if value == "" {
			value = lipgloss.NewStyle().
				Foreground(r.theme.ColorTextDim).
				Italic(true).
				Render("(empty — press 'e' to edit)")
		}
		box = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(r.theme.ColorIdle).
			Padding(0, 1).
			Render(value)
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, label, box)
}

func (r *Renderer) renderFooter(vs ViewState) string {
	keys := vs.Active.HelpKeys(r.theme)
	keys = append(keys,
		r.theme.KeyCap.Render("e")+" task",
		r.theme.KeyCap.Render("tab")+" mode",
		r.theme.KeyCap.Render("q")+" quit",
	)
	help := r.theme.Help.Render(strings.Join(keys, "   "))

	var extra string
	switch {
	case vs.ErrMsg != "":
		extra = r.theme.Error.Render("✗ " + vs.ErrMsg)
	case vs.Toast != "":
		extra = r.theme.Toast.Render("✓ " + vs.Toast)
	}
	if extra != "" {
		return lipgloss.JoinVertical(lipgloss.Left, extra, help)
	}
	return help
}
