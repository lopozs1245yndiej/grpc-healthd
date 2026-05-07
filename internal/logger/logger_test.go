package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/your-org/grpc-healthd/internal/logger"
)

func TestNew_ReturnsLogger(t *testing.T) {
	l := logger.New("info")
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNew_DefaultsToInfo(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithWriter("", &buf)

	l.Debug("should be suppressed")
	l.Info("should appear")

	lines := splitLines(buf.Bytes())
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line, got %d", len(lines))
	}
	assertMsg(t, lines[0], "should appear")
}

func TestNew_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithWriter("debug", &buf)

	l.Debug("debug message")
	l.Info("info message")

	lines := splitLines(buf.Bytes())
	if len(lines) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(lines))
	}
}

func TestNew_ErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithWriter("error", &buf)

	l.Info("suppressed")
	l.Warn("also suppressed")
	l.Error("visible")

	lines := splitLines(buf.Bytes())
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line, got %d", len(lines))
	}
	assertMsg(t, lines[0], "visible")
}

func TestNew_UnrecognisedLevelDefaultsToInfo(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewWithWriter("verbose", &buf)

	l.Debug("suppressed")
	l.Info("visible")

	lines := splitLines(buf.Bytes())
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line, got %d", len(lines))
	}
}

// splitLines returns non-empty lines from a byte slice.
func splitLines(b []byte) [][]byte {
	var out [][]byte
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(bytes.TrimSpace(line)) > 0 {
			out = append(out, line)
		}
	}
	return out
}

// assertMsg checks that a JSON log line contains the expected message.
func assertMsg(t *testing.T, line []byte, want string) {
	t.Helper()
	var entry map[string]any
	if err := json.Unmarshal(line, &entry); err != nil {
		t.Fatalf("failed to parse log line as JSON: %v", err)
	}
	if got, ok := entry["msg"].(string); !ok || got != want {
		t.Errorf("expected msg %q, got %q", want, got)
	}
}
