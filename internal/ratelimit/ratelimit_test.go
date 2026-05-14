package ratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/grpc-healthd/internal/ratelimit"
)

func TestAllow_WithinLimit(t *testing.T) {
	l := ratelimit.New(10, time.Second, 5)
	for i := 0; i < 5; i++ {
		if !l.Allow("127.0.0.1") {
			t.Fatalf("expected Allow=true on request %d", i+1)
		}
	}
}

func TestAllow_ExceedsCapacity(t *testing.T) {
	l := ratelimit.New(10, time.Second, 3)
	for i := 0; i < 3; i++ {
		l.Allow("10.0.0.1")
	}
	if l.Allow("10.0.0.1") {
		t.Fatal("expected Allow=false after capacity exhausted")
	}
}

func TestAllow_DifferentKeys(t *testing.T) {
	l := ratelimit.New(10, time.Second, 1)
	if !l.Allow("192.168.1.1") {
		t.Fatal("expected first request from 192.168.1.1 to be allowed")
	}
	if !l.Allow("192.168.1.2") {
		t.Fatal("expected first request from 192.168.1.2 to be allowed")
	}
	if l.Allow("192.168.1.1") {
		t.Fatal("expected second request from 192.168.1.1 to be denied")
	}
}

func TestAllow_RefillAfterInterval(t *testing.T) {
	l := ratelimit.New(2, 50*time.Millisecond, 2)
	l.Allow("host")
	l.Allow("host")
	if l.Allow("host") {
		t.Fatal("expected deny after capacity exhausted")
	}
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("host") {
		t.Fatal("expected allow after refill interval")
	}
}

func TestMiddleware_AllowsRequest(t *testing.T) {
	l := ratelimit.New(10, time.Second, 5)
	handler := l.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestMiddleware_BlocksExcessRequests(t *testing.T) {
	l := ratelimit.New(10, time.Second, 2)
	handler := l.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.2:9999"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.2:9999"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rr.Code)
	}
}
