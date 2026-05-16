package buildinfo

import (
	"encoding/json"
	"net/http"
)

// Handler returns an HTTP handler that serves build information as JSON.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		info := Get()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})
}
