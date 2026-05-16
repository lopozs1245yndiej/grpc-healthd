package audit

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that serves the audit log as JSON.
// GET /audit returns all entries; other methods return 405.
func Handler(l *Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		entries := l.Entries()
		if entries == nil {
			entries = []Entry{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})
}
