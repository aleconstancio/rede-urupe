/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package api

import (
	"crypto/subtle"
	"net/http"
	"sync"
	"time"
)

// requireAuth validates Bearer token against ADMIN_API_KEY.
// If AdminAPIKey is empty, auth is bypassed (open mode).
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.AdminAPIKey == "" {
			next(w, r)
			return
		}
		token := r.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		if subtle.ConstantTimeCompare([]byte(token), []byte(s.cfg.AdminAPIKey)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// RateLimiter implements per-key token bucket rate limiting.
type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*bucket
	limit    int
	interval time.Duration
}

type bucket struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter creates a rate limiter allowing `limit` requests per `interval`.
func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*bucket),
		limit:    limit,
		interval: interval,
	}
}

// Allow checks if a request from key is within rate limits.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.clients[key]
	now := time.Now()
	if !exists || now.Sub(b.lastReset) > rl.interval {
		rl.clients[key] = &bucket{tokens: rl.limit - 1, lastReset: now}
		return true
	}
	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// Middleware wraps an http.HandlerFunc with rate limiting.
func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			ip = fwd
		}
		if !rl.Allow(ip) {
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
