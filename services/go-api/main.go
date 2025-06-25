package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func setupTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
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
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	ctx := context.Background()
	tp, err := setupTracerProvider(ctx)
	if err != nil {
		log.Fatalf("failed to setup tracer provider: %v", err)
	}
	defer func() { _ = tp.Shutdown(ctx) }()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	r.Use(MetricsMiddleware())

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/example", func(c *gin.Context) {
		req, _ := http.NewRequestWithContext(c.Request.Context(), "GET", "https://example.com", nil)
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp.Body.Close()
		c.JSON(http.StatusOK, gin.H{"status": resp.Status})
	})

	if err := r.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
