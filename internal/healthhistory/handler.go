package healthhistory

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that exposes recorded history as JSON.
//
// GET /history          — all events
// GET /history?service= — events for a specific service
func Handler(h *History) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var events []Event
		if svc := r.URL.Query().Get("service"); svc != "" {
			events = h.EntriesFor(svc)
		} else {
			events = h.Entries()
		}

		if events == nil {
			events = []Event{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})
}
