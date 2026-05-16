package readiness_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jamesread/grpc-healthd/internal/readiness"
)

func TestNew_NotReadyByDefault(t *testing.T) {
	tr := readiness.New()
	if tr.IsReady() {
		t.Fatal("expected tracker to be not ready by default")
	}
}

func TestSetReady(t *testing.T) {
	tr := readiness.New()
	tr.SetReady()
	if !tr.IsReady() {
		t.Fatal("expected tracker to be ready after SetReady")
	}
}

func TestSetNotReady(t *testing.T) {
	tr := readiness.New()
	tr.SetReady()
	tr.SetNotReady()
	if tr.IsReady() {
		t.Fatal("expected tracker to be not ready after SetNotReady")
	}
}

func TestHandler_NotReady(t *testing.T) {
	tr := readiness.New()
	h := readiness.Handler(tr)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
	if rec.Body.String() != "not ready" {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestHandler_Ready(t *testing.T) {
	tr := readiness.New()
	tr.SetReady()
	h := readiness.Handler(tr)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	tr := readiness.New()
	h := readiness.Handler(tr)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/readyz", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
