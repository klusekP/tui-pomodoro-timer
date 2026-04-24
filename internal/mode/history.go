package mode

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/klusekP/time_pr_tracking/internal/domain"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

const historyTableRows = 20

type History struct {
	all        []domain.Session
	totalToday time.Duration

	filter        textinput.Model
	editingFilter bool
}

func NewHistory() *History {
	ti := textinput.New()
	ti.Placeholder = "task name…"
	ti.CharLimit = 80
	ti.Width = 40
	ti.Prompt = ""
	return &History{filter: ti}
}

func (h *History) Name() string { return "Sessions" }

func (h *History) State() ui.RunState { return ui.StateIdle }

func (h *History) OnTick() tea.Cmd { return nil }

func (h *History) UsesTaskInput() bool { return false }

func (h *History) CapturingKeys() bool { return h.editingFilter }

func (h *History) HandleKey(msg tea.KeyMsg, _ string) tea.Cmd {
	if h.editingFilter {
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			h.editingFilter = false
			h.filter.Blur()
			return nil
		}
		var cmd tea.Cmd
		h.filter, cmd = h.filter.Update(msg)
		return cmd
	}

	switch msg.String() {
	case "f", "/":
		h.editingFilter = true
		h.filter.Focus()
		return textinput.Blink
	case "c":
		if strings.TrimSpace(h.filter.Value()) == "" {
			return nil
		}
		h.filter.SetValue("")
		return ui.Toast("Filter cleared")
	}
	return nil
}

func (h *History) HelpKeys(theme *ui.Theme) []string {
	switch {
	case h.editingFilter:
		return []string{
			theme.KeyCap.Render("enter") + " apply",
			theme.KeyCap.Render("esc") + " done",
		}
	case strings.TrimSpace(h.filter.Value()) != "":
		return []string{
			theme.KeyCap.Render("f") + " edit filter",
			theme.KeyCap.Render("c") + " clear",
		}
	default:
		return []string{
			theme.KeyCap.Render("f") + " filter",
		}
	}
}

func (h *History) SetSessions(recent []domain.Session, totalToday time.Duration) {
	h.all = recent
	h.totalToday = totalToday
}

func (h *History) query() string {
	return strings.ToLower(strings.TrimSpace(h.filter.Value()))
}

func (h *History) filtered() []domain.Session {
	q := h.query()
	if q == "" {
		return h.all
	}
	out := make([]domain.Session, 0, len(h.all))
	for _, s := range h.all {
		if strings.Contains(strings.ToLower(s.Task), q) {
			out = append(out, s)
		}
	}
	return out
}

type kindStat struct {
	count int
	total time.Duration
}

func summarize(sessions []domain.Session) (total time.Duration, byKind map[domain.SessionKind]kindStat) {
	byKind = map[domain.SessionKind]kindStat{}
	for _, s := range sessions {
		d := s.Duration()
		total += d
		st := byKind[s.Kind]
		st.count++
		st.total += d
		byKind[s.Kind] = st
	}
	return total, byKind
}

func (h *History) RenderBody(width int, theme *ui.Theme, _ *ui.BigFont) string {
	panelW := width - 4
	if panelW < 20 {
		panelW = 20
	}
	h.filter.Width = ui.Clamp(panelW-20, 20, 60)

	filter := h.renderFilter(theme)

	if len(h.all) == 0 {
		title := theme.PanelTitle.Render("Sessions")
		empty := lipgloss.NewStyle().
			Foreground(theme.ColorTextDim).
			Italic(true).
			Render("  No saved sessions yet — stop the Stopwatch to save one.")
		inner := lipgloss.JoinVertical(lipgloss.Left, title, empty)
		panel := theme.Panel.Copy().Width(panelW).Render(inner)
		return lipgloss.JoinVertical(lipgloss.Left, filter, "", panel)
	}

	filtered := h.filtered()
	summary := h.renderSummary(panelW, theme, filtered)
	table := h.renderTable(panelW, theme, filtered)
	return lipgloss.JoinVertical(lipgloss.Left, filter, "", summary, "", table)
}

