package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns an http.Handler that serves Prometheus metrics from reg.
// If reg is nil, the default prometheus registry is used.
func Handler(reg *prometheus.Registry) http.Handler {
	if reg == nil {
		return promhttp.Handler()
	}
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		EnableOpenMetrics: false,
	})
}

// NewServer creates a minimal *http.ServeMux that exposes /metrics.
func NewServer(addr string, reg *prometheus.Registry) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", Handler(reg))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
