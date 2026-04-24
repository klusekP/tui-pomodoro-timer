package timer

import (
	"time"

	"github.com/klusekP/time_pr_tracking/internal/clock"
)

type State int

const (
	Idle State = iota
	Running
	Paused
)

type Timer struct {
	clock        clock.Clock
	state        State
	startedAt    time.Time
	accumulated  time.Duration
	sessionStart time.Time
}

func New(c clock.Clock) *Timer {
	if c == nil {
		c = clock.System{}
	}
	return &Timer{clock: c}
}

func (t *Timer) State() State { return t.state }

func (t *Timer) IsIdle() bool { return t.state == Idle }

func (t *Timer) SessionStart() time.Time { return t.sessionStart }

func (t *Timer) Start() {
	now := t.clock.Now()
	switch t.state {
	case Idle:
		t.accumulated = 0
		t.sessionStart = now
		t.startedAt = now
		t.state = Running
	case Paused:
		t.startedAt = now
		t.state = Running
	}
}

func (t *Timer) Pause() {
	if t.state != Running {
		return
	}
	t.accumulated += t.clock.Now().Sub(t.startedAt)
	t.state = Paused
}

func (t *Timer) Toggle() {
	if t.state == Running {
		t.Pause()
	} else {
		t.Start()
	}
}

func (t *Timer) Elapsed() time.Duration {
	switch t.state {
	case Running:
		return t.accumulated + t.clock.Now().Sub(t.startedAt)
	case Paused:
		return t.accumulated
	default:
		return 0
	}
}

func (t *Timer) Stop() (start, end time.Time, elapsed time.Duration) {
	end = t.clock.Now()
	if t.state == Running {
		t.accumulated += end.Sub(t.startedAt)
	}
	start = t.sessionStart
	elapsed = t.accumulated
	t.Reset()
	return
}

func (t *Timer) Reset() {
	t.state = Idle
	t.startedAt = time.Time{}
	t.accumulated = 0
	t.sessionStart = time.Time{}
}
