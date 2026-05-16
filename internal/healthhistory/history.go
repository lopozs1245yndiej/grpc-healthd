// Package healthhistory tracks status transition events for monitored services.
package healthhistory

import (
	"sync"
	"time"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// Event records a single status transition for a service.
type Event struct {
	Service   string                          `json:"service"`
	Status    grpc_health_v1.HealthCheckResponse_ServingStatus `json:"status"`
	Timestamp time.Time                       `json:"timestamp"`
}

// History maintains a bounded ring-buffer of status transition events.
type History struct {
	mu      sync.Mutex
	events  []Event
	maxSize int
}

// New returns a History that retains at most maxSize events.
// If maxSize is less than 1 it defaults to 100.
func New(maxSize int) *History {
	if maxSize < 1 {
		maxSize = 100
	}
	return &History{
		events:  make([]Event, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record appends a new status-transition event. The oldest event is evicted
// once the buffer is full.
func (h *History) Record(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ev := Event{
		Service:   service,
		Status:    status,
		Timestamp: time.Now().UTC(),
	}

	if len(h.events) >= h.maxSize {
		// Evict oldest by shifting left.
		copy(h.events, h.events[1:])
		h.events[len(h.events)-1] = ev
		return
	}
	h.events = append(h.events, ev)
}

// Entries returns a shallow copy of all recorded events in insertion order.
func (h *History) Entries() []Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]Event, len(h.events))
	copy(out, h.events)
	return out
}

// EntriesFor returns events recorded for the given service name.
func (h *History) EntriesFor(service string) []Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	var out []Event
	for _, e := range h.events {
		if e.Service == service {
			out = append(out, e)
		}
	}
	return out
}
