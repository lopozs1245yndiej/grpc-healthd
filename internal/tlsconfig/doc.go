// Package tlsconfig loads TLS credentials for grpc-healthd servers.
//
// Configuration is driven by two environment variables:
//
//	GRPC_HEALTHD_TLS_CERT – path to a PEM-encoded X.509 certificate
//	GRPC_HEALTHD_TLS_KEY  – path to the corresponding PEM-encoded private key
//
// When both variables are set, [Load] returns a *tls.Config that enforces
// TLS 1.2 as the minimum protocol version. When either variable is absent
// the daemon starts without TLS.
package tlsconfig
