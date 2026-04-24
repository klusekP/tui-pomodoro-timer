package timer

import (
	"testing"
	"time"
)

type fakeClock struct {
	now time.Time
}

func (f *fakeClock) Now() time.Time { return f.now }

func (f *fakeClock) advance(d time.Duration) { f.now = f.now.Add(d) }

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}
}

func TestNewTimerIsIdle(t *testing.T) {
	tm := New(newFakeClock())
	if tm.State() != Idle {
		t.Errorf("state = %v, want Idle", tm.State())
	}
	if !tm.IsIdle() {
		t.Error("IsIdle should be true")
	}
	if tm.Elapsed() != 0 {
		t.Errorf("Elapsed = %v, want 0", tm.Elapsed())
	}
}

func TestNewTimerWithNilClockFallsBackToSystem(t *testing.T) {
	tm := New(nil)
	if tm == nil {
		t.Fatal("nil timer")
	}
	if tm.clock == nil {
		t.Error("clock should fall back to System")
	}
}

func TestStartAndElapsedWhileRunning(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)
	tm.Start()

	if tm.State() != Running {
		t.Errorf("state = %v, want Running", tm.State())
	}
	if tm.SessionStart().IsZero() {
		t.Error("sessionStart should be set")
	}

	fc.advance(3 * time.Second)
	if got := tm.Elapsed(); got != 3*time.Second {
		t.Errorf("Elapsed = %v, want 3s", got)
	}
}

func TestPauseFreezesElapsed(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)
	tm.Start()
	fc.advance(10 * time.Second)
	tm.Pause()

	if tm.State() != Paused {
		t.Errorf("state = %v, want Paused", tm.State())
	}
	fc.advance(1 * time.Hour)
	if got := tm.Elapsed(); got != 10*time.Second {
		t.Errorf("Elapsed after pause = %v, want 10s", got)
	}
}

func TestResumeContinuesFromPause(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)
	tm.Start()
	fc.advance(5 * time.Second)
	tm.Pause()
	fc.advance(10 * time.Second)
	tm.Start()
	fc.advance(2 * time.Second)

	if got := tm.Elapsed(); got != 7*time.Second {
		t.Errorf("Elapsed = %v, want 7s", got)
	}
}

func TestToggleCycles(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)

	tm.Toggle()
	if tm.State() != Running {
		t.Errorf("first toggle: %v", tm.State())
	}
	tm.Toggle()
	if tm.State() != Paused {
		t.Errorf("second toggle: %v", tm.State())
	}
	tm.Toggle()
	if tm.State() != Running {
		t.Errorf("third toggle: %v", tm.State())
	}
}

func TestStopReturnsTimesAndResets(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)
	tm.Start()
	expectedStart := fc.now
	fc.advance(90 * time.Second)

	start, end, elapsed := tm.Stop()

	if !start.Equal(expectedStart) {
		t.Errorf("start = %v, want %v", start, expectedStart)
	}
	if !end.Equal(fc.now) {
		t.Errorf("end = %v, want %v", end, fc.now)
	}
	if elapsed != 90*time.Second {
		t.Errorf("elapsed = %v, want 90s", elapsed)
	}
	if tm.State() != Idle {
		t.Errorf("state after stop = %v, want Idle", tm.State())
	}
	if tm.Elapsed() != 0 {
		t.Errorf("Elapsed after stop = %v, want 0", tm.Elapsed())
	}
}

func TestStopWhilePausedReturnsAccumulated(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)
	tm.Start()
	fc.advance(7 * time.Second)
	tm.Pause()
	fc.advance(10 * time.Minute)

	_, _, elapsed := tm.Stop()
	if elapsed != 7*time.Second {
		t.Errorf("elapsed = %v, want 7s", elapsed)
	}
}

func TestResetFromAnyStateGoesIdle(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)
	tm.Start()
	fc.advance(time.Minute)
	tm.Reset()

	if tm.State() != Idle {
		t.Errorf("state = %v, want Idle", tm.State())
	}
	if tm.Elapsed() != 0 {
		t.Errorf("Elapsed = %v, want 0", tm.Elapsed())
	}
	if !tm.SessionStart().IsZero() {
		t.Error("sessionStart should be zero after reset")
	}
}

func TestPauseWhileIdleIsNoop(t *testing.T) {
	fc := newFakeClock()
	tm := New(fc)
	tm.Pause()
	if tm.State() != Idle {
		t.Errorf("state = %v, want Idle", tm.State())
	}
}
