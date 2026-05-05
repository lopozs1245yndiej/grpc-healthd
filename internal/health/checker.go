package health

import (
	"context"
	"sync"
	"time"
)

// Status represents the health status of a service.
type Status int

const (
	StatusUnknown    Status = iota
	StatusServing
	StatusNotServing
)

func (s Status) String() string {
	switch s {
	case StatusServing:
		return "SERVING"
	case StatusNotServing:
		return "NOT_SERVING"
	default:
		return "UNKNOWN"
	}
}

// ServiceStatus holds the current status and last checked time for a service.
type ServiceStatus struct {
	Status    Status
	CheckedAt time.Time
}

// Checker manages health statuses for named services.
type Checker struct {
	mu       sync.RWMutex
	services map[string]*ServiceStatus
}

// NewChecker creates a new Checker instance.
func NewChecker() *Checker {
	return &Checker{
		services: make(map[string]*ServiceStatus),
	}
}

// SetStatus updates the health status for the given service name.
func (c *Checker) SetStatus(service string, status Status) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[service] = &ServiceStatus{
		Status:    status,
		CheckedAt: time.Now(),
	}
}

// GetStatus returns the current health status for the given service.
// Returns StatusUnknown if the service has not been registered.
func (c *Checker) GetStatus(_ context.Context, service string) Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ss, ok := c.services[service]
	if !ok {
		return StatusUnknown
	}
	return ss.Status
}

// ListServices returns a snapshot of all registered service statuses.
func (c *Checker) ListServices() map[string]ServiceStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]ServiceStatus, len(c.services))
	for k, v := range c.services {
		result[k] = *v
	}
	return result
}
