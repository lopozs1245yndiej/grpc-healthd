// Package logger provides a structured logger for grpc-healthd.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Level constants mirror slog levels for convenience.
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// New creates a new *slog.Logger configured from the given level string.
// Accepted values (case-insensitive): "debug", "info", "warn", "error".
// Defaults to Info when the value is empty or unrecognised.
func New(levelStr string) *slog.Logger {
	level := parseLevel(levelStr)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

// NewWithWriter creates a logger that writes JSON to the supplied writer.
// Useful for tests that want to capture log output.
func NewWithWriter(levelStr string, w interface{ Write([]byte) (int, error) }) *slog.Logger {
	level := parseLevel(levelStr)
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

// parseLevel converts a string to a slog.Level, defaulting to Info.
func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
