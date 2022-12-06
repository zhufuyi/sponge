package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var traceName = "unknown"

// SetTraceName each service corresponds to a traceName
func SetTraceName(name string) {
	if name != "" {
		traceName = name
	}
}

// NewSpan create a span, to end a span you must call span.End()
func NewSpan(ctx context.Context, spanName string, tags map[string]interface{}) (context.Context, trace.Span) {
	var opts []trace.SpanStartOption

	for k, v := range tags {
		var tag attribute.KeyValue
		switch v.(type) {
		case nil:
			continue
		case bool:
			tag = attribute.Bool(k, v.(bool))
		case string:
			tag = attribute.String(k, v.(string))
		case []string:
			tag = attribute.StringSlice(k, v.([]string))
		case int:
			tag = attribute.Int(k, v.(int))
		case []int:
			tag = attribute.IntSlice(k, v.([]int))
		case int64:
			tag = attribute.Int64(k, v.(int64))
		case []int64:
			tag = attribute.Int64Slice(k, v.([]int64))
		case float64:
			tag = attribute.Float64(k, v.(float64))
		case []float64:
			tag = attribute.Float64Slice(k, v.([]float64))
		default:
			tag = attribute.String(k, fmt.Sprintf("%+v", v))
		}
		opts = append(opts, trace.WithAttributes(tag))
	}

	return otel.Tracer(traceName).Start(ctx, spanName, opts...)
}
