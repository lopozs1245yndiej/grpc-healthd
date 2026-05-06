package signals_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/example/grpc-healthd/internal/signals"
)

func TestNotifyContext_CancelledByFunction(t *testing.T) {
	ctx, stop := signals.NotifyContext(context.Background())
	defer stop()

	stop() // simulate manual cancellation

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after stop()")
	}
}

func TestNotifyContext_CancelledBySIGTERM(t *testing.T) {
	ctx, stop := signals.NotifyContext(context.Background())
	defer stop()

	// Send SIGTERM to ourselves.
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatalf("failed to send SIGTERM: %v", err)
	}

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after SIGTERM")
	}
}

func TestAwaitShutdown_ReturnsWhenContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		signals.AwaitShutdown(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// expected
	case <-time.After(time.Second):
		t.Fatal("AwaitShutdown did not return after context cancellation")
	}
}
