// Package watcher implements periodic health probing for registered services.
//
// A Probe implementation (e.g. HTTPChecker) is polled on a fixed interval by a
// ServiceWatcher. The result is written back into a health.Checker so that the
// gRPC health server always reflects the latest observed state.
//
// Multiple watchers can be managed together with Manager:
//
//	var mgr watcher.Manager
//	mgr.Add(watcher.NewServiceWatcher(
//		"my-service",
//		&watcher.HTTPChecker{URL: "http://localhost:8080/healthz"},
//		10*time.Second,
//		checker,
//	))
//	mgr.Start(ctx)
//	mgr.Wait()
package watcher
