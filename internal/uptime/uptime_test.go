package uptime_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/grpc-healthd/internal/uptime"
)

func TestNew_StartTimeIsUTC(t *testing.T) {
	before := time.Now().UTC()
	tr := uptime.New()
	after := time.Now().UTC()

	st := tr.StartTime()
	if st.Location() != time.UTC {
		t.Errorf("expected UTC location, got %v", st.Location())
	}
	if st.Before(before) || st.After(after) {
		t.Errorf("start time %v not between %v and %v", st, before, after)
	}
}

func TestUptime_NonNegative(t *testing.T) {
	tr := uptime.New()
	time.Sleep(10 * time.Millisecond)
	if tr.Uptime() < 0 {
		t.Error("uptime should be non-negative")
	}
}

func TestUptime_IncreasesOverTime(t *testing.T) {
	start := time.Now().UTC().Add(-5 * time.Second)
	tr := uptime.NewWithTime(start)
	if tr.Uptime() < 4*time.Second {
		t.Errorf("expected uptime >= 4s, got %v", tr.Uptime())
	}
}

func TestSnapshot_Fields(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	tr := uptime.NewWithTime(start)

	snap := tr.Snapshot()
	if snap.StartTime != "2024-01-15T10:00:00Z" {
		t.Errorf("unexpected start_time: %s", snap.StartTime)
	}
	if snap.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds should be non-negative, got %d", snap.UptimeSeconds)
	}
}

func TestNewWithTime_PreservesUTC(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	local := time.Date(2024, 6, 1, 12, 0, 0, 0, loc)
	tr := uptime.NewWithTime(local)

	if tr.StartTime().Location() != time.UTC {
		t.Errorf("expected UTC, got %v", tr.StartTime().Location())
	}
}
