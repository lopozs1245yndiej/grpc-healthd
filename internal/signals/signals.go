// Package signals provides graceful shutdown handling for the daemon.
package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"
)

// NotifyContext returns a context that is cancelled when an OS termination
// signal (SIGINT or SIGTERM) is received. The returned stop function should
// be called to release resources when the context is no longer needed.
func NotifyContext(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
}

// AwaitShutdown blocks until ctx is done (e.g. a signal was received) and
// then logs the reason for shutdown.
func AwaitShutdown(ctx context.Context) {
	<-ctx.Done()
	slog.Info("shutdown signal received", "reason", context.Cause(ctx))
}
