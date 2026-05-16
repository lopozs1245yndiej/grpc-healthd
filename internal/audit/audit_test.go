package audit_test

import (
	"testing"
	"time"

	"github.com/example/grpc-healthd/internal/audit"
)

func TestNew_DefaultMaxSize(t *testing.T) {
	l := audit.New(0)
	if l == nil {
		t.Fatal("expected non-nil Log")
	}
}

func TestRecord_SingleEntry(t *testing.T) {
	l := audit.New(10)
	l.Record("127.0.0.1", "set_status", "svc-a", "SERVING")

	if l.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", l.Len())
	}

	entries := l.Entries()
	e := entries[0]
	if e.RemoteIP != "127.0.0.1" {
		t.Errorf("RemoteIP: got %q, want %q", e.RemoteIP, "127.0.0.1")
	}
	if e.Action != "set_status" {
		t.Errorf("Action: got %q, want %q", e.Action, "set_status")
	}
	if e.Service != "svc-a" {
		t.Errorf("Service: got %q, want %q", e.Service, "svc-a")
	}
	if e.Detail != "SERVING" {
		t.Errorf("Detail: got %q, want %q", e.Detail, "SERVING")
	}
	if e.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	before := time.Now().UTC()
	l := audit.New(10)
	l.Record("::1", "list_services", "", "")
	after := time.Now().UTC()

	e := l.Entries()[0]
	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Errorf("Timestamp %v not between %v and %v", e.Timestamp, before, after)
	}
	if e.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC location, got %v", e.Timestamp.Location())
	}
}

func TestRecord_Eviction(t *testing.T) {
	l := audit.New(3)
	for i := 0; i < 5; i++ {
		l.Record("10.0.0.1", "set_status", "svc", "SERVING")
	}
	if l.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", l.Len())
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	l := audit.New(10)
	l.Record("1.2.3.4", "action", "svc", "detail")

	a := l.Entries()
	a[0].Action = "tampered"

	b := l.Entries()
	if b[0].Action == "tampered" {
		t.Error("Entries() should return an independent copy")
	}
}

func TestLen_Empty(t *testing.T) {
	l := audit.New(10)
	if l.Len() != 0 {
		t.Errorf("expected 0, got %d", l.Len())
	}
}
