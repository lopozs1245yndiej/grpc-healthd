package uptime

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.HandlerFunc that serves uptime information as JSON.
func Handler(tr *Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := tr.Snapshot()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
