package alarm

import (
	"errors"
	"testing"
)

type spyAlarm struct {
	played int
	err    error
}

func (s *spyAlarm) Play() error {
	s.played++
	return s.err
}

func TestNoopPlay(t *testing.T) {
	if err := (Noop{}).Play(); err != nil {
		t.Errorf("Noop.Play = %v, want nil", err)
	}
}

func TestCompositePlaysAll(t *testing.T) {
	a, b := &spyAlarm{}, &spyAlarm{}
	c := NewComposite(a, b)
	if err := c.Play(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.played != 1 || b.played != 1 {
		t.Errorf("plays: a=%d b=%d", a.played, b.played)
	}
}

func TestCompositeSkipsNil(t *testing.T) {
	a := &spyAlarm{}
	c := NewComposite(nil, a, nil)
	if err := c.Play(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.played != 1 {
		t.Errorf("a.played = %d, want 1", a.played)
	}
}

func TestCompositeSwallowsInnerErrors(t *testing.T) {
	a := &spyAlarm{err: errors.New("boom")}
	b := &spyAlarm{}
	c := NewComposite(a, b)
	if err := c.Play(); err != nil {
		t.Errorf("Composite.Play = %v, want nil", err)
	}
	if a.played != 1 || b.played != 1 {
		t.Errorf("both should play even when first fails: a=%d b=%d", a.played, b.played)
	}
}

func TestNewSystemReturnsUsableAlarm(t *testing.T) {
	a := NewSystem()
	if a == nil {
		t.Fatal("NewSystem returned nil")
	}
	if _, ok := a.(Composite); !ok {
		t.Errorf("expected Composite, got %T", a)
	}
}
