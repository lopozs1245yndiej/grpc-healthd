// Package healthhistory provides a bounded in-memory ring-buffer that records
// health status transition events for monitored services.
//
// Events can be retrieved in full or filtered by service name via the HTTP
// handler exposed by Handler.
package healthhistory
