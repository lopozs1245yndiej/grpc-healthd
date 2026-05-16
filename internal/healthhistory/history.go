package healthhistory

import (
	"sync"
	"time"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// Event records a single health-status transition.
type Event struct {
	Service   string                          `json:"service"`
	Status    grpc_health_v1.HealthCheckResponse_ServingStatus `json:"status"`
	StatusStr string                          `json:"status_text"`
	Timestamp time.Time                       `json:"timestamp"`
}

// History is a bounded, thread-safe ring buffer of health events.
type History struct {
	mu      sync.Mutex
	events  []Event
	maxSize int
}

// New returns a History that retains at most maxSize events.
// If maxSize is <= 0 it defaults to 200.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 200
	}
	return &History{maxSize: maxSize}
}

// Record appends a new event, evicting the oldest when the buffer is full.
func (h *History) Record(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()

	e := Event{
		Service:   service,
		Status:    status,
		StatusStr: status.String(),
		Timestamp: time.Now().UTC(),
	}

	if len(h.events) >= h.maxSize {
		h.events = h.events[1:]
	}
	h.events = append(h.events, e)
}

// Events returns a copy of all recorded events, optionally filtered by
// service name (empty string means "all services").
func (h *History) Events(service string) []Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]Event, 0, len(h.events))
	for _, e := range h.events {
		if service == "" || e.Service == service {
			out = append(out, e)
		}
	}
	return out
}
