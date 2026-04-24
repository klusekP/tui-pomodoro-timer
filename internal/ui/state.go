package ui

type RunState int

const (
	StateIdle RunState = iota
	StateRunning
	StatePaused
	StateFinished
)
