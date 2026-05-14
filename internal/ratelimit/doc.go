// Package ratelimit provides a lightweight, in-process token-bucket rate
// limiter designed for use with the grpc-healthd admin and metrics HTTP
// servers.
//
// Each unique key (typically a remote IP address) maintains its own token
// bucket. Tokens are replenished lazily when Allow is called, avoiding the
// need for background goroutines.
//
// Usage:
//
//	limiter := ratelimit.New(10, time.Second, 20)
//	http.Handle("/admin/", limiter.Middleware(adminHandler))
package ratelimit
