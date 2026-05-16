package healthhistory

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that serves health-transition events as
// JSON. An optional "service" query parameter filters the results.
func Handler(h *History) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		service := r.URL.Query().Get("service")
		events := h.Events(service)

		// Return an empty JSON array rather than null.
		if events == nil {
			events = []Event{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	})
}
