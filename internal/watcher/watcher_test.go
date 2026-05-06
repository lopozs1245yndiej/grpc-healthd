package watcher_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/example/grpc-healthd/internal/health"
	"github.com/example/grpc-healthd/internal/watcher"
)

func TestHTTPChecker_Serving(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := &watcher.HTTPChecker{URL: ts.URL}
	status := c.Check(context.Background())
	if status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("expected SERVING, got %s", status)
	}
}

func TestHTTPChecker_NotServing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := &watcher.HTTPChecker{URL: ts.URL}
	status := c.Check(context.Background())
	if status != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("expected NOT_SERVING, got %s", status)
	}
}

func TestHTTPChecker_Unreachable(t *testing.T) {
	c := &watcher.HTTPChecker{URL: "http://127.0.0.1:1"}
	status := c.Check(context.Background())
	if status != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("expected NOT_SERVING, got %s", status)
	}
}

type stubProbe struct{ status grpc_health_v1.HealthCheckResponse_ServingStatus }

func (s *stubProbe) Check(_ context.Context) grpc_health_v1.HealthCheckResponse_ServingStatus {
	return s.status
}

func TestServiceWatcher_UpdatesChecker(t *testing.T) {
	checker := health.NewChecker()
	probe := &stubProbe{status: grpc_health_v1.HealthCheckResponse_SERVING}

	sw := watcher.NewServiceWatcher("svc", probe, 20*time.Millisecond, checker)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go sw.Run(ctx)
	<-ctx.Done()

	status := checker.GetStatus("svc")
	if status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("expected SERVING after watcher run, got %s", status)
	}
}

func TestManager_StartsAllWatchers(t *testing.T) {
	checker := health.NewChecker()
	probe := &stubProbe{status: grpc_health_v1.HealthCheckResponse_SERVING}

	var mgr watcher.Manager
	mgr.Add(watcher.NewServiceWatcher("a", probe, 20*time.Millisecond, checker))
	mgr.Add(watcher.NewServiceWatcher("b", probe, 20*time.Millisecond, checker))

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	mgr.Start(ctx)
	mgr.Wait()

	for _, svc := range []string{"a", "b"} {
		if s := checker.GetStatus(svc); s != grpc_health_v1.HealthCheckResponse_SERVING {
			t.Errorf("service %q: expected SERVING, got %s", svc, s)
		}
	}
}
