package ratelimit

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Capacity  int     `json:"capacity"`
	RatePerSec float64 `json:"rate_per_sec"`
	Enabled   bool    `json:"enabled"`
}

// Handler returns an http.Handler that exposes the current rate-limit
// configuration as a JSON endpoint. It is intended to be mounted on the
// admin/metrics server so operators can inspect live settings.
func Handler(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := statusResponse{
			Capacity:   cfg.Capacity,
			RatePerSec: cfg.RatePerSec,
			Enabled:    cfg.Enabled,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
}
