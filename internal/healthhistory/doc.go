// Package healthhistory records a rolling window of health state
// transitions for each watched service.
//
// Events are stored in memory up to a configurable maximum. The Handler
// exposes them over HTTP as JSON, with optional filtering by service name.
package healthhistory
