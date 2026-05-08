// Package admin provides HTTP administration endpoints for the grpc-healthd
// daemon. It exposes two routes:
//
//   - GET  /admin/services  – lists all registered services and their current
//     health status as a JSON array.
//
//   - POST /admin/set       – accepts a JSON body {"service": "…", "status": "…"}
//     and overrides the health status of the named service in the shared
//     health.Checker.
//
// The package is intentionally lightweight and has no external dependencies
// beyond the standard library and the internal health package.
package admin
