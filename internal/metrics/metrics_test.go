package metrics_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

func newRegistry() (*metrics.Metrics, *prometheus.Registry) {
	reg := prometheus.NewRegistry()
	m := metrics.New(reg)
	return m, reg
}

func gatherGauge(t *testing.T, reg *prometheus.Registry, name, service string) float64 {
	t.Helper()
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather error: %v", err)
	}
	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		for _, m := range mf.GetMetric() {
			for _, lp := range m.GetLabel() {
				if lp.GetName() == "service" && lp.GetValue() == service {
					return m.GetGauge().GetValue()
				}
			}
		}
	}
	t.Fatalf("metric %q for service %q not found", name, service)
	return 0
}

func gatherCounter(t *testing.T, reg *prometheus.Registry, name, service, status string) float64 {
	t.Helper()
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather error: %v", err)
	}
	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		for _, m := range mf.GetMetric() {
			labels := map[string]string{}
			for _, lp := range m.GetLabel() {
				labels[lp.GetName()] = lp.GetValue()
			}
			if labels["service"] == service && labels["status"] == status {
				return m.GetCounter().GetValue()
			}
		}
	}
	return 0
}

func TestRecordServing(t *testing.T) {
	m, reg := newRegistry()
	m.RecordServing("svc-a", 0.05)

	if v := gatherGauge(t, reg, "grpc_healthd_service_status", "svc-a"); v != 1 {
		t.Errorf("expected gauge 1, got %v", v)
	}
	if v := gatherCounter(t, reg, "grpc_healthd_checks_total", "svc-a", "SERVING"); v != 1 {
		t.Errorf("expected counter 1, got %v", v)
	}
}

func TestRecordNotServing(t *testing.T) {
	m, reg := newRegistry()
	m.RecordNotServing("svc-b", 0.1)

	if v := gatherGauge(t, reg, "grpc_healthd_service_status", "svc-b"); v != 0 {
		t.Errorf("expected gauge 0, got %v", v)
	}
	if v := gatherCounter(t, reg, "grpc_healthd_checks_total", "svc-b", "NOT_SERVING"); v != 1 {
		t.Errorf("expected counter 1, got %v", v)
	}
}

func TestRecordMultiple(t *testing.T) {
	m, reg := newRegistry()
	m.RecordServing("svc-c", 0.01)
	m.RecordServing("svc-c", 0.02)
	m.RecordNotServing("svc-c", 0.03)

	if v := gatherCounter(t, reg, "grpc_healthd_checks_total", "svc-c", "SERVING"); v != 2 {
		t.Errorf("expected SERVING counter 2, got %v", v)
	}
	if v := gatherCounter(t, reg, "grpc_healthd_checks_total", "svc-c", "NOT_SERVING"); v != 1 {
		t.Errorf("expected NOT_SERVING counter 1, got %v", v)
	}
	// Last call was NOT_SERVING so gauge should be 0.
	if v := gatherGauge(t, reg, "grpc_healthd_service_status", "svc-c"); v != 0 {
		t.Errorf("expected gauge 0, got %v", v)
	}
}

func TestHistogramPresent(t *testing.T) {
	m, reg := newRegistry()
	m.RecordServing("svc-d", 0.007)

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	var found *dto.MetricFamily
	for _, mf := range mfs {
		if mf.GetName() == "grpc_healthd_check_duration_seconds" {
			found = mf
			break
		}
	}
	if found == nil {
		t.Fatal("histogram metric not found")
	}
	if len(found.GetMetric()) == 0 {
		t.Fatal("no histogram samples recorded")
	}
}
