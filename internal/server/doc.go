// Package server provides the gRPC server implementation for grpc-healthd.
//
// It wires together the gRPC health protocol (grpc.health.v1) with the
// internal health.Checker, exposing a Check endpoint that reflects the
// current status of each registered service.
//
// Usage:
//
//	cfg := config.DefaultConfig()
//	checker := health.NewChecker()
//	srv := server.New(cfg, checker)
//
//	// Register services before starting.
//	checker.SetStatus("my-service", grpc_health_v1.HealthCheckResponse_SERVING)
//
//	// Start blocks until Stop is called.
//	go srv.Start()
//	defer srv.Stop()
package server
