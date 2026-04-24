package mode

import (
	"strings"
	"testing"
	"time"

	"github.com/klusekP/time_pr_tracking/internal/domain"
)

func sampleSessions() []domain.Session {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	return []domain.Session{
		{Kind: domain.KindStopwatch, Task: "frontend login", Elapsed: 90 * time.Minute, EndedAt: base},
		{Kind: domain.KindStopwatch, Task: "backend api", Elapsed: 45 * time.Minute, EndedAt: base.Add(time.Hour)},
		{Kind: domain.KindPomodoro, Task: "frontend styles", Elapsed: 25 * time.Minute, EndedAt: base.Add(2 * time.Hour)},
		{Kind: domain.KindPomodoro, Task: "meeting", Elapsed: 30 * time.Minute, EndedAt: base.Add(3 * time.Hour)},
	}
}

func TestSummarize(t *testing.T) {
	total, byKind := summarize(sampleSessions())

	want := 190 * time.Minute
	if total != want {
		t.Errorf("total = %v, want %v", total, want)
	}
	if got := byKind[domain.KindStopwatch]; got.count != 2 || got.total != 135*time.Minute {
		t.Errorf("stopwatch = %+v, want count=2 total=2h15m", got)
	}
	if got := byKind[domain.KindPomodoro]; got.count != 2 || got.total != 55*time.Minute {
		t.Errorf("pomodoro = %+v, want count=2 total=55m", got)
	}
}

func TestSummarizeEmpty(t *testing.T) {
	total, byKind := summarize(nil)
	if total != 0 {
		t.Errorf("total = %v, want 0", total)
	}
	if len(byKind) != 0 {
		t.Errorf("byKind = %v, want empty", byKind)
	}
}

func TestHistoryFilterEmpty(t *testing.T) {
	h := NewHistory()
	h.SetSessions(sampleSessions(), 0)
	if got := len(h.filtered()); got != 4 {
		t.Errorf("empty filter should return all: got %d", got)
	}
}

func TestHistoryFilterSubstring(t *testing.T) {
	h := NewHistory()
	h.SetSessions(sampleSessions(), 0)
	h.filter.SetValue("frontend")

	got := h.filtered()
	if len(got) != 2 {
		t.Fatalf("len(filtered) = %d, want 2", len(got))
	}
	for _, s := range got {
		if !strings.Contains(strings.ToLower(s.Task), "frontend") {
			t.Errorf("unexpected session in results: %q", s.Task)
		}
	}
}

func TestHistoryFilterCaseInsensitive(t *testing.T) {
	h := NewHistory()
	h.SetSessions(sampleSessions(), 0)
	h.filter.SetValue("FRONTEND")
	if got := len(h.filtered()); got != 2 {
		t.Errorf("case-insensitive: got %d, want 2", got)
	}
}

func TestHistoryFilterNoMatches(t *testing.T) {
	h := NewHistory()
	h.SetSessions(sampleSessions(), 0)
	h.filter.SetValue("does-not-exist")
	if got := h.filtered(); len(got) != 0 {
		t.Errorf("want empty, got %d", len(got))
	}
}

func TestHistoryFilterTrimsWhitespace(t *testing.T) {
	h := NewHistory()
	h.SetSessions(sampleSessions(), 0)
	h.filter.SetValue("   ")
	if got := len(h.filtered()); got != 4 {
		t.Errorf("whitespace-only filter should be treated as empty: got %d", got)
	}
}

func TestHistoryCapturingKeysToggle(t *testing.T) {
	h := NewHistory()
	if h.CapturingKeys() {
		t.Fatal("fresh history should not capture keys")
	}

	h.HandleKey(runeKey('f'), "")
	if !h.CapturingKeys() {
		t.Error("f should enter edit mode")
	}

	h.HandleKey(enterKey(), "")
	if h.CapturingKeys() {
		t.Error("enter should exit edit mode")
	}
}

func TestHistoryFilterInputCapturesTyping(t *testing.T) {
	h := NewHistory()
	h.HandleKey(runeKey('f'), "")
	h.HandleKey(runeKey('h'), "")
	h.HandleKey(runeKey('i'), "")

	if got := h.filter.Value(); got != "hi" {
		t.Errorf("filter value = %q, want %q", got, "hi")
	}
}

func TestHistoryClearFilter(t *testing.T) {
	h := NewHistory()
	h.filter.SetValue("xyz")

	cmd := h.HandleKey(runeKey('c'), "")
	if h.filter.Value() != "" {
		t.Errorf("filter not cleared: %q", h.filter.Value())
	}
	if cmd == nil {
		t.Error("expected toast cmd after clear")
	}
}

func TestHistoryClearNoopWhenEmpty(t *testing.T) {
	h := NewHistory()
	if cmd := h.HandleKey(runeKey('c'), ""); cmd != nil {
		t.Error("clearing empty filter should be a no-op (nil cmd)")
	}
}

func TestHistorySetSessionsReplacesState(t *testing.T) {
	h := NewHistory()
	h.SetSessions(sampleSessions(), 2*time.Hour)
	if len(h.all) != 4 {
		t.Errorf("all = %d, want 4", len(h.all))
	}
	if h.totalToday != 2*time.Hour {
		t.Errorf("totalToday = %v, want 2h", h.totalToday)
	}

	h.SetSessions(nil, 0)
	if len(h.all) != 0 {
		t.Errorf("all after empty = %d, want 0", len(h.all))
	}
	if h.totalToday != 0 {
		t.Errorf("totalToday after empty = %v, want 0", h.totalToday)
	}
}
