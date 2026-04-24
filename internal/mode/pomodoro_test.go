package mode

import (
	"testing"
	"time"

	"github.com/klusekP/time_pr_tracking/internal/alarm"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

func TestPomodoroInitialState(t *testing.T) {
	p := NewPomodoro(newFakeClock(), alarm.Noop{}, 25*time.Minute)

	if p.State() != ui.StateIdle {
		t.Errorf("state = %v, want Idle", p.State())
	}
	if p.Target() != 25*time.Minute {
		t.Errorf("target = %v, want 25m", p.Target())
	}
	if p.Elapsed() != 0 {
		t.Errorf("elapsed = %v, want 0", p.Elapsed())
	}
	if p.Remaining() != 25*time.Minute {
		t.Errorf("remaining = %v, want 25m", p.Remaining())
	}
	if p.Ratio() != 0 {
		t.Errorf("ratio = %v, want 0", p.Ratio())
	}
	if p.UsesTaskInput() {
		t.Error("pomodoro should not use task input")
	}
	if p.CapturingKeys() {
		t.Error("pomodoro should not capture keys")
	}
}

func TestPomodoroTargetClampsBelowMin(t *testing.T) {
	p := NewPomodoro(newFakeClock(), nil, time.Second)
	if p.Target() != time.Minute {
		t.Errorf("target = %v, want 1m (min)", p.Target())
	}
}

func TestPomodoroTargetClampsAboveMax(t *testing.T) {
	p := NewPomodoro(newFakeClock(), nil, 100*time.Hour)
	if p.Target() != 8*time.Hour {
		t.Errorf("target = %v, want 8h (max)", p.Target())
	}
}

func TestPomodoroNilAlarmFallsBackToNoop(t *testing.T) {
	p := NewPomodoro(newFakeClock(), nil, 25*time.Minute)
	if _, ok := p.alarm.(alarm.Noop); !ok {
		t.Errorf("alarm = %T, want alarm.Noop", p.alarm)
	}
}

func TestPomodoroStartAndTickBeforeTarget(t *testing.T) {
	fc := newFakeClock()
	p := NewPomodoro(fc, alarm.Noop{}, time.Minute)

	p.HandleKey(spaceKey(), "")
	if p.State() != ui.StateRunning {
		t.Fatalf("state = %v, want Running", p.State())
	}

	fc.advance(30 * time.Second)
	if cmd := p.OnTick(); cmd != nil {
		t.Error("OnTick before target should be nil")
	}
	if p.State() != ui.StateRunning {
		t.Errorf("state = %v, want still Running", p.State())
	}
}

func TestPomodoroFinishesViaOnTick(t *testing.T) {
	fc := newFakeClock()
	p := NewPomodoro(fc, alarm.Noop{}, time.Minute)
	p.HandleKey(spaceKey(), "")

	fc.advance(61 * time.Second)
	cmd := p.OnTick()
	if cmd == nil {
		t.Fatal("OnTick past target should return alarm+toast cmd")
	}

	if p.State() != ui.StateFinished {
		t.Errorf("state = %v, want Finished", p.State())
	}
	if p.Elapsed() != time.Minute {
		t.Errorf("elapsed = %v, want clamped to target (1m)", p.Elapsed())
	}
	if p.Remaining() != 0 {
		t.Errorf("remaining = %v, want 0", p.Remaining())
	}
	if p.Ratio() != 1 {
		t.Errorf("ratio = %v, want 1", p.Ratio())
	}
}

func TestPomodoroOnTickNoopAfterFinish(t *testing.T) {
	fc := newFakeClock()
	p := NewPomodoro(fc, alarm.Noop{}, time.Minute)
	p.HandleKey(spaceKey(), "")
	fc.advance(61 * time.Second)
	p.OnTick()

	if cmd := p.OnTick(); cmd != nil {
		t.Error("OnTick after finish should be nil")
	}
}

func TestPomodoroSpaceAfterFinishRestartsFresh(t *testing.T) {
	fc := newFakeClock()
	p := NewPomodoro(fc, alarm.Noop{}, time.Minute)
	p.HandleKey(spaceKey(), "")
	fc.advance(61 * time.Second)
	p.OnTick()

	p.HandleKey(spaceKey(), "")
	if p.State() != ui.StateRunning {
		t.Errorf("state = %v, want Running", p.State())
	}
	if p.Elapsed() != 0 {
		t.Errorf("elapsed after restart = %v, want 0", p.Elapsed())
	}
}

func TestPomodoroResetKey(t *testing.T) {
	fc := newFakeClock()
	p := NewPomodoro(fc, alarm.Noop{}, time.Minute)
	p.HandleKey(spaceKey(), "")
	fc.advance(30 * time.Second)

	p.HandleKey(runeKey('r'), "")
	if p.State() != ui.StateIdle {
		t.Errorf("state = %v, want Idle", p.State())
	}
	if p.Elapsed() != 0 {
		t.Errorf("elapsed = %v, want 0", p.Elapsed())
	}
}

func TestPomodoroStopIdleEmitsToast(t *testing.T) {
	p := NewPomodoro(newFakeClock(), alarm.Noop{}, time.Minute)
	cmd := p.HandleKey(runeKey('s'), "")
	if cmd == nil {
		t.Fatal("expected toast cmd")
	}
	if _, ok := cmd().(ui.ToastMsg); !ok {
		t.Errorf("expected ToastMsg, got %T", cmd())
	}
}

func TestPomodoroPresetKeysSetTarget(t *testing.T) {
	cases := []struct {
		key  rune
		want time.Duration
	}{
		{'w', 25 * time.Minute},
		{'b', 5 * time.Minute},
		{'l', 15 * time.Minute},
	}
	for _, c := range cases {
		p := NewPomodoro(newFakeClock(), alarm.Noop{}, time.Hour)
		p.HandleKey(runeKey(c.key), "")
		if p.Target() != c.want {
			t.Errorf("key %q → target %v, want %v", string(c.key), p.Target(), c.want)
		}
	}
}

func TestPomodoroPlusMinusAdjustsTarget(t *testing.T) {
	p := NewPomodoro(newFakeClock(), alarm.Noop{}, 20*time.Minute)
	p.HandleKey(runeKey('+'), "")
	if p.Target() != 25*time.Minute {
		t.Errorf("after +: %v, want 25m", p.Target())
	}
	p.HandleKey(runeKey('-'), "")
	p.HandleKey(runeKey('-'), "")
	if p.Target() != 15*time.Minute {
		t.Errorf("after two -: %v, want 15m", p.Target())
	}
}

func TestPomodoroPlusDoesNotExceedMax(t *testing.T) {
	p := NewPomodoro(newFakeClock(), alarm.Noop{}, 8*time.Hour)
	p.HandleKey(runeKey('+'), "")
	if p.Target() != 8*time.Hour {
		t.Errorf("target should stay at max: %v", p.Target())
	}
}

func TestPomodoroMinusDoesNotGoBelowMin(t *testing.T) {
	p := NewPomodoro(newFakeClock(), alarm.Noop{}, time.Minute)
	p.HandleKey(runeKey('-'), "")
	if p.Target() != time.Minute {
		t.Errorf("target should stay at min: %v", p.Target())
	}
}
