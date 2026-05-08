package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grpc-healthd/internal/admin"
	"github.com/grpc-healthd/internal/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func newHandler(t *testing.T) (*admin.Handler, *health.Checker) {
	t.Helper()
	c := health.NewChecker()
	return admin.New(c), c
}

func TestListServices_Empty(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ListServices(rec, httptest.NewRequest(http.MethodGet, "/admin/services", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out []map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty list, got %v", out)
	}
}

func TestListServices_WithEntries(t *testing.T) {
	h, c := newHandler(t)
	c.SetStatus("svc-a", grpc_health_v1.HealthCheckResponse_SERVING)
	c.SetStatus("svc-b", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	rec := httptest.NewRecorder()
	h.ListServices(rec, httptest.NewRequest(http.MethodGet, "/admin/services", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out []map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&out)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestListServices_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ListServices(rec, httptest.NewRequest(http.MethodPost, "/admin/services", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestSetStatus_OK(t *testing.T) {
	h, c := newHandler(t)
	body := `{"service":"svc-x","status":"SERVING"}`
	rec := httptest.NewRecorder()
	h.SetStatus(rec, httptest.NewRequest(http.MethodPost, "/admin/set", strings.NewReader(body)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	st, _ := c.GetStatus("svc-x")
	if st != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("expected SERVING, got %v", st)
	}
}

func TestSetStatus_InvalidJSON(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.SetStatus(rec, httptest.NewRequest(http.MethodPost, "/admin/set", bytes.NewBufferString("not-json")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestSetStatus_UnknownStatus(t *testing.T) {
	h, _ := newHandler(t)
	body := `{"service":"svc-y","status":"MAYBE"}`
	rec := httptest.NewRecorder()
	h.SetStatus(rec, httptest.NewRequest(http.MethodPost, "/admin/set", strings.NewReader(body)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestNewServer_Routes(t *testing.T) {
	h, _ := newHandler(t)
	srv := admin.NewServer(":0", h)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}
