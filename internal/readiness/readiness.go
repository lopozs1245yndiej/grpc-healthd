// Package readiness tracks whether the daemon itself is ready to serve traffic.
// Readiness is determined by whether at least one watched service has a known
// health status (i.e. the first watcher poll has completed successfully).
package readiness

import (
	"net/http"
	"sync"
)

// Tracker records daemon readiness state.
type Tracker struct {
	mu    sync.RWMutex
	ready bool
}

// New returns a new Tracker in the not-ready state.
func New() *Tracker {
	return &Tracker{}
}

// SetReady marks the daemon as ready.
func (t *Tracker) SetReady() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ready = true
}

// SetNotReady marks the daemon as not ready.
func (t *Tracker) SetNotReady() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ready = false
}

// IsReady reports whether the daemon is currently ready.
func (t *Tracker) IsReady() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ready
}

// Handler returns an HTTP handler that responds 200 OK when ready and
// 503 Service Unavailable when not ready. Suitable for use as a Kubernetes
// readinessProbe HTTP endpoint.
func Handler(t *Tracker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if t.IsReady() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("not ready"))
	})
}
