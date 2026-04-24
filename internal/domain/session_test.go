package domain

import (
	"testing"
	"time"
)

func TestSessionKindLabel(t *testing.T) {
	cases := []struct {
		kind SessionKind
		want string
	}{
		{KindStopwatch, "Stopwatch"},
		{KindPomodoro, "Pomodoro"},
		{"", "Stopwatch"},
		{"custom-kind", "custom-kind"},
	}
	for _, c := range cases {
		if got := c.kind.Label(); got != c.want {
			t.Errorf("Label(%q) = %q, want %q", string(c.kind), got, c.want)
		}
	}
}

func TestSessionDuration(t *testing.T) {
	s := Session{Elapsed: 5 * time.Minute}
	if got := s.Duration(); got != 5*time.Minute {
		t.Errorf("Duration = %v, want %v", got, 5*time.Minute)
	}
}

func TestSessionDurationZero(t *testing.T) {
	s := Session{}
	if got := s.Duration(); got != 0 {
		t.Errorf("Duration = %v, want 0", got)
	}
}
