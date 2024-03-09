package tracer

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// alias, for other structs, the following code does not need to change the names of the resourceOptions
type resourceOptions = resourceConfig

// ResourceOption modifying struct field values by means of an interface
type ResourceOption interface {
	apply(*resourceOptions)
}

type resourceOptionFunc func(*resourceOptions)

func (o resourceOptionFunc) apply(cfg *resourceOptions) {
	o(cfg)
}

// set obj fields value
func apply(obj *resourceOptions, opts ...ResourceOption) {
	for _, opt := range opts {
		opt.apply(obj)
	}
}

// WithServiceName set service name
func WithServiceName(name string) ResourceOption {
	return resourceOptionFunc(func(o *resourceOptions) {
		o.serviceName = name
	})
}

// WithServiceVersion set service version
func WithServiceVersion(version string) ResourceOption {
	return resourceOptionFunc(func(o *resourceOptions) {
		o.serviceVersion = version
	})
}

// WithEnvironment set service environment
func WithEnvironment(environment string) ResourceOption {
	return resourceOptionFunc(func(o *resourceOptions) {
		o.environment = environment
	})
}

// WithAttributes set service attributes
func WithAttributes(attributes map[string]string) ResourceOption {
	return resourceOptionFunc(func(o *resourceOptions) {
		o.attributes = attributes
	})
}

type resourceConfig struct {
	serviceName    string
	serviceVersion string
	environment    string

	attributes map[string]string
}

// NewResource returns a resource describing this application.
func NewResource(opts ...ResourceOption) *resource.Resource {
	// default values
	rc := &resourceConfig{
		serviceName:    "demo-service",
		serviceVersion: "v0.0.0",
		environment:    "dev",
	}
	apply(rc, opts...)

	kvs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(rc.serviceName),
		semconv.ServiceVersionKey.String(rc.serviceVersion),
		attribute.String("env", rc.environment),
	}
	for k, v := range rc.attributes {
		kvs = append(kvs, attribute.String(k, v))
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, kvs...),
	)
	if err != nil {
		panic(err)
	}

	return r
}
