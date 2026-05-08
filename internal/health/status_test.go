package health_test

import (
	"testing"

	"github.com/grpc-healthd/internal/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func TestParseStatus_KnownValues(t *testing.T) {
	cases := []struct {
		input    string
		want     grpc_health_v1.HealthCheckResponse_ServingStatus
		wantOK   bool
	}{
		{"SERVING", grpc_health_v1.HealthCheckResponse_SERVING, true},
		{"NOT_SERVING", grpc_health_v1.HealthCheckResponse_NOT_SERVING, true},
		{"UNKNOWN", grpc_health_v1.HealthCheckResponse_UNKNOWN, true},
		{"SERVICE_UNKNOWN", grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN, true},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := health.ParseStatus(tc.input)
			if ok != tc.wantOK {
				t.Fatalf("ok: want %v, got %v", tc.wantOK, ok)
			}
			if got != tc.want {
				t.Fatalf("status: want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestParseStatus_UnknownValue(t *testing.T) {
	_, ok := health.ParseStatus("MAYBE")
	if ok {
		t.Fatal("expected ok=false for unknown status string")
	}
}

func TestParseStatus_EmptyString(t *testing.T) {
	_, ok := health.ParseStatus("")
	if ok {
		t.Fatal("expected ok=false for empty string")
	}
}

func TestParseStatus_CaseSensitive(t *testing.T) {
	_, ok := health.ParseStatus("serving")
	if ok {
		t.Fatal("expected ok=false for lowercase input")
	}
}
