// Package watcher provides periodic health status updates for registered services.
package watcher

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/example/grpc-healthd/internal/health"
)

// HTTPChecker performs an HTTP GET and reports serving if the response is 2xx.
type HTTPChecker struct {
	URL string
}

// Check performs the HTTP health probe.
func (h *HTTPChecker) Check(ctx context.Context) grpc_health_v1.HealthCheckResponse_ServingStatus {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, nil)
	if err != nil {
		return grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return grpc_health_v1.HealthCheckResponse_SERVING
	}
	return grpc_health_v1.HealthCheckResponse_NOT_SERVING
}

// Probe is implemented by any type that can report a serving status.
type Probe interface {
	Check(ctx context.Context) grpc_health_v1.HealthCheckResponse_ServingStatus
}

// ServiceWatcher polls a Probe at a fixed interval and updates a health.Checker.
type ServiceWatcher struct {
	service  string
	probe    Probe
	interval time.Duration
	checker  *health.Checker
}

// NewServiceWatcher creates a watcher for the given service.
func NewServiceWatcher(service string, probe Probe, interval time.Duration, checker *health.Checker) *ServiceWatcher {
	return &ServiceWatcher{
		service:  service,
		probe:    probe,
		interval: interval,
		checker:  checker,
	}
}

// Run starts polling until ctx is cancelled.
func (sw *ServiceWatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(sw.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			status := sw.probe.Check(ctx)
			sw.checker.SetStatus(sw.service, status)
			log.Printf("watcher: service=%q status=%s", sw.service, status)
		}
	}
}

// Manager owns multiple ServiceWatchers and starts/stops them together.
type Manager struct {
	watchers []*ServiceWatcher
	wg       sync.WaitGroup
}

// Add registers a watcher with the manager.
func (m *Manager) Add(sw *ServiceWatcher) {
	m.watchers = append(m.watchers, sw)
}

// Start launches all registered watchers.
func (m *Manager) Start(ctx context.Context) {
	for _, sw := range m.watchers {
		m.wg.Add(1)
		go func(w *ServiceWatcher) {
			defer m.wg.Done()
			w.Run(ctx)
		}(sw)
	}
}

// Wait blocks until all watchers have stopped.
func (m *Manager) Wait() { m.wg.Wait() }
