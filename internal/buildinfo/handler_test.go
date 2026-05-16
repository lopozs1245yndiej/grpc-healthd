package buildinfo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicholasgasior/grpc-healthd/internal/buildinfo"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	buildinfo.Version = "1.2.3"
	buildinfo.Commit = "abc1234"
	buildinfo.BuildDate = "2024-01-01T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	buildinfo.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var info buildinfo.Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if info.Version != "1.2.3" {
		t.Errorf("expected version 1.2.3, got %s", info.Version)
	}
	if info.Commit != "abc1234" {
		t.Errorf("expected commit abc1234, got %s", info.Commit)
	}
	if info.BuildDate != "2024-01-01T00:00:00Z" {
		t.Errorf("expected build date 2024-01-01T00:00:00Z, got %s", info.BuildDate)
	}
	if info.GoVersion == "" {
		t.Error("expected non-empty go version")
	}
}

func TestHandler_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	buildinfo.Handler().ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/version", nil)
		rec := httptest.NewRecorder()

		buildinfo.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", method, rec.Code)
		}
	}
}

func TestGet_RetrievedIsUTC(t *testing.T) {
	info := buildinfo.Get()
	if info.Retrieved.Location().String() != "UTC" {
		t.Errorf("expected UTC, got %s", info.Retrieved.Location())
	}
}
