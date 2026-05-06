// Package metrics provides Prometheus metrics collection for grpc-healthd.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the health daemon.
type Metrics struct {
	CheckTotal    *prometheus.CounterVec
	CheckDuration *prometheus.HistogramVec
	ServiceStatus *prometheus.GaugeVec
}

// New creates and registers a new set of Prometheus metrics.
func New(reg prometheus.Registerer) *Metrics {
	factory := promauto.With(reg)

	return &Metrics{
		CheckTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_healthd_checks_total",
				Help: "Total number of health checks performed.",
			},
			[]string{"service", "status"},
		),
		CheckDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_healthd_check_duration_seconds",
				Help:    "Duration of health check probes in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service"},
		),
		ServiceStatus: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "grpc_healthd_service_status",
				Help: "Current status of a watched service (1=SERVING, 0=NOT_SERVING).",
			},
			[]string{"service"},
		),
	}
}

// RecordServing records a SERVING result for the given service.
func (m *Metrics) RecordServing(service string, durationSec float64) {
	m.CheckTotal.WithLabelValues(service, "SERVING").Inc()
	m.CheckDuration.WithLabelValues(service).Observe(durationSec)
	m.ServiceStatus.WithLabelValues(service).Set(1)
}

// RecordNotServing records a NOT_SERVING result for the given service.
func (m *Metrics) RecordNotServing(service string, durationSec float64) {
	m.CheckTotal.WithLabelValues(service, "NOT_SERVING").Inc()
	m.CheckDuration.WithLabelValues(service).Observe(durationSec)
	m.ServiceStatus.WithLabelValues(service).Set(0)
}
