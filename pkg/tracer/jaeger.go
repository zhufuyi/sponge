package tracer

import (
	"go.opentelemetry.io/otel/exporters/jaeger" //nolint
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

// JaegerOption set fields
type JaegerOption func(*jaegerOptions)

type jaegerOptions struct {
	username string
	password string
}

func (o *jaegerOptions) apply(opts ...JaegerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default setting
func defaultJaegerOptions() *jaegerOptions {
	return &jaegerOptions{}
}

// WithUsername set username
func WithUsername(username string) JaegerOption {
	return func(o *jaegerOptions) {
		o.username = username
	}
}

// WithPassword set password
func WithPassword(password string) JaegerOption {
	return func(o *jaegerOptions) {
		o.password = password
	}
}

// NewJaegerExporter use jaeger collector as exporter, e.g. default url=http://localhost:14268/api/traces
func NewJaegerExporter(url string, opts ...JaegerOption) (sdkTrace.SpanExporter, error) {
	ceps := []jaeger.CollectorEndpointOption{
		jaeger.WithEndpoint(url),
	}

	o := defaultJaegerOptions()
	o.apply(opts...)
	if o.username != "" {
		ceps = append(ceps, jaeger.WithUsername(o.username))
	}
	if o.password != "" {
		ceps = append(ceps, jaeger.WithPassword(o.password))
	}

	endpointOption := jaeger.WithCollectorEndpoint(ceps...)

	return jaeger.New(endpointOption)
}

// NewJaegerAgentExporter use jaeger agent as exporter, e.g. host=localhost port=6831
func NewJaegerAgentExporter(host string, port string) (sdkTrace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(host),
			jaeger.WithAgentPort(port),
		),
	)
}
