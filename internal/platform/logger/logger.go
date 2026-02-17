package logger

import (
	"context"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339

	env := os.Getenv("APP_ENV")

	if env == "production" {
		Log = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		Log = zerolog.New(consoleWriter).
			With().
			Timestamp().
			Caller().
			Logger()
	}
}

// FromContext returns a logger with request ID from context if available
func FromContext(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return Log
	}
	
	requestID := middleware.GetReqID(ctx)
	if requestID != "" {
		return Log.With().Str("request_id", requestID).Logger()
	}
	
	return Log
}
