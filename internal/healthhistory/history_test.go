package healthhistory

import (
	"testing"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func TestNew_DefaultMaxSize(t *testing.T) {
	h := New(0)
	if h.maxSize != 200 {
		t.Fatalf("expected default maxSize 200, got %d", h.maxSize)
	}
}

func TestRecord_SingleEntry(t *testing.T) {
	h := New(10)
	h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)
	events := h.Events("")
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Service != "svc" {
		t.Errorf("unexpected service %q", events[0].Service)
	}
	if events[0].Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("unexpected status %v", events[0].Status)
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	h := New(10)
	h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)
	events := h.Events("")
	if events[0].Timestamp.Location().String() != "UTC" {
		t.Errorf("timestamp not UTC: %v", events[0].Timestamp.Location())
	}
}

func TestRecord_Eviction(t *testing.T) {
	h := New(3)
	h.Record("a", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("b", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("c", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("d", grpc_health_v1.HealthCheckResponse_SERVING)
	events := h.Events("")
	if len(events) != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", len(events))
	}
	if events[0].Service != "b" {
		t.Errorf("expected oldest to be 'b', got %q", events[0].Service)
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	h := New(10)
	h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)
	events := h.Events("")
	events[0].Service = "mutated"
	original := h.Events("")
	if original[0].Service == "mutated" {
		t.Error("Events should return a copy, not a reference to internal slice")
	}
}

func TestEvents_FilterByService(t *testing.T) {
	h := New(20)
	h.Record("alpha", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("beta", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	h.Record("alpha", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	results := h.Events("alpha")
	if len(results) != 2 {
		t.Fatalf("expected 2 events for alpha, got %d", len(results))
	}
}
