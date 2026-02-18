package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/atozi-ai/gateway/internal/handlers"
	"github.com/atozi-ai/gateway/internal/platform/logger"
	"github.com/atozi-ai/gateway/internal/ratelimit"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func getEnvFloat(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func Start() {
	logger.Init()
	logger.Log.Info().Msg("Initializing application...")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	rateLimitConfig := ratelimit.RateLimitConfig{
		RequestsPerSecond: getEnvFloat("RATE_LIMIT_REQUESTS_PER_SECOND", 10),
		RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 0),
		RequestsPerHour:   getEnvInt("RATE_LIMIT_REQUESTS_PER_HOUR", 0),
		RequestsPerDay:    getEnvInt("RATE_LIMIT_REQUESTS_PER_DAY", 0),
		Burst:             getEnvInt("RATE_LIMIT_BURST", 20),
	}

	rateLimiter := ratelimit.NewRateLimiter(rateLimitConfig)
	logger.Log.Info().
		Float64("requests_per_second", rateLimitConfig.RequestsPerSecond).
		Int("requests_per_minute", rateLimitConfig.RequestsPerMinute).
		Int("requests_per_hour", rateLimitConfig.RequestsPerHour).
		Int("requests_per_day", rateLimitConfig.RequestsPerDay).
		Int("burst", rateLimitConfig.Burst).
		Msg("Rate limiting enabled")

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	chatHandler := handlers.NewChatHandler()

	r.Route("/api/v1", func(r chi.Router) {
		ratelimit.RegisterRateLimiter(r, rateLimiter)
		chatHandler.RegisterRoutes(r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Log.Info().Str("port", port).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Log.Info().Msg("Server exited")

}
