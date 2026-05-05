// Package health provides a thread-safe service health status registry
// for use with gRPC health check endpoints.
//
// The Checker type maintains a map of named services and their current
// health statuses (SERVING, NOT_SERVING, or UNKNOWN). It is safe for
// concurrent use and is intended to be shared between the gRPC server
// handler and any background probing goroutines.
//
// Basic usage:
//
//	checker := health.NewChecker()
//	checker.SetStatus("my-service", health.StatusServing)
//	status := checker.GetStatus(ctx, "my-service")
package health
