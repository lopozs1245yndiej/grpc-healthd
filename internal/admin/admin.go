package admin

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-healthd/internal/health"
)

// Handler exposes HTTP admin endpoints for managing service health states.
type Handler struct {
	checker *health.Checker
}

// New creates a new admin Handler backed by the given Checker.
func New(c *health.Checker) *Handler {
	return &Handler{checker: c}
}

// NewServer returns an *http.Server wired with the admin routes.
func NewServer(addr string, h *Handler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/services", h.ListServices)
	mux.HandleFunc("/admin/set", h.SetStatus)
	return &http.Server{Addr: addr, Handler: mux}
}

// serviceEntry is the JSON shape returned by ListServices.
type serviceEntry struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// ListServices handles GET /admin/services and returns all known services.
func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	svcs := h.checker.ListServices()
	entries := make([]serviceEntry, 0, len(svcs))
	for _, svc := range svcs {
		st, _ := h.checker.GetStatus(svc)
		entries = append(entries, serviceEntry{Service: svc, Status: st.String()})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

// setRequest is the expected JSON body for SetStatus.
type setRequest struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// SetStatus handles POST /admin/set and overrides a service's health status.
func (h *Handler) SetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req setRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	st, ok := health.ParseStatus(req.Status)
	if !ok {
		http.Error(w, "unknown status value", http.StatusBadRequest)
		return
	}
	h.checker.SetStatus(req.Service, st)
	w.WriteHeader(http.StatusNoContent)
}
