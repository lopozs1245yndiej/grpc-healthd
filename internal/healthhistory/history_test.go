package healthhistory

import (
	"testing"
	"time"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func TestNew_DefaultMaxSize(t *testing.T) {
	h := New()
	if h.maxSize != defaultMaxSize {
		t.Fatalf("expected maxSize %d, got %d", defaultMaxSize, h.maxSize)
	}
}

func TestRecord_SingleEntry(t *testing.T) {
	h := New()
	h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)

	events := h.Events("")
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Service != "svc" {
		t.Errorf("expected service 'svc', got %q", events[0].Service)
	}
	if events[0].Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("unexpected status: %v", events[0].Status)
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	h := New()
	before := time.Now().UTC()
	h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)
	after := time.Now().UTC()

	events := h.Events("")
	ts := events[0].Timestamp
	if ts.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", ts.Location())
	}
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
}

func TestRecord_Eviction(t *testing.T) {
	h := NewWithSize(3)
	for i := 0; i < 5; i++ {
		h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)
	}
	events := h.Events("")
	if len(events) != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", len(events))
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	h := New()
	h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)

	a := h.Events("")
	a[0].Service = "mutated"

	b := h.Events("")
	if b[0].Service == "mutated" {
		t.Error("Events should return a copy, not expose internal state")
	}
}

func TestEvents_FilterByService(t *testing.T) {
	h := New()
	h.Record("alpha", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("beta", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	h.Record("alpha", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	alpha := h.Events("alpha")
	if len(alpha) != 2 {
		t.Fatalf("expected 2 alpha events, got %d", len(alpha))
	}
	for _, e := range alpha {
		if e.Service != "alpha" {
			t.Errorf("filter returned wrong service: %q", e.Service)
		}
	}
}

func TestNewWithSize_ZeroUsesDefault(t *testing.T) {
	h := NewWithSize(0)
	if h.maxSize != defaultMaxSize {
		t.Errorf("expected default maxSize %d, got %d", defaultMaxSize, h.maxSize)
	}
}
