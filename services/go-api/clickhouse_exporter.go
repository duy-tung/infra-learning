package main

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// clickhouseSpanExporter adapts the collector ClickHouse exporter to the OTel SDK.
type clickhouseSpanExporter struct {
	exp exporter.Traces
}

type nopHost struct{}

func (nopHost) GetExtensions() map[component.ID]component.Component { return nil }

func newClickhouseSpanExporter(ctx context.Context, cfg *clickhouseexporter.Config) (*clickhouseSpanExporter, error) {
	factory := clickhouseexporter.NewFactory()
	set := exporter.Settings{
		ID:                component.NewID(factory.Type()),
		TelemetrySettings: component.TelemetrySettings{Logger: zap.NewNop()},
		BuildInfo:         component.NewDefaultBuildInfo(),
	}
	exp, err := factory.CreateTraces(ctx, set, cfg)
	if err != nil {
		return nil, err
	}
	if err := exp.Start(ctx, nopHost{}); err != nil {
		return nil, err
	}
	return &clickhouseSpanExporter{exp: exp}, nil
}

func (e *clickhouseSpanExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	td := spansToTraces(spans)
	return e.exp.ConsumeTraces(ctx, td)
}

func (e *clickhouseSpanExporter) Shutdown(ctx context.Context) error {
	return e.exp.Shutdown(ctx)
}
