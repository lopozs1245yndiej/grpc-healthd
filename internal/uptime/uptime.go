// Package uptime tracks the daemon start time and exposes uptime metrics.
package uptime

import (
	"sync"
	"time"
)

// Tracker records when the daemon started and provides uptime calculations.
type Tracker struct {
	mu        sync.RWMutex
	startTime time.Time
}

// New creates a new Tracker with the current UTC time as the start time.
func New() *Tracker {
	return &Tracker{
		startTime: time.Now().UTC(),
	}
}

// NewWithTime creates a Tracker with an explicit start time (useful for testing).
func NewWithTime(t time.Time) *Tracker {
	return &Tracker{
		startTime: t.UTC(),
	}
}

// StartTime returns the UTC time at which the daemon started.
func (t *Tracker) StartTime() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.startTime
}

// Uptime returns the duration since the daemon started.
func (t *Tracker) Uptime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return time.Since(t.startTime).Truncate(time.Second)
}

// Info holds a snapshot of uptime information for serialisation.
type Info struct {
	StartTime string `json:"start_time"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

// Snapshot returns a point-in-time Info struct.
func (t *Tracker) Snapshot() Info {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return Info{
		StartTime:     t.startTime.Format(time.RFC3339),
		UptimeSeconds: int64(time.Since(t.startTime).Seconds()),
	}
}