func (h *History) renderFilter(theme *ui.Theme) string {
	label := theme.InputLabel.Render("Filter:")
	var box string
	if h.editingFilter {
		box = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.ColorAccent).
			Padding(0, 1).
			Render(h.filter.View())
	} else {
		value := strings.TrimSpace(h.filter.Value())
		if value == "" {
			value = lipgloss.NewStyle().
				Foreground(theme.ColorTextDim).
				Italic(true).
				Render("(press 'f' to filter by task name)")
		}
		box = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.ColorIdle).
			Padding(0, 1).
			Render(value)
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, label, box)
}

func (h *History) renderSummary(panelW int, theme *ui.Theme, sessions []domain.Session) string {
	total, byKind := summarize(sessions)
	q := strings.TrimSpace(h.filter.Value())

	var title string
	if q != "" {
		title = theme.PanelTitle.Render(fmt.Sprintf("Summary   (filter: %q)", q))
	} else {
		title = theme.PanelTitle.Render("Summary")
	}

	dim := lipgloss.NewStyle().Foreground(theme.ColorTextDim)
	val := lipgloss.NewStyle().Foreground(theme.ColorText).Bold(true)
	line := func(label, value string) string {
		return "  " + dim.Render(fmt.Sprintf("%-18s", label)) + val.Render(value)
	}

	lines := []string{title}
	if q != "" {
		lines = append(lines,
			line("Matches:", fmt.Sprintf("%d sessions", len(sessions))),
			line("Total:", ui.FormatHours(total)),
		)
	} else {
		lines = append(lines,
			line("Today:", ui.FormatHours(h.totalToday)),
			line("Total (loaded):", fmt.Sprintf("%s   (%d sessions)",
				ui.FormatHours(total), len(sessions))),
		)
	}

	order := []domain.SessionKind{domain.KindStopwatch, domain.KindPomodoro}
	seen := map[domain.SessionKind]bool{}
	for _, k := range order {
		if st, ok := byKind[k]; ok {
			lines = append(lines, line("  "+k.Label()+":",
				fmt.Sprintf("%s   (%d)", ui.FormatHours(st.total), st.count)))
			seen[k] = true
		}
	}
	for k, st := range byKind {
		if seen[k] {
			continue
		}
		lines = append(lines, line("  "+k.Label()+":",
			fmt.Sprintf("%s   (%d)", ui.FormatHours(st.total), st.count)))
	}

	inner := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return theme.Panel.Copy().Width(panelW).Render(inner)
}

func (h *History) renderTable(panelW int, theme *ui.Theme, sessions []domain.Session) string {
	q := strings.TrimSpace(h.filter.Value())

	var titleText string
	if q != "" {
		titleText = fmt.Sprintf("Matching sessions (%d)", len(sessions))
	} else {
		titleText = fmt.Sprintf("Recent sessions (%d)", len(sessions))
	}
	title := theme.PanelTitle.Render(titleText)

	if len(sessions) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(theme.ColorTextDim).
			Italic(true).
			Render("  No matching sessions.")
		inner := lipgloss.JoinVertical(lipgloss.Left, title, empty)
		return theme.Panel.Copy().Width(panelW).Render(inner)
	}

	kindW := 10
	durW := 10
	timeW := 19
	innerW := panelW - 4
	taskW := innerW - kindW - durW - timeW - 6
	if taskW < 10 {
		taskW = 10
	}

	head := lipgloss.NewStyle().
		Foreground(theme.ColorTextDim).
		Bold(true).
		Render(fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s",
			kindW, "Mode", taskW, "Task", durW, "Time", timeW, "Ended"))

	shown := sessions
	truncated := 0
	if len(shown) > historyTableRows {
		truncated = len(shown) - historyTableRows
		shown = shown[:historyTableRows]
	}

	rows := make([]string, 0, len(shown)+3)
	rows = append(rows, title, head)
	for _, s := range shown {
		rows = append(rows, fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s",
			kindW, s.Kind.Label(),
			taskW, ui.Truncate(ui.DisplayTask(s.Task), taskW),
			durW, ui.FormatHMS(s.Duration()),
			timeW, s.EndedAt.Format("2006-01-02 15:04"),
		))
	}
	if truncated > 0 {
		rows = append(rows, lipgloss.NewStyle().
			Foreground(theme.ColorTextDim).
			Italic(true).
			Render(fmt.Sprintf("  … and %d more", truncated)))
	}
	inner := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return theme.Panel.Copy().Width(panelW).Render(inner)
}
