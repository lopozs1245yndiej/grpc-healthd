// Package main is the entry point for the grpc-healthd sidecar daemon.
// It wires together configuration, logging, health checking, gRPC serving,
// metrics exposition, admin HTTP API, and graceful shutdown.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/your-org/grpc-healthd/internal/admin"
	"github.com/your-org/grpc-healthd/internal/config"
	"github.com/your-org/grpc-healthd/internal/health"
	"github.com/your-org/grpc-healthd/internal/logger"
	"github.com/your-org/grpc-healthd/internal/metrics"
	"github.com/your-org/grpc-healthd/internal/server"
	"github.com/your-org/grpc-healthd/internal/signals"
	"github.com/your-org/grpc-healthd/internal/watcher"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "grpc-healthd: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration from environment variables.
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Initialise structured logger.
	log := logger.New(cfg.LogLevel)
	log.Info("starting grpc-healthd",
		"grpc_addr", cfg.GRPCAddr,
		"metrics_addr", cfg.MetricsAddr,
		"admin_addr", cfg.AdminAddr,
	)

	// Root context that is cancelled on SIGTERM / SIGINT.
	ctx, stop := signals.NotifyContext(context.Background())
	defer stop()

	// Health checker holds per-service status.
	checker := health.NewChecker()

	// Prometheus metrics.
	m := metrics.New(nil)

	// gRPC health server.
	grpcSrv, err := server.New(cfg.GRPCAddr, checker)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	// Admin HTTP server (list/set service status).
	adminHandler := admin.New(checker)
	adminSrv := admin.NewServer(cfg.AdminAddr, adminHandler)

	// Metrics HTTP server.
	metricsSrv := metrics.NewServer(cfg.MetricsAddr, m)

	// Service watcher — polls upstream health endpoints and updates checker.
	sw := watcher.NewServiceWatcher(checker, m, log)
	for _, svc := range cfg.Services {
		if err := sw.Register(svc); err != nil {
			log.Warn("failed to register service watcher", "service", svc.Name, "error", err)
		}
	}

	// Start background components.
	errCh := make(chan error, 3)

	go func() {
		log.Info("gRPC server listening", "addr", cfg.GRPCAddr)
		if err := grpcSrv.Serve(); err != nil {
			errCh <- fmt.Errorf("grpc server: %w", err)
		}
	}()

	go func() {
		log.Info("metrics server listening", "addr", cfg.MetricsAddr)
		if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("metrics server: %w", err)
		}
	}()

	go func() {
		log.Info("admin server listening", "addr", cfg.AdminAddr)
		if err := adminSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("admin server: %w", err)
		}
	}()

	go sw.Run(ctx)

	// Block until shutdown signal or a server error.
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received, stopping")
	}

	// Graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	grpcSrv.GracefulStop()

	if err := metricsSrv.Shutdown(shutdownCtx); err != nil {
		log.Warn("metrics server shutdown error", "error", err)
	}
	if err := adminSrv.Shutdown(shutdownCtx); err != nil {
		log.Warn("admin server shutdown error", "error", err)
	}

	log.Info("grpc-healthd stopped")
	return nil
}
