package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func setupTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	cfg := clickhouseexporter.NewFactory().CreateDefaultConfig().(*clickhouseexporter.Config)
	cfg.Endpoint = os.Getenv("CLICKHOUSE_ENDPOINT")
	if db := os.Getenv("CLICKHOUSE_DATABASE"); db != "" {
		cfg.Database = db
	}
	cfg.Username = os.Getenv("CLICKHOUSE_USERNAME")
	if pw := os.Getenv("CLICKHOUSE_PASSWORD"); pw != "" {
		cfg.Password = configopaque.String(pw)
	}
	chExporter, err := newClickhouseSpanExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}
	ratio := 1.0
	if v := os.Getenv("OTEL_SAMPLER_RATIO"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			ratio = f
		}
	}
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "go-api"
	}
	serviceVersion := os.Getenv("OTEL_SERVICE_VERSION")
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(chExporter),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func parseCustomTags() []attribute.KeyValue {
	var attrs []attribute.KeyValue
	if tags := os.Getenv("OTEL_CUSTOM_TAGS"); tags != "" {
		for _, p := range strings.Split(tags, ",") {
			kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
			if len(kv) == 2 {
				attrs = append(attrs, attribute.String(kv[0], kv[1]))
			}
		}
	}
	return attrs
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	ctx := context.Background()
	tp, err := setupTracerProvider(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup tracer provider")
	}
	defer func() { _ = tp.Shutdown(ctx) }()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	r.Use(MetricsMiddleware())
	r.Use(LoggingMiddleware())

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	promPort := os.Getenv("OTEL_EXPORTER_PROMETHEUS_PORT")
	if promPort == "" {
		promPort = "9464"
	}
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		addr := ":" + promPort
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("metrics server error")
		}
	}()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/example", func(c *gin.Context) {
		serviceName := os.Getenv("OTEL_SERVICE_NAME")
		if serviceName == "" {
			serviceName = "go-api"
		}
		attrs := []attribute.KeyValue{
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("operation", "example"),
		}
		attrs = append(attrs, parseCustomTags()...)
		ctx, span := otel.Tracer("example").Start(c.Request.Context(), "example")
		span.SetAttributes(attrs...)
		defer span.End()

		logger := LoggerWithTrace(ctx)
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
		resp, err := client.Do(req)
		if err != nil {
			logger.Error().Err(err).Msg("failed to call example.com")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp.Body.Close()
		logger.Info().Str("status", resp.Status).Msg("example.com response")
		c.JSON(http.StatusOK, gin.H{"status": resp.Status})
	})

	if err := r.Run(); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
