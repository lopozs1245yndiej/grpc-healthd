package healthhistory

import (
	"sync"
	"time"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// Event records a single health status transition for a service.
type Event struct {
	Service   string                          `json:"service"`
	Status    grpc_health_v1.HealthCheckResponse_ServingStatus `json:"status"`
	Timestamp time.Time                       `json:"timestamp"`
}

// History maintains a bounded ring-buffer of health status events.
type History struct {
	mu      sync.RWMutex
	events  []Event
	maxSize int
}

// New creates a History with the given maximum number of retained events.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &History{maxSize: maxSize}
}

// Record appends a new event to the history, evicting the oldest entry when
// the buffer is at capacity.
func (h *History) Record(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := Event{
		Service:   service,
		Status:    status,
		Timestamp: time.Now().UTC(),
	}
	if len(h.events) >= h.maxSize {
		h.events = h.events[1:]
	}
	h.events = append(h.events, e)
}

// Events returns a copy of all recorded events, optionally filtered by service
// name (pass an empty string to retrieve all events).
func (h *History) Events(service string) []Event {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Event, 0, len(h.events))
	for _, e := range h.events {
		if service == "" || e.Service == service {
			out = append(out, e)
		}
	}
	return out
}
