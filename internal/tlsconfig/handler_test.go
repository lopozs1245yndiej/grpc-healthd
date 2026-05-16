package tlsconfig_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/grpc-healthd/internal/tlsconfig"
)

func TestHandler_Disabled(t *testing.T) {
	cfg := tlsconfig.Config{Enabled: false}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tls", nil)
	tlsconfig.Handler(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", body["enabled"])
	}
	if _, ok := body["cert_file"]; ok {
		t.Error("cert_file should be omitted when disabled")
	}
}

func TestHandler_Enabled(t *testing.T) {
	cfg := tlsconfig.Config{Enabled: true, CertFile: "/certs/tls.crt", KeyFile: "/certs/tls.key"}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tls", nil)
	tlsconfig.Handler(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["enabled"] != true {
		t.Errorf("expected enabled=true, got %v", body["enabled"])
	}
	if body["cert_file"] != "/certs/tls.crt" {
		t.Errorf("unexpected cert_file: %v", body["cert_file"])
	}
	if body["key_file"] != "/certs/tls.key" {
		t.Errorf("unexpected key_file: %v", body["key_file"])
	}
}

func TestHandler_ContentType(t *testing.T) {
	cfg := tlsconfig.Config{Enabled: false}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tls", nil)
	tlsconfig.Handler(cfg).ServeHTTP(rec, req)
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	cfg := tlsconfig.Config{Enabled: false}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tls", nil)
	tlsconfig.Handler(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
