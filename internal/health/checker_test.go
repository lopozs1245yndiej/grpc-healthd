package health_test

import (
	"context"
	"testing"

	"github.com/your-org/grpc-healthd/internal/health"
)

func TestNewChecker(t *testing.T) {
	c := health.NewChecker()
	if c == nil {
		t.Fatal("expected non-nil Checker")
	}
}

func TestGetStatus_Unknown(t *testing.T) {
	c := health.NewChecker()
	status := c.GetStatus(context.Background(), "unknown-service")
	if status != health.StatusUnknown {
		t.Errorf("expected UNKNOWN, got %s", status)
	}
}

func TestSetAndGetStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   health.Status
		expected string
	}{
		{"serving", health.StatusServing, "SERVING"},
		{"not-serving", health.StatusNotServing, "NOT_SERVING"},
		{"unknown", health.StatusUnknown, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := health.NewChecker()
			c.SetStatus("my-service", tt.status)
			got := c.GetStatus(context.Background(), "my-service")
			if got != tt.status {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
			if got.String() != tt.expected {
				t.Errorf("expected string %q, got %q", tt.expected, got.String())
			}
		})
	}
}

func TestSetStatus_Overwrite(t *testing.T) {
	c := health.NewChecker()
	c.SetStatus("svc", health.StatusServing)
	c.SetStatus("svc", health.StatusNotServing)
	got := c.GetStatus(context.Background(), "svc")
	if got != health.StatusNotServing {
		t.Errorf("expected NOT_SERVING after overwrite, got %s", got)
	}
}

func TestListServices(t *testing.T) {
	c := health.NewChecker()
	c.SetStatus("alpha", health.StatusServing)
	c.SetStatus("beta", health.StatusNotServing)

	services := c.ListServices()
	if len(services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(services))
	}
	if services["alpha"].Status != health.StatusServing {
		t.Errorf("expected alpha SERVING")
	}
	if services["beta"].Status != health.StatusNotServing {
		t.Errorf("expected beta NOT_SERVING")
	}
}

func TestListServices_Empty(t *testing.T) {
	c := health.NewChecker()
	services := c.ListServices()
	if len(services) != 0 {
		t.Errorf("expected empty map, got %d entries", len(services))
	}
}
