package ratelimit

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds rate-limiter configuration sourced from environment variables.
type Config struct {
	// Rate is the number of tokens added per Interval (default: 10).
	Rate int
	// Interval is the token refill period (default: 1s).
	Interval time.Duration
	// Capacity is the maximum burst size (default: 20).
	Capacity int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Rate:     10,
		Interval: time.Second,
		Capacity: 20,
	}
}

// LoadFromEnv reads RATELIMIT_RATE, RATELIMIT_INTERVAL_MS, and
// RATELIMIT_CAPACITY from the environment, falling back to defaults.
func LoadFromEnv() (Config, error) {
	cfg := DefaultConfig()

	if v := os.Getenv("RATELIMIT_RATE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("ratelimit: invalid RATELIMIT_RATE %q: must be a positive integer", v)
		}
		cfg.Rate = n
	}

	if v := os.Getenv("RATELIMIT_INTERVAL_MS"); v != "" {
		ms, err := strconv.Atoi(v)
		if err != nil || ms <= 0 {
			return cfg, fmt.Errorf("ratelimit: invalid RATELIMIT_INTERVAL_MS %q: must be a positive integer", v)
		}
		cfg.Interval = time.Duration(ms) * time.Millisecond
	}

	if v := os.Getenv("RATELIMIT_CAPACITY"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("ratelimit: invalid RATELIMIT_CAPACITY %q: must be a positive integer", v)
		}
		cfg.Capacity = n
	}

	return cfg, nil
}

// NewFromConfig constructs a Limiter from the given Config.
func NewFromConfig(cfg Config) *Limiter {
	return New(cfg.Rate, cfg.Interval, cfg.Capacity)
}
