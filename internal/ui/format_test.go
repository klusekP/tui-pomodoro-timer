package ui

import (
	"testing"
	"time"
)

func TestFormatHMS(t *testing.T) {
	cases := []struct {
		name string
		in   time.Duration
		want string
	}{
		{"zero", 0, "00:00:00"},
		{"seconds", 45 * time.Second, "00:00:45"},
		{"minutes+seconds", 2*time.Minute + 5*time.Second, "00:02:05"},
		{"hours+minutes+seconds", 3*time.Hour + 12*time.Minute + 9*time.Second, "03:12:09"},
		{"over 24h still counted", 25 * time.Hour, "25:00:00"},
		{"negative clamped to zero", -5 * time.Second, "00:00:00"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := FormatHMS(c.in); got != c.want {
				t.Errorf("FormatHMS(%v) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestFormatHours(t *testing.T) {
	cases := []struct {
		name string
		in   time.Duration
		want string
	}{
		{"zero", 0, "0h 00m"},
		{"sub-minute floors to zero", 30 * time.Second, "0h 00m"},
		{"five minutes", 5 * time.Minute, "0h 05m"},
		{"90 minutes", 90 * time.Minute, "1h 30m"},
		{"exactly 5h", 5 * time.Hour, "5h 00m"},
		{"negative clamped", -time.Hour, "0h 00m"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := FormatHours(c.in); got != c.want {
				t.Errorf("FormatHours(%v) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		s    string
		max  int
		want string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello world", 8, "hello w…"},
		{"hello", 0, ""},
		{"hello", 1, "h"},
		{"żółw", 2, "ż…"},
		{"żółw", 4, "żółw"},
		{"", 5, ""},
	}
	for _, c := range cases {
		if got := Truncate(c.s, c.max); got != c.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", c.s, c.max, got, c.want)
		}
	}
}

func TestDisplayTask(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "Untitled"},
		{"   ", "Untitled"},
		{"hello", "hello"},
		{"  hello  ", "hello"},
		{"\t\nhi\t", "hi"},
	}
	for _, c := range cases {
		if got := DisplayTask(c.in); got != c.want {
			t.Errorf("DisplayTask(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestClamp(t *testing.T) {
	if got := Clamp(5, 0, 10); got != 5 {
		t.Errorf("middle: got %d", got)
	}
	if got := Clamp(-3, 0, 10); got != 0 {
		t.Errorf("below lo: got %d", got)
	}
	if got := Clamp(50, 0, 10); got != 10 {
		t.Errorf("above hi: got %d", got)
	}
	if got := Clamp(0, 0, 10); got != 0 {
		t.Errorf("boundary lo: got %d", got)
	}
	if got := Clamp(10, 0, 10); got != 10 {
		t.Errorf("boundary hi: got %d", got)
	}
}

func TestClampDuration(t *testing.T) {
	lo := time.Minute
	hi := time.Hour
	if got := ClampDuration(30*time.Minute, lo, hi); got != 30*time.Minute {
		t.Errorf("middle: got %v", got)
	}
	if got := ClampDuration(0, lo, hi); got != lo {
		t.Errorf("below lo: got %v", got)
	}
	if got := ClampDuration(2*time.Hour, lo, hi); got != hi {
		t.Errorf("above hi: got %v", got)
	}
}
