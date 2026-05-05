package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/health"
)

// Server wraps the gRPC server and health service handler.
type Server struct {
	grpcServer *grpc.Server
	checker    *health.Checker
	cfg        *config.Config
}

// New creates a new Server with the given config and checker.
func New(cfg *config.Config, checker *health.Checker) *Server {
	grpcServer := grpc.NewServer()
	s := &Server{
		grpcServer: grpcServer,
		checker:    checker,
		cfg:        cfg,
	}
	grpc_health_v1.RegisterHealthServer(grpcServer, &healthHandler{checker: checker})
	return s
}

// Start begins listening and serving gRPC requests.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	return s.grpcServer.Serve(lis)
}

// Stop gracefully shuts down the gRPC server.
func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

// healthHandler implements grpc_health_v1.HealthServer.
type healthHandler struct {
	grpc_health_v1.UnimplementedHealthServer
	checker *health.Checker
}

// Check returns the health status for the requested service.
func (h *healthHandler) Check(
	_ context.Context,
	req *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	status := h.checker.GetStatus(req.Service)
	return &grpc_health_v1.HealthCheckResponse{Status: status}, nil
}

// Watch is not supported and delegates to the unimplemented base.
func (h *healthHandler) Watch(
	req *grpc_health_v1.HealthCheckRequest,
	stream grpc_health_v1.Health_WatchServer,
) error {
	return h.UnimplementedHealthServer.Watch(req, stream)
}
