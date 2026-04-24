package repository

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/klusekP/time_pr_tracking/internal/domain"
)

func newTestRepo(t *testing.T) *SQLite {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	r, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("NewSQLite: %v", err)
	}
	t.Cleanup(func() { _ = r.Close() })
	return r
}

func TestSQLiteSaveRoundtrip(t *testing.T) {
	r := newTestRepo(t)

	now := time.Now()
	in := domain.Session{
		Kind:      domain.KindStopwatch,
		Task:      "refactor",
		Elapsed:   5 * time.Minute,
		StartedAt: now.Add(-5 * time.Minute),
		EndedAt:   now,
	}
	saved, err := r.Save(in)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if saved.ID == 0 {
		t.Error("saved.ID should be non-zero")
	}

	rec, err := r.Recent(10)
	if err != nil {
		t.Fatalf("Recent: %v", err)
	}
	if len(rec) != 1 {
		t.Fatalf("len(Recent) = %d, want 1", len(rec))
	}
	got := rec[0]
	if got.Task != "refactor" {
		t.Errorf("task = %q", got.Task)
	}
	if got.Kind != domain.KindStopwatch {
		t.Errorf("kind = %q", got.Kind)
	}
	if got.Elapsed != 5*time.Minute {
		t.Errorf("elapsed = %v", got.Elapsed)
	}
	if got.ID != saved.ID {
		t.Errorf("id mismatch: got %d, saved %d", got.ID, saved.ID)
	}
}

func TestSQLiteRecentOrderedByEndedAtDesc(t *testing.T) {
	r := newTestRepo(t)
	base := time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC)

	for i := 0; i < 3; i++ {
		ended := base.Add(time.Duration(i) * time.Hour)
		_, err := r.Save(domain.Session{
			Kind:      domain.KindStopwatch,
			Task:      fmt.Sprintf("task-%d", i),
			Elapsed:   time.Minute,
			StartedAt: ended.Add(-time.Minute),
			EndedAt:   ended,
		})
		if err != nil {
			t.Fatalf("Save %d: %v", i, err)
		}
	}

	rec, err := r.Recent(10)
	if err != nil {
		t.Fatalf("Recent: %v", err)
	}
	if len(rec) != 3 {
		t.Fatalf("len = %d, want 3", len(rec))
	}
	if rec[0].Task != "task-2" || rec[1].Task != "task-1" || rec[2].Task != "task-0" {
		t.Errorf("unexpected order: %q, %q, %q", rec[0].Task, rec[1].Task, rec[2].Task)
	}
}

func TestSQLiteRecentRespectsLimit(t *testing.T) {
	r := newTestRepo(t)
	base := time.Now()
	for i := 0; i < 5; i++ {
		_, _ = r.Save(domain.Session{
			Kind:      domain.KindStopwatch,
			Task:      fmt.Sprintf("t%d", i),
			Elapsed:   time.Minute,
			StartedAt: base,
			EndedAt:   base.Add(time.Duration(i) * time.Second),
		})
	}
	rec, err := r.Recent(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec) != 2 {
		t.Errorf("len = %d, want 2", len(rec))
	}
}

func TestSQLiteTotalTodaySumsOnlyToday(t *testing.T) {
	r := newTestRepo(t)

	now := time.Now()
	yesterday := now.Add(-25 * time.Hour)

	_, _ = r.Save(domain.Session{Kind: domain.KindStopwatch, Elapsed: 30 * time.Minute, StartedAt: now, EndedAt: now})
	_, _ = r.Save(domain.Session{Kind: domain.KindStopwatch, Elapsed: 45 * time.Minute, StartedAt: now, EndedAt: now})
	_, _ = r.Save(domain.Session{Kind: domain.KindPomodoro, Elapsed: 2 * time.Hour, StartedAt: yesterday, EndedAt: yesterday})

	total, err := r.TotalToday()
	if err != nil {
		t.Fatal(err)
	}
	if total != 75*time.Minute {
		t.Errorf("total = %v, want 1h15m", total)
	}
}

func TestSQLiteTotalTodayEmptyReturnsZero(t *testing.T) {
	r := newTestRepo(t)
	total, err := r.TotalToday()
	if err != nil {
		t.Fatal(err)
	}
	if total != 0 {
		t.Errorf("total = %v, want 0", total)
	}
}

func TestSQLiteSaveDefaultsKindToStopwatch(t *testing.T) {
	r := newTestRepo(t)
	now := time.Now()
	saved, err := r.Save(domain.Session{Task: "no-kind", Elapsed: time.Minute, StartedAt: now, EndedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Kind != domain.KindStopwatch {
		t.Errorf("saved.Kind = %q, want default stopwatch", saved.Kind)
	}
	rec, _ := r.Recent(1)
	if rec[0].Kind != domain.KindStopwatch {
		t.Errorf("loaded kind = %q, want stopwatch", rec[0].Kind)
	}
}

func TestSQLiteCloseNilSafe(t *testing.T) {
	var r *SQLite
	if err := r.Close(); err != nil {
		t.Errorf("Close on nil = %v, want nil", err)
	}
}
