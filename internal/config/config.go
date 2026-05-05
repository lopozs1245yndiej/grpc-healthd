package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the runtime configuration for grpc-healthd.
type Config struct {
	// Port on which the gRPC health server listens.
	Port int

	// CheckInterval is how often registered services are probed.
	CheckInterval time.Duration

	// ServicesFile is the path to the YAML file listing services to monitor.
	ServicesFile string

	// LogLevel controls verbosity: debug, info, warn, error.
	LogLevel string
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Port:          50051,
		CheckInterval: 10 * time.Second,
		ServicesFile:  "/etc/grpc-healthd/services.yaml",
		LogLevel:      "info",
	}
}

// LoadFromEnv overrides defaults with values from environment variables.
//
// Supported variables:
//
//	GRPC_HEALTHD_PORT          – integer listen port
//	GRPC_HEALTHD_CHECK_INTERVAL – duration string (e.g. "30s")
//	GRPC_HEALTHD_SERVICES_FILE  – path to services YAML
//	GRPC_HEALTHD_LOG_LEVEL      – log level string
func LoadFromEnv() (*Config, error) {
	cfg := DefaultConfig()

	if v := os.Getenv("GRPC_HEALTHD_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid GRPC_HEALTHD_PORT %q: %w", v, err)
		}
		if p < 1 || p > 65535 {
			return nil, fmt.Errorf("GRPC_HEALTHD_PORT %d out of range [1, 65535]", p)
		}
		cfg.Port = p
	}

	if v := os.Getenv("GRPC_HEALTHD_CHECK_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid GRPC_HEALTHD_CHECK_INTERVAL %q: %w", v, err)
		}
		if d <= 0 {
			return nil, fmt.Errorf("GRPC_HEALTHD_CHECK_INTERVAL must be positive")
		}
		cfg.CheckInterval = d
	}

	if v := os.Getenv("GRPC_HEALTHD_SERVICES_FILE"); v != "" {
		cfg.ServicesFile = v
	}

	if v := os.Getenv("GRPC_HEALTHD_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}

	return cfg, nil
}
