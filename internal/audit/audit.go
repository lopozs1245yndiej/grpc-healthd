// Package audit provides a simple in-memory audit log for recording
// administrative actions performed against the health daemon.
package audit

import (
	"sync"
	"time"
)

// Entry represents a single recorded audit event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	RemoteIP  string    `json:"remote_ip"`
	Action    string    `json:"action"`
	Service   string    `json:"service,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

// Log is a bounded, thread-safe in-memory audit log.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New creates a new Log that retains at most maxSize entries.
// Older entries are evicted when the capacity is exceeded.
func New(maxSize int) *Log {
	if maxSize <= 0 {
		maxSize = 256
	}
	return &Log{maxSize: maxSize}
}

// Record appends a new entry to the log.
func (l *Log) Record(remoteIP, action, service, detail string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		RemoteIP:  remoteIP,
		Action:    action,
		Service:   service,
		Detail:    detail,
	}

	l.entries = append(l.entries, entry)
	if len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}
}

// Entries returns a snapshot of all current log entries, oldest first.
func (l *Log) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	snap := make([]Entry, len(l.entries))
	copy(snap, l.entries)
	return snap
}

// Len returns the number of entries currently held in the log.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}
