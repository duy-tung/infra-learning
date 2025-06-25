package main

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func spansToTraces(spans []sdktrace.ReadOnlySpan) ptrace.Traces {
	td := ptrace.NewTraces()
	for _, s := range spans {
		rs := td.ResourceSpans().AppendEmpty()
		resAttrs := rs.Resource().Attributes()
		for _, kv := range s.Resource().Attributes() {
			setAttribute(resAttrs, kv)
		}
		rs.SetSchemaUrl(s.Resource().SchemaURL())
		ss := rs.ScopeSpans().AppendEmpty()
		scope := s.InstrumentationScope()
		ss.Scope().SetName(scope.Name)
		ss.Scope().SetVersion(scope.Version)
		sp := ss.Spans().AppendEmpty()
		sp.SetName(s.Name())
		sp.SetTraceID(pcommon.TraceID(s.SpanContext().TraceID()))
		sp.SetSpanID(pcommon.SpanID(s.SpanContext().SpanID()))
		if pid := s.Parent().SpanID(); pid.IsValid() {
			sp.SetParentSpanID(pcommon.SpanID(pid))
		}
		sp.SetStartTimestamp(pcommon.NewTimestampFromTime(s.StartTime()))
		sp.SetEndTimestamp(pcommon.NewTimestampFromTime(s.EndTime()))
		sp.SetKind(ptrace.SpanKind(s.SpanKind()))
		attrMap := sp.Attributes()
		for _, kv := range s.Attributes() {
			setAttribute(attrMap, kv)
		}
	}
	return td
}

func setAttribute(m pcommon.Map, kv attribute.KeyValue) {
	switch kv.Value.Type() {
	case attribute.BOOL:
		m.PutBool(string(kv.Key), kv.Value.AsBool())
	case attribute.INT64:
		m.PutInt(string(kv.Key), kv.Value.AsInt64())
	case attribute.FLOAT64:
		m.PutDouble(string(kv.Key), kv.Value.AsFloat64())
	case attribute.STRING:
		m.PutStr(string(kv.Key), kv.Value.AsString())
	default:
		m.PutStr(string(kv.Key), kv.Value.Emit())
	}
}
