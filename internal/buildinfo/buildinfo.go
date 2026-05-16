// Package buildinfo exposes compile-time version metadata for grpc-healthd.
package buildinfo

import "time"

// Variables are set at build time via -ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Info holds the build-time metadata for the running binary.
type Info struct {
	Version   string    `json:"version"`
	Commit    string    `json:"commit"`
	BuildDate string    `json:"build_date"`
	GoVersion string    `json:"go_version"`
	Retrieved time.Time `json:"retrieved_at"`
}

// Get returns the current build information.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: goVersion(),
		Retrieved: time.Now().UTC(),
	}
}
