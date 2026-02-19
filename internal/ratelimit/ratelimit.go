package ratelimit

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/atozi-ai/gateway/internal/platform/logger"
	"github.com/go-chi/chi/v5"
	"golang.org/x/time/rate"
)

type RateLimitConfig struct {
	RequestsPerSecond float64
	RequestsPerMinute int
	RequestsPerHour   int
	RequestsPerDay    int
	Burst             int
	MaxClients        int // Maximum number of unique clients to track
}

type RateLimiter struct {
	clients         map[string]*clientLimiter
	mu              sync.RWMutex
	config          RateLimitConfig
	cleanupInterval time.Duration
}

type clientLimiter struct {
	secondLimiter *rate.Limiter
	minuteWindow  *windowCounter
	hourWindow    *windowCounter
	dayWindow     *windowCounter
	lastSeen      time.Time
}

type windowCounter struct {
	count     int
	resetTime time.Time
	mu        sync.Mutex
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		clients:         make(map[string]*clientLimiter),
		config:          config,
		cleanupInterval: 5 * time.Minute,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) getClient(key string) *clientLimiter {
	rl.mu.RLock()
	if client, exists := rl.clients[key]; exists {
		client.lastSeen = time.Now()
		rl.mu.RUnlock()
		return client
	}
	rl.mu.RUnlock()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check after acquiring write lock
	if client, exists := rl.clients[key]; exists {
		client.lastSeen = time.Now()
		return client
	}

	// Evict oldest clients if at capacity
	if rl.config.MaxClients > 0 && len(rl.clients) >= rl.config.MaxClients {
		rl.evictOldest(rl.config.MaxClients / 4) // Evict 25% of capacity
	}

	now := time.Now()
	client := &clientLimiter{
		secondLimiter: rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.Burst),
		minuteWindow: &windowCounter{
			resetTime: now.Truncate(time.Minute).Add(time.Minute),
		},
		hourWindow: &windowCounter{
			resetTime: now.Truncate(time.Hour).Add(time.Hour),
		},
		dayWindow: &windowCounter{
			resetTime: now.Truncate(24 * time.Hour).Add(24 * time.Hour),
		},
		lastSeen: now,
	}
	rl.clients[key] = client

	return client
}

func (rl *RateLimiter) evictOldest(count int) {
	type clientEntry struct {
		key      string
		lastSeen time.Time
	}

	entries := make([]clientEntry, 0, len(rl.clients))
	for k, c := range rl.clients {
		entries = append(entries, clientEntry{key: k, lastSeen: c.lastSeen})
	}

	// Sort by lastSeen ascending (oldest first)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].lastSeen.After(entries[j].lastSeen) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Remove oldest clients
	for i := 0; i < count && i < len(entries); i++ {
		delete(rl.clients, entries[i].key)
	}

	if count > 0 && len(entries) > 0 {
		logger.Log.Info().
			Int("evicted", min(count, len(entries))).
			Int("remaining", len(rl.clients)).
			Msg("Rate limiter evicted old clients")
	}
}

func (rl *RateLimiter) Allow(key string) (bool, string) {
	client := rl.getClient(key)

	if !client.secondLimiter.Allow() {
		return false, "Rate limit exceeded (per second)"
	}

	if rl.config.RequestsPerMinute > 0 {
		allowed, reason := client.minuteWindow.allow(rl.config.RequestsPerMinute, time.Minute)
		if !allowed {
			return false, "Rate limit exceeded (per minute): " + reason
		}
	}

	if rl.config.RequestsPerHour > 0 {
		allowed, reason := client.hourWindow.allow(rl.config.RequestsPerHour, time.Hour)
		if !allowed {
			return false, "Rate limit exceeded (per hour): " + reason
		}
	}

	if rl.config.RequestsPerDay > 0 {
		allowed, reason := client.dayWindow.allow(rl.config.RequestsPerDay, 24*time.Hour)
		if !allowed {
			return false, "Rate limit exceeded (per day): " + reason
		}
	}

	return true, ""
}

func (wc *windowCounter) allow(limit int, window time.Duration) (bool, string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	now := time.Now()

	if now.After(wc.resetTime) {
		wc.count = 1
		wc.resetTime = now.Truncate(window).Add(window)
		return true, ""
	}

	if wc.count >= limit {
		retryAfter := wc.resetTime.Sub(now).Seconds()
		return false, "retry after " + formatRetryAfter(retryAfter)
	}

	wc.count++
	return true, ""
}

func formatRetryAfter(seconds float64) string {
	if seconds < 60 {
		return "< 1 minute"
	}
	if seconds < 3600 {
		return string(rune(int(seconds/60))) + " minutes"
	}
	return string(rune(int(seconds/3600))) + " hours"
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for key, client := range rl.clients {
			if time.Since(client.lastSeen) > 10*time.Minute {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimit(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := extractAPIKey(r)
			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			allowed, reason := rl.Allow(apiKey)
			if !allowed {
				logger.Log.Warn().
					Str("api_key", truncate(apiKey, 8)).
					Str("reason", reason).
					Msg("Rate limit exceeded")
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":{"message":"Rate limit exceeded","type":"rate_limit_error","code":"rate_limit_exceeded"}}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractAPIKey(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const bearerPrefix = "Bearer "
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return strings.TrimPrefix(authHeader, bearerPrefix)
	}

	return authHeader
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

func RegisterRateLimiter(r chi.Router, rl *RateLimiter) {
	r.Use(RateLimit(rl))
}
