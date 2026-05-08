package health

import (
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// statusNames maps canonical string representations to their proto enum values.
var statusNames = map[string]grpc_health_v1.HealthCheckResponse_ServingStatus{
	"UNKNOWN":     grpc_health_v1.HealthCheckResponse_UNKNOWN,
	"SERVING":     grpc_health_v1.HealthCheckResponse_SERVING,
	"NOT_SERVING": grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	"SERVICE_UNKNOWN": grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
}

// ParseStatus converts a string such as "SERVING" to the corresponding
// ServingStatus enum value. The second return value is false when the string
// does not match any known status.
func ParseStatus(s string) (grpc_health_v1.HealthCheckResponse_ServingStatus, bool) {
	v, ok := statusNames[s]
	return v, ok
}
