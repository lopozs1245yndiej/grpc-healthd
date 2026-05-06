package metrics_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

func TestHandler_MetricsEndpoint(t *testing.T) {
	reg := prometheus.NewRegistry()
	m := metrics.New(reg)
	m.RecordServing("svc-handler", 0.005)

	ts := httptest.NewServer(metrics.Handler(reg))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "grpc_healthd_checks_total") {
		t.Error("response does not contain grpc_healthd_checks_total")
	}
	if !strings.Contains(string(body), "grpc_healthd_service_status") {
		t.Error("response does not contain grpc_healthd_service_status")
	}
}

func TestNewServer_Healthz(t *testing.T) {
	reg := prometheus.NewRegistry()
	_ = metrics.New(reg)

	srv := metrics.NewServer("127.0.0.1:0", reg)
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Errorf("expected body 'ok', got %q", string(body))
	}
}

func TestHandler_NilRegistry(t *testing.T) {
	h := metrics.Handler(nil)
	if h == nil {
		t.Fatal("expected non-nil handler for nil registry")
	}
}
