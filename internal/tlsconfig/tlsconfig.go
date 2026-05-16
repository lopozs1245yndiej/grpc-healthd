// Package tlsconfig provides helpers for loading TLS credentials
// from environment variables or explicit file paths for use with
// the gRPC server and admin HTTP server.
package tlsconfig

import (
	"crypto/tls"
	"fmt"
	"os"
)

// Config holds paths to TLS material.
type Config struct {
	// CertFile is the path to the PEM-encoded certificate.
	CertFile string
	// KeyFile is the path to the PEM-encoded private key.
	KeyFile string
	// Enabled reports whether TLS should be used.
	Enabled bool
}

// DefaultConfig returns a Config populated from environment variables.
//
//	GRPC_HEALTHD_TLS_CERT – path to certificate file (default: "")
//	GRPC_HEALTHD_TLS_KEY  – path to key file (default: "")
//
// TLS is considered enabled when both variables are non-empty.
func DefaultConfig() Config {
	cert := os.Getenv("GRPC_HEALTHD_TLS_CERT")
	key := os.Getenv("GRPC_HEALTHD_TLS_KEY")
	return Config{
		CertFile: cert,
		KeyFile:  key,
		Enabled:  cert != "" && key != "",
	}
}

// Load reads the certificate and key files and returns a *tls.Config
// suitable for use with either a gRPC or HTTP server.
// It returns an error when the files cannot be parsed.
func Load(cfg Config) (*tls.Config, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("tlsconfig: load key pair: %w", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
