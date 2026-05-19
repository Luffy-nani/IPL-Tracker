package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// This is for the logs--> Used for debugging...
// Same as OTEL
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		fmt.Printf("[%s] %s %s — %v\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path, time.Since(start))
	}
}

//Rate limiter

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

var limiter = &RateLimiter{
	visitors: make(map[string]*visitor),
	limit:    10,
	window:   time.Minute,
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]

	if !exists || time.Since(v.lastSeen) > rl.window {
		rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
		return true
	}

	if v.count >= rl.limit {
		return false
	}

	v.count++
	return true
}

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // to get the ip address
		if !limiter.allow(ip) {
			writeJSON(w, http.StatusTooManyRequests, map[string]string{
				"error": "rate limit exceeded, try again in a minute",
			})
			return
		}
		next(w, r)
	}
}

// Chaining the middleware
func Chain(h http.HandlerFunc) http.HandlerFunc {
	return loggingMiddleware(rateLimitMiddleware(h))
}
