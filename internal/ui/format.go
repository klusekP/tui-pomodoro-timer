package ui

import (
	"fmt"
	"strings"
	"time"
)

func FormatHMS(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int64(d / time.Second)
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func FormatHours(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	totalMin := int64(d / time.Minute)
	h := totalMin / 60
	m := totalMin % 60
	return fmt.Sprintf("%dh %02dm", h, m)
}

func Truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 1 {
		return string(r[:max])
	}
	return string(r[:max-1]) + "…"
}

func DisplayTask(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "Untitled"
	}
	return s
}

func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func ClampDuration(v, lo, hi time.Duration) time.Duration {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
