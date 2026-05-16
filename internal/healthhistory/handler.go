package healthhistory

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that exposes the health event history as
// JSON. An optional "service" query parameter filters results by service name.
func Handler(h *History) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		service := r.URL.Query().Get("service")
		events := h.Events(service)
		if events == nil {
			events = []Event{}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	})
}
