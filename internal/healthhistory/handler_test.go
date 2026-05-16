package healthhistory

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func newHistory(maxSize int) *History { return New(maxSize) }

func TestHandler_EmptyHistory(t *testing.T) {
	h := newHistory(10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	Handler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected empty slice, got %d events", len(events))
	}
}

func TestHandler_AllEvents(t *testing.T) {
	h := newHistory(10)
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("svc-b", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	Handler(h).ServeHTTP(rec, req)

	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestHandler_FilterByService(t *testing.T) {
	h := newHistory(10)
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("svc-b", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history?service=svc-a", nil)
	Handler(h).ServeHTTP(rec, req)

	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events for svc-a, got %d", len(events))
	}
	for _, e := range events {
		if e.Service != "svc-a" {
			t.Errorf("unexpected service %q in filtered results", e.Service)
		}
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := newHistory(10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/history", nil)
	Handler(h).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_ContentTypeJSON(t *testing.T) {
	h := newHistory(10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	Handler(h).ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}
