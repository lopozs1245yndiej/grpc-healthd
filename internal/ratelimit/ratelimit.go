// Package ratelimit provides a simple token-bucket rate limiter for
// protecting admin and health endpoints from excessive request rates.
package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

// Limiter is a per-IP token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int           // tokens added per interval
	interval time.Duration // refill interval
	capacity int           // max tokens per bucket
}

type bucket struct {
	tokens   int
	lastSeen time.Time
}

// New creates a Limiter that allows rate requests per interval with the given
// burst capacity. A background goroutine is not required; refills are lazy.
func New(rate int, interval time.Duration, capacity int) *Limiter {
	return &Limiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		interval: interval,
		capacity: capacity,
	}
}

// Allow returns true if the given key (e.g. remote IP) is within the rate limit.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.capacity, lastSeen: now}
		l.buckets[key] = b
	}

	// Lazy refill: add tokens proportional to elapsed intervals.
	elapsed := now.Sub(b.lastSeen)
	if elapsed >= l.interval {
		added := int(elapsed/l.interval) * l.rate
		b.tokens += added
		if b.tokens > l.capacity {
			b.tokens = l.capacity
		}
		b.lastSeen = now
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// Middleware wraps an http.Handler and enforces the rate limit per remote address.
func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := remoteIP(r)
		if !l.Allow(ip) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// remoteIP extracts the IP portion from r.RemoteAddr.
func remoteIP(r *http.Request) string {
	addr := r.RemoteAddr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}
