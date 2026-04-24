package mode

import (
	"strings"
	"testing"
	"time"

	"github.com/klusekP/time_pr_tracking/internal/domain"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

func TestStopwatchInitialState(t *testing.T) {
	s := NewStopwatch(newFakeClock(), &fakeRepo{})
	if s.State() != ui.StateIdle {
		t.Errorf("state = %v, want Idle", s.State())
	}
	if s.Elapsed() != 0 {
		t.Errorf("elapsed = %v, want 0", s.Elapsed())
	}
	if !s.UsesTaskInput() {
		t.Error("stopwatch should use task input")
	}
	if s.CapturingKeys() {
		t.Error("stopwatch should not capture keys")
	}
}

func TestStopwatchStartPauseResumeWithFakeClock(t *testing.T) {
	fc := newFakeClock()
	s := NewStopwatch(fc, &fakeRepo{})

	s.HandleKey(spaceKey(), "")
	if s.State() != ui.StateRunning {
		t.Fatalf("after start: state=%v", s.State())
	}

	fc.advance(40 * time.Second)
	s.HandleKey(spaceKey(), "")
	if s.State() != ui.StatePaused {
		t.Fatalf("after pause: state=%v", s.State())
	}
	if s.Elapsed() != 40*time.Second {
		t.Errorf("elapsed after pause = %v, want 40s", s.Elapsed())
	}

	fc.advance(5 * time.Minute)
	if s.Elapsed() != 40*time.Second {
		t.Errorf("elapsed must not grow while paused: %v", s.Elapsed())
	}

	s.HandleKey(spaceKey(), "")
	fc.advance(20 * time.Second)
	if s.Elapsed() != 60*time.Second {
		t.Errorf("elapsed after resume+20s = %v, want 60s", s.Elapsed())
	}
}

func TestStopwatchStopAndSaveProducesSavedMsg(t *testing.T) {
	fc := newFakeClock()
	repo := &fakeRepo{}
	s := NewStopwatch(fc, repo)

	s.HandleKey(spaceKey(), "")
	fc.advance(2 * time.Minute)

	cmd := s.HandleKey(runeKey('s'), "  Login flow  ")
	if cmd == nil {
		t.Fatal("expected non-nil save cmd")
	}
	msg, ok := cmd().(ui.SavedMsg)
	if !ok {
		t.Fatalf("expected SavedMsg, got %T", cmd())
	}
	if msg.Err != nil {
		t.Fatalf("save error: %v", msg.Err)
	}
	if msg.Session.Task != "Login flow" {
		t.Errorf("task = %q, want %q (trimmed)", msg.Session.Task, "Login flow")
	}
	if msg.Session.Elapsed != 2*time.Minute {
		t.Errorf("elapsed = %v, want 2m", msg.Session.Elapsed)
	}
	if msg.Session.Kind != domain.KindStopwatch {
		t.Errorf("kind = %v, want stopwatch", msg.Session.Kind)
	}
	if msg.Session.ID == 0 {
		t.Error("saved session should have an ID")
	}
	if len(repo.saved) != 1 {
		t.Errorf("repo.saved = %d, want 1", len(repo.saved))
	}
	if s.State() != ui.StateIdle {
		t.Errorf("state after save = %v, want Idle", s.State())
	}
}

func TestStopwatchStopIdleEmitsToastAndSkipsRepo(t *testing.T) {
	repo := &fakeRepo{}
	s := NewStopwatch(newFakeClock(), repo)

	cmd := s.HandleKey(runeKey('s'), "irrelevant")
	if cmd == nil {
		t.Fatal("expected toast cmd")
	}
	msg, ok := cmd().(ui.ToastMsg)
	if !ok {
		t.Fatalf("expected ToastMsg, got %T", cmd())
	}
	if !strings.Contains(string(msg), "Nothing") {
		t.Errorf("toast text = %q", string(msg))
	}
	if len(repo.saved) != 0 {
		t.Errorf("repo should be untouched, got %d saves", len(repo.saved))
	}
}

func TestStopwatchResetKey(t *testing.T) {
	fc := newFakeClock()
	s := NewStopwatch(fc, &fakeRepo{})
	s.HandleKey(spaceKey(), "")
	fc.advance(time.Minute)

	cmd := s.HandleKey(runeKey('r'), "")
	if cmd == nil {
		t.Fatal("expected toast cmd after reset")
	}
	if _, ok := cmd().(ui.ToastMsg); !ok {
		t.Errorf("expected ToastMsg, got %T", cmd())
	}
	if s.State() != ui.StateIdle {
		t.Errorf("state after reset = %v, want Idle", s.State())
	}
	if s.Elapsed() != 0 {
		t.Errorf("elapsed after reset = %v, want 0", s.Elapsed())
	}
}

func TestStopwatchSaveErrorIsPropagated(t *testing.T) {
	repo := &fakeRepo{saveErr: errBoom}
	s := NewStopwatch(newFakeClock(), repo)
	s.HandleKey(spaceKey(), "")

	cmd := s.HandleKey(runeKey('s'), "task")
	msg, ok := cmd().(ui.SavedMsg)
	if !ok {
		t.Fatalf("expected SavedMsg, got %T", cmd())
	}
	if msg.Err == nil {
		t.Error("expected error from save to propagate in SavedMsg.Err")
	}
}

var errBoom = errTest("boom")

type errTest string

func (e errTest) Error() string { return string(e) }
