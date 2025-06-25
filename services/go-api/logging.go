package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

// LoggerWithTrace returns a logger enriched with trace and span IDs from ctx.
func LoggerWithTrace(ctx context.Context) zerolog.Logger {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	return log.With().
		Str("trace_id", sc.TraceID().String()).
		Str("span_id", sc.SpanID().String()).
		Logger()
}

// LoggingMiddleware logs each request with trace information.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := LoggerWithTrace(c.Request.Context())
		logger.Info().
			Str("method", c.Request.Method).
			Str("path", c.FullPath()).
			Msg("request started")

		c.Next()

		logger = LoggerWithTrace(c.Request.Context())
		logger.Info().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", c.FullPath()).
			Msg("request completed")
	}
}
