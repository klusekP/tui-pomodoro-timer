package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/klusekP/time_pr_tracking/internal/domain"
	"github.com/klusekP/time_pr_tracking/internal/mode"
	"github.com/klusekP/time_pr_tracking/internal/repository"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

const recentLimit = 2000

type Model struct {
	repo     repository.Sessions
	renderer *Renderer

	width, height int

	modes     []mode.Mode
	activeIdx int

	taskInput   textinput.Model
	editingTask bool

	total time.Duration

	toast     string
	toastTime time.Time
	errMsg    string
}

func NewModel(repo repository.Sessions, renderer *Renderer, modes []mode.Mode) Model {
	ti := textinput.New()
	ti.Placeholder = "What are you working on?"
	ti.CharLimit = 80
	ti.Width = 40
	ti.Prompt = ""

	return Model{
		repo:      repo,
		renderer:  renderer,
		modes:     modes,
		activeIdx: 0,
		taskInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(ui.Tick(), loadRecentCmd(m.repo))
}

func loadRecentCmd(repo repository.Sessions) tea.Cmd {
	return func() tea.Msg {
		sessions, err := repo.Recent(recentLimit)
		if err != nil {
			return ui.SessionsLoadedMsg{Err: err}
		}
		total, err := repo.TotalToday()
		if err != nil {
			return ui.SessionsLoadedMsg{Sessions: sessions, Err: err}
		}
		return ui.SessionsLoadedMsg{Sessions: sessions, Total: total}
	}
}

func (m *Model) broadcastSessions(sessions []domain.Session, total time.Duration) {
	for _, md := range m.modes {
		if aware, ok := md.(mode.SessionsAware); ok {
			aware.SetSessions(sessions, total)
		}
	}
}

func (m Model) active() mode.Mode { return m.modes[m.activeIdx] }

func (m Model) viewState() ViewState {
	return ViewState{
		Width:       m.width,
		Height:      m.height,
		Active:      m.active(),
		Modes:       m.modes,
		ActiveIdx:   m.activeIdx,
		Task:        m.taskInput,
		EditingTask: m.editingTask,
		TotalToday:  m.total,
		Toast:       m.toast,
		ErrMsg:      m.errMsg,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.taskInput.Width = ui.Clamp(msg.Width-20, 20, 60)
		return m, nil

	case ui.TickMsg:
		var cmds []tea.Cmd
		for _, md := range m.modes {
			if c := md.OnTick(); c != nil {
				cmds = append(cmds, c)
			}
		}
		if m.toast != "" && time.Since(m.toastTime) > 3*time.Second {
			m.toast = ""
		}
		cmds = append(cmds, ui.Tick())
		return m, tea.Batch(cmds...)

	case ui.SessionsLoadedMsg:
		if msg.Err != nil {
			m.errMsg = msg.Err.Error()
		} else {
			m.total = msg.Total
			m.errMsg = ""
			m.broadcastSessions(msg.Sessions, msg.Total)
		}
		return m, nil

	case ui.SavedMsg:
		if msg.Err != nil {
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.toast = fmt.Sprintf("Saved: %s  (%s)",
			ui.DisplayTask(msg.Session.Task), ui.FormatHMS(msg.Session.Duration()))
		m.toastTime = time.Now()
		return m, loadRecentCmd(m.repo)

	case ui.ToastMsg:
		if string(msg) != "" {
			m.toast = string(msg)
			m.toastTime = time.Now()
		}
		return m, nil

	case ui.AlarmDoneMsg:
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	if m.editingTask {
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			m.editingTask = false
			m.taskInput.Blur()
			return m, nil
		}
		var cmd tea.Cmd
		m.taskInput, cmd = m.taskInput.Update(msg)
		return m, cmd
	}

	if m.active().CapturingKeys() {
		return m, m.active().HandleKey(msg, m.taskInput.Value())
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "tab":
		m.activeIdx = (m.activeIdx + 1) % len(m.modes)
		return m, nil
	case "e", "i":

		if !m.active().UsesTaskInput() {
			return m, nil
		}
		m.editingTask = true
		m.taskInput.Focus()
		return m, textinput.Blink
	}

	for i := range m.modes {
		if msg.String() == fmt.Sprintf("%d", i+1) {
			m.activeIdx = i
			return m, nil
		}
	}

	return m, m.active().HandleKey(msg, m.taskInput.Value())
}

func (m Model) View() string {
	return m.renderer.Render(m.viewState())
}
