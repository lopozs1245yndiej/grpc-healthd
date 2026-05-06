// Package signals provides OS signal handling utilities for grpc-healthd.
//
// It wraps the standard library's signal.NotifyContext to listen for
// SIGINT and SIGTERM, enabling clean graceful shutdown of the daemon and
// all its components (gRPC server, metrics server, service watchers).
//
// Typical usage:
//
//	ctx, stop := signals.NotifyContext(context.Background())
//	defer stop()
//	// ... start servers ...
//	signals.AwaitShutdown(ctx)
package signals
