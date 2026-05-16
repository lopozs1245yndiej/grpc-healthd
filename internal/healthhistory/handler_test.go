package healthhistory_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/yourorg/grpc-healthd/internal/healthhistory"
)

func newHistory(t *testing.T) *healthhistory.History {
	t.Helper()
	return healthhistory.New(50)
}

func TestHandler_EmptyHistory(t *testing.T) {
	h := newHistory(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)

	healthhistory.Handler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var events []healthhistory.Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected empty slice, got %d events", len(events))
	}
}

func TestHandler_AllEvents(t *testing.T) {
	h := newHistory(t)
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("svc-b", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	healthhistory.Handler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var events []healthhistory.Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestHandler_FilterByService(t *testing.T) {
	h := newHistory(t)
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("svc-b", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history?service=svc-a", nil)
	healthhistory.Handler(h).ServeHTTP(rec, req)

	var events []healthhistory.Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events for svc-a, got %d", len(events))
	}
	for _, e := range events {
		if e.Service != "svc-a" {
			t.Errorf("unexpected service %q in filtered results", e.Service)
		}
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := newHistory(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/history", nil)
	healthhistory.Handler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_ContentTypeJSON(t *testing.T) {
	h := newHistory(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	healthhistory.Handler(h).ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}
