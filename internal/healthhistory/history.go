package healthhistory

import (
	"sync"
	"time"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

const defaultMaxSize = 100

// Event represents a single health status transition recorded for a service.
type Event struct {
	Service   string                          `json:"service"`
	Status    grpc_health_v1.HealthCheckResponse_ServingStatus `json:"status"`
	Timestamp time.Time                       `json:"timestamp"`
}

// History maintains a bounded, ordered log of health status change events.
type History struct {
	mu      sync.RWMutex
	events  []Event
	maxSize int
}

// New creates a History with the default maximum size.
func New() *History {
	return &History{maxSize: defaultMaxSize}
}

// NewWithSize creates a History with a custom maximum number of events.
func NewWithSize(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = defaultMaxSize
	}
	return &History{maxSize: maxSize}
}

// Record appends a new event to the history, evicting the oldest entry when
// the buffer is full.
func (h *History) Record(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()

	event := Event{
		Service:   service,
		Status:    status,
		Timestamp: time.Now().UTC(),
	}

	if len(h.events) >= h.maxSize {
		h.events = h.events[1:]
	}
	h.events = append(h.events, event)
}

// Events returns a shallow copy of all recorded events, optionally filtered
// by service name when filter is non-empty.
func (h *History) Events(filter string) []Event {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]Event, 0, len(h.events))
	for _, e := range h.events {
		if filter == "" || e.Service == filter {
			result = append(result, e)
		}
	}
	return result
}
