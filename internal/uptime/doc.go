// Package uptime provides a Tracker that records the daemon start time
// and computes the elapsed uptime. It also exposes an HTTP handler that
// returns a JSON snapshot containing the start timestamp and uptime in
// seconds, suitable for mounting on the admin HTTP server.
package uptime
