package config_test

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Port != 50051 {
		t.Errorf("expected default port 50051, got %d", cfg.Port)
	}
	if cfg.CheckInterval != 10*time.Second {
		t.Errorf("expected default check interval 10s, got %v", cfg.CheckInterval)
	}
	if cfg.ServicesFile == "" {
		t.Error("expected non-empty default services file path")
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log level 'info', got %q", cfg.LogLevel)
	}
}

func TestLoadFromEnv_Defaults(t *testing.T) {
	// No env vars set — should return defaults without error.
	cfg, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := config.DefaultConfig()
	if cfg.Port != def.Port {
		t.Errorf("port mismatch: got %d, want %d", cfg.Port, def.Port)
	}
}

func TestLoadFromEnv_CustomValues(t *testing.T) {
	t.Setenv("GRPC_HEALTHD_PORT", "9090")
	t.Setenv("GRPC_HEALTHD_CHECK_INTERVAL", "30s")
	t.Setenv("GRPC_HEALTHD_SERVICES_FILE", "/tmp/services.yaml")
	t.Setenv("GRPC_HEALTHD_LOG_LEVEL", "debug")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.CheckInterval != 30*time.Second {
		t.Errorf("expected interval 30s, got %v", cfg.CheckInterval)
	}
	if cfg.ServicesFile != "/tmp/services.yaml" {
		t.Errorf("unexpected services file: %s", cfg.ServicesFile)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected log level 'debug', got %q", cfg.LogLevel)
	}
}

func TestLoadFromEnv_InvalidPort(t *testing.T) {
	t.Setenv("GRPC_HEALTHD_PORT", "not-a-number")
	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error for invalid port, got nil")
	}
}

func TestLoadFromEnv_OutOfRangePort(t *testing.T) {
	t.Setenv("GRPC_HEALTHD_PORT", "99999")
	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error for out-of-range port, got nil")
	}
}

func TestLoadFromEnv_InvalidDuration(t *testing.T) {
	t.Setenv("GRPC_HEALTHD_CHECK_INTERVAL", "bad-duration")
	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}
}

func TestLoadFromEnv_ZeroDuration(t *testing.T) {
	t.Setenv("GRPC_HEALTHD_CHECK_INTERVAL", "0s")
	_, err := config.LoadFromEnv()
	if err == nil {
		t.Fatal("expected error for zero duration, got nil")
	}
}
