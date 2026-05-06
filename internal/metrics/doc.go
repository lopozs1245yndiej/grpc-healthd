// Package metrics exposes Prometheus instrumentation for grpc-healthd.
//
// It tracks the total number of health checks performed, their duration,
// and the current serving status of each watched service.
//
// Usage:
//
//	reg := prometheus.NewRegistry()
//	m := metrics.New(reg)
//	m.RecordServing("my-service", 0.012)
package metrics
