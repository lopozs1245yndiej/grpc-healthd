package ratelimit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsConfig(t *testing.T) {
	cfg := Config{
		Capacity:   50,
		RatePerSec: 10.0,
		Enabled:    true,
	}

	h := Handler(cfg)
	req := httptest.NewRequest(http.MethodGet, "/ratelimit", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp statusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Capacity != cfg.Capacity {
		t.Errorf("capacity: want %d, got %d", cfg.Capacity, resp.Capacity)
	}
	if resp.RatePerSec != cfg.RatePerSec {
		t.Errorf("rate_per_sec: want %f, got %f", cfg.RatePerSec, resp.RatePerSec)
	}
	if resp.Enabled != cfg.Enabled {
		t.Errorf("enabled: want %v, got %v", cfg.Enabled, resp.Enabled)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := Handler(Config{})
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/ratelimit", nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", method, rr.Code)
		}
	}
}

func TestHandler_ContentTypeJSON(t *testing.T) {
	h := Handler(Config{Capacity: 10, RatePerSec: 1.0, Enabled: false})
	req := httptest.NewRequest(http.MethodGet, "/ratelimit", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestHandler_DisabledConfig(t *testing.T) {
	cfg := Config{
		Capacity:   0,
		RatePerSec: 0,
		Enabled:    false,
	}

	h := Handler(cfg)
	req := httptest.NewRequest(http.MethodGet, "/ratelimit", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	var resp statusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Enabled {
		t.Error("expected enabled=false")
	}
}
