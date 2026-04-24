package domain

import "time"

type SessionKind string

const (
	KindStopwatch SessionKind = "stopwatch"
	KindPomodoro  SessionKind = "pomodoro"
)

func (k SessionKind) Label() string {
	switch k {
	case KindPomodoro:
		return "Pomodoro"
	case KindStopwatch, "":
		return "Stopwatch"
	default:
		return string(k)
	}
}

type Session struct {
	ID        int64
	Kind      SessionKind
	Task      string
	Elapsed   time.Duration
	StartedAt time.Time
	EndedAt   time.Time
}

func (s Session) Duration() time.Duration { return s.Elapsed }
