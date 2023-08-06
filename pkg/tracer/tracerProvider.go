// Package tracer is a library wrapped in go.opentelemetry.io/otel.
package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var tp *trace.TracerProvider

// Init Initialize tracer, parameter fraction is fraction, default is 1.0, value >= 1.0 means all links are sampled,
// value <= 0 means all are not sampled, 0 < value < 1 only samples percentage
func Init(exporter trace.SpanExporter, res *resource.Resource, fractions ...float64) {
	var fraction = 1.0
	if len(fractions) > 0 {
		if fractions[0] <= 0 {
			fraction = 0
		} else if fractions[0] < 1 {
			fraction = fractions[0]
		}
	}

	tp = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(fraction))), // sampling rate
	)
	// register the TracerProvider as global so that any future imports of package go.opentelemetry.io/otel/trace will use it by default.
	otel.SetTracerProvider(tp)
	// propagation of context across processes
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

// Close tracer
func Close(ctx context.Context) error {
	if tp == nil {
		return nil
	}
	return tp.Shutdown(ctx)
}

// InitWithConfig Initialize tracer according to configuration, fraction is fraction, default is 1.0, value >= 1.0 means all links are sampled,
// value <= 0 means all are not sampled, 0 < value < 1 only samples percentage
func InitWithConfig(appName string, appEnv string, appVersion string,
	jaegerAgentHost string, jaegerAgentPort string, jaegerSamplingRate float64) {
	res := NewResource(
		WithServiceName(appName),
		WithEnvironment(appEnv),
		WithServiceVersion(appVersion),
	)

	// initializing tracing
	exporter, err := NewJaegerAgentExporter(jaegerAgentHost, jaegerAgentPort)
	if err != nil {
		panic("init trace error:" + err.Error())
	}

	Init(exporter, res, jaegerSamplingRate)

	SetTraceName(appName)
}
