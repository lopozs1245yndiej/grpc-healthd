package version_test

import (
	"testing"

	"github.com/your-org/grpc-healthd/internal/version"
)

func TestParse_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  version.Version
	}{
		{"1.2.3", version.Version{Major: 1, Minor: 2, Patch: 3}},
		{"v1.2.3", version.Version{Major: 1, Minor: 2, Patch: 3}},
		{"0.0.0", version.Version{}},
		{"10.20.30", version.Version{Major: 10, Minor: 20, Patch: 30}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := version.Parse(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParse_Invalid(t *testing.T) {
	inputs := []string{"1.2", "1.2.x", "", "abc", "1.2.3.4"}
	for _, s := range inputs {
		t.Run(s, func(t *testing.T) {
			_, err := version.Parse(s)
			if err == nil {
				t.Errorf("Parse(%q): expected error, got nil", s)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	v := version.Version{Major: 2, Minor: 4, Patch: 1}
	if got := v.String(); got != "2.4.1" {
		t.Errorf("String() = %q, want %q", got, "2.4.1")
	}
}

func TestAtLeast(t *testing.T) {
	cases := []struct {
		v, other string
		want     bool
	}{
		{"2.0.0", "1.9.9", true},
		{"1.0.0", "1.0.0", true},
		{"1.0.0", "1.0.1", false},
		{"1.1.0", "1.0.9", true},
		{"0.0.1", "0.0.2", false},
	}
	for _, tc := range cases {
		t.Run(tc.v+"_"+tc.other, func(t *testing.T) {
			v, _ := version.Parse(tc.v)
			other, _ := version.Parse(tc.other)
			if got := v.AtLeast(other); got != tc.want {
				t.Errorf("AtLeast(%v, %v) = %v, want %v", v, other, got, tc.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		v, other string
		want     int
	}{
		{"1.0.0", "1.0.0", 0},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0", "2.0.0", -1},
	}
	for _, tc := range cases {
		t.Run(tc.v+"_vs_"+tc.other, func(t *testing.T) {
			v, _ := version.Parse(tc.v)
			other, _ := version.Parse(tc.other)
			if got := v.Compare(other); got != tc.want {
				t.Errorf("Compare(%v, %v) = %d, want %d", v, other, got, tc.want)
			}
		})
	}
}
