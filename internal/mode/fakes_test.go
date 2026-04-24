package mode

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/klusekP/time_pr_tracking/internal/domain"
)

type fakeClock struct {
	now time.Time
}

func (f *fakeClock) Now() time.Time { return f.now }

func (f *fakeClock) advance(d time.Duration) { f.now = f.now.Add(d) }

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)}
}

type fakeRepo struct {
	saved   []domain.Session
	saveErr error
}

func (r *fakeRepo) Save(s domain.Session) (domain.Session, error) {
	if r.saveErr != nil {
		return s, r.saveErr
	}
	s.ID = int64(len(r.saved) + 1)
	r.saved = append(r.saved, s)
	return s, nil
}

func (r *fakeRepo) Recent(int) ([]domain.Session, error) { return nil, nil }
func (r *fakeRepo) TotalToday() (time.Duration, error)   { return 0, nil }
func (r *fakeRepo) Close() error                         { return nil }

func runeKey(ch rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
}

func spaceKey() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
}

func enterKey() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyEnter}
}
