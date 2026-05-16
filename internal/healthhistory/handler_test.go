package healthhistory

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func newHistory() *History { return New(50) }

func TestHandler_EmptyHistory(t *testing.T) {
	h := newHistory()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	Handler(h).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestHandler_AllEvents(t *testing.T) {
	h := newHistory()
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("svc-b", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	Handler(h).ServeHTTP(rec, req)
	var events []Event
	_ = json.NewDecoder(rec.Body).Decode(&events)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestHandler_FilterByService(t *testing.T) {
	h := newHistory()
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_SERVING)
	h.Record("svc-b", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	h.Record("svc-a", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history?service=svc-a", nil)
	Handler(h).ServeHTTP(rec, req)
	var events []Event
	_ = json.NewDecoder(rec.Body).Decode(&events)
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
	h := newHistory()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/history", nil)
	Handler(h).ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandler_ContentTypeJSON(t *testing.T) {
	h := newHistory()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	Handler(h).ServeHTTP(rec, req)
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}

func TestHandler_TimestampsAreUTC(t *testing.T) {
	h := newHistory()
	h.Record("svc", grpc_health_v1.HealthCheckResponse_SERVING)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	Handler(h).ServeHTTP(rec, req)
	var events []Event
	_ = json.NewDecoder(rec.Body).Decode(&events)
	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	if events[0].Timestamp.Location().String() != "UTC" {
		t.Errorf("expected UTC timestamp, got %v", events[0].Timestamp.Location())
	}
}
