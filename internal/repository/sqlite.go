package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/klusekP/time_pr_tracking/internal/domain"
)

type Sessions interface {
	Save(session domain.Session) (domain.Session, error)
	Recent(limit int) ([]domain.Session, error)
	TotalToday() (time.Duration, error)
	Close() error
}

type SQLite struct {
	db *sql.DB
}

func NewSQLite(dbPath string) (*SQLite, error) {
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		dir := filepath.Join(home, ".time_tracker")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("cannot create directory %s: %w", dir, err)
		}
		dbPath = filepath.Join(dir, "sessions.db")
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}

	repo := &SQLite{db: db}
	if err := repo.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *SQLite) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		kind         TEXT    NOT NULL DEFAULT 'stopwatch',
		task         TEXT    NOT NULL,
		duration_ms  INTEGER NOT NULL,
		started_at   DATETIME NOT NULL,
		ended_at     DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_sessions_ended_at ON sessions(ended_at DESC);
	`
	if _, err := r.db.Exec(schema); err != nil {
		return fmt.Errorf("cannot initialize schema: %w", err)
	}

	if _, err := r.db.Exec(
		`ALTER TABLE sessions ADD COLUMN kind TEXT NOT NULL DEFAULT 'stopwatch'`,
	); err != nil && !strings.Contains(err.Error(), "duplicate column") {
		return fmt.Errorf("kind column migration failed: %w", err)
	}
	return nil
}

func (r *SQLite) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}

func (r *SQLite) Save(session domain.Session) (domain.Session, error) {
	kind := session.Kind
	if kind == "" {
		kind = domain.KindStopwatch
	}
	res, err := r.db.Exec(
		`INSERT INTO sessions (kind, task, duration_ms, started_at, ended_at) VALUES (?, ?, ?, ?, ?)`,
		string(kind),
		session.Task,
		session.Elapsed.Milliseconds(),
		session.StartedAt.UTC(),
		session.EndedAt.UTC(),
	)
	if err != nil {
		return session, fmt.Errorf("saving session failed: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return session, err
	}
	session.ID = id
	session.Kind = kind
	return session, nil
}

func (r *SQLite) Recent(limit int) ([]domain.Session, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.db.Query(
		`SELECT id, kind, task, duration_ms, started_at, ended_at
		 FROM sessions ORDER BY ended_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("fetching sessions failed: %w", err)
	}
	defer rows.Close()

	var out []domain.Session
	for rows.Next() {
		var (
			s    domain.Session
			kind string
			ms   int64
		)
		if err := rows.Scan(&s.ID, &kind, &s.Task, &ms, &s.StartedAt, &s.EndedAt); err != nil {
			return nil, err
		}
		s.Kind = domain.SessionKind(kind)
		s.Elapsed = time.Duration(ms) * time.Millisecond
		s.StartedAt = s.StartedAt.Local()
		s.EndedAt = s.EndedAt.Local()
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SQLite) TotalToday() (time.Duration, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	var total sql.NullInt64
	err := r.db.QueryRow(
		`SELECT COALESCE(SUM(duration_ms), 0) FROM sessions WHERE ended_at >= ? AND ended_at < ?`,
		start.UTC(), end.UTC(),
	).Scan(&total)
	if err != nil {
		return 0, err
	}
	return time.Duration(total.Int64) * time.Millisecond, nil
}
