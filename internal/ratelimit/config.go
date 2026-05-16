package ratelimit

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the token-bucket rate-limit parameters.
type Config struct {
	// Capacity is the maximum number of tokens in the bucket.
	Capacity int
	// RatePerSec is the number of tokens added per second.
	RatePerSec float64
	// Enabled controls whether rate limiting is active.
	Enabled bool
}

// DefaultConfig returns a Config with sensible production defaults.
func DefaultConfig() Config {
	return Config{
		Capacity:   100,
		RatePerSec: 10.0,
		Enabled:    true,
	}
}

// LoadFromEnv reads rate-limit configuration from environment variables,
// falling back to DefaultConfig values when variables are absent.
//
// Environment variables:
//
//	RATELIMIT_CAPACITY    – integer token-bucket capacity (default 100)
//	RATELIMIT_RATE_PER_SEC – float tokens added per second (default 10.0)
//	RATELIMIT_ENABLED     – "true"/"false" (default true)
func LoadFromEnv() (Config, error) {
	cfg := DefaultConfig()

	if v := os.Getenv("RATELIMIT_CAPACITY"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("ratelimit: invalid RATELIMIT_CAPACITY %q: must be a positive integer", v)
		}
		cfg.Capacity = n
	}

	if v := os.Getenv("RATELIMIT_RATE_PER_SEC"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || f <= 0 {
			return cfg, fmt.Errorf("ratelimit: invalid RATELIMIT_RATE_PER_SEC %q: must be a positive number", v)
		}
		cfg.RatePerSec = f
	}

	if v := os.Getenv("RATELIMIT_ENABLED"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return cfg, fmt.Errorf("ratelimit: invalid RATELIMIT_ENABLED %q: must be true or false", v)
		}
		cfg.Enabled = b
	}

	return cfg, nil
}

// NewFromConfig constructs a Limiter from the provided Config.
// If cfg.Enabled is false a no-op limiter is returned.
func NewFromConfig(cfg Config) *Limiter {
	if !cfg.Enabled {
		// Return a limiter with effectively unlimited capacity.
		return New(1<<31-1, float64(1<<31-1))
	}
	return New(cfg.Capacity, cfg.RatePerSec)
}
