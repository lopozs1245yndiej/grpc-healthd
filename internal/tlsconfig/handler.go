package tlsconfig

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"cert_file,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
}

// Handler returns an http.Handler that reports the current TLS configuration.
// It responds only to GET requests and returns JSON.
//
// Example response (TLS enabled):
//
//	{"enabled":true,"cert_file":"/etc/certs/tls.crt","key_file":"/etc/certs/tls.key"}
//
// Example response (TLS disabled):
//
//	{"enabled":false}
func Handler(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		resp := statusResponse{Enabled: cfg.Enabled}
		if cfg.Enabled {
			resp.CertFile = cfg.CertFile
			resp.KeyFile = cfg.KeyFile
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
}
