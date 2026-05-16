package uptime_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nicholasgasior/grpc-healthd/internal/uptime"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	start := time.Now().UTC().Add(-30 * time.Second)
	tr := uptime.NewWithTime(start)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uptime", nil)
	uptime.Handler(tr)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var info uptime.Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if info.UptimeSeconds < 29 {
		t.Errorf("expected uptime_seconds >= 29, got %d", info.UptimeSeconds)
	}
	if info.StartTime == "" {
		t.Error("start_time should not be empty")
	}
}

func TestHandler_ContentType(t *testing.T) {
	tr := uptime.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uptime", nil)
	uptime.Handler(tr)(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	tr := uptime.New()
	for _, method := range []string{http.MethodPost, http.MethodDelete, http.MethodPut} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/uptime", nil)
		uptime.Handler(tr)(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", method, rec.Code)
		}
	}
}

func TestHandler_StartTimeFormat(t *testing.T) {
	start := time.Date(2024, 3, 10, 8, 30, 0, 0, time.UTC)
	tr := uptime.NewWithTime(start)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uptime", nil)
	uptime.Handler(tr)(rec, req)

	var info uptime.Info
	json.NewDecoder(rec.Body).Decode(&info) //nolint:errcheck
	if info.StartTime != "2024-03-10T08:30:00Z" {
		t.Errorf("unexpected start_time format: %s", info.StartTime)
	}
}
