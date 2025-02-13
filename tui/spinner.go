package tui

import "time"

type Spinner struct {
	Frames []string
	FPS    time.Duration
}

type TickMsg struct {
	time time.Time
	tag  int
	ID   int
}
