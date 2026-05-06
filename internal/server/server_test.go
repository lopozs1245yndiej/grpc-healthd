package server_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/health"
	"github.com/yourorg/grpc-healthd/internal/server"
)

func freePort(t *testing.T) int {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port
	lis.Close()
	return port
}

func startServer(t *testing.T) (*server.Server, *health.Checker, int) {
	t.Helper()
	port := freePort(t)
	cfg := &config.Config{Host: "127.0.0.1", Port: port}
	checker := health.NewChecker()
	srv := server.New(cfg, checker)

	go func() {
		if err := srv.Start(); err != nil {
			// server stopped — expected on GracefulStop
		}
	}()
	time.Sleep(50 * time.Millisecond)
	return srv, checker, port
}

func clientConn(t *testing.T, port int) *grpc.ClientConn {
	t.Helper()
	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	return conn
}

func TestServer_CheckUnknownService(t *testing.T) {
	srv, _, port := startServer(t)
	defer srv.Stop()

	conn := clientConn(t, port)
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "unknown"})
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_UNKNOWN {
		t.Errorf("expected UNKNOWN, got %v", resp.Status)
	}
}

func TestServer_CheckServingService(t *testing.T) {
	srv, checker, port := startServer(t)
	defer srv.Stop()

	checker.SetStatus("my-service", grpc_health_v1.HealthCheckResponse_SERVING)

	conn := clientConn(t, port)
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "my-service"})
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Errorf("expected SERVING, got %v", resp.Status)
	}
}

func TestServer_CheckNotServingService(t *testing.T) {
	srv, checker, port := startServer(t)
	defer srv.Stop()

	checker.SetStatus("degraded-service", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	conn := clientConn(t, port)
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "degraded-service"})
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Errorf("expected NOT_SERVING, got %v", resp.Status)
	}
}
