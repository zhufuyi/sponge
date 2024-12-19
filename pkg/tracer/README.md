## tracer

Tracer library wrapped in [go.opentelemetry.io/otel](https://github.com/open-telemetry/opentelemetry-go).

<br>

## Example of use

Initialize the trace, specifying exporter and resource.

```go
import "github.com/go-dev-frame/sponge/pkg/tracer"

func initTrace() {
	// exporter := tracer.NewConsoleExporter() // output to terminal

	// exporter, f, err := tracer.NewFileExporter("trace.json") // output to file

	// exporter, err := tracer.NewJaegerExporter("http://localhost:14268/api/traces") // output to jaeger, using collector http
	exporter, err := tracer.NewJaegerAgentExporter("192.168.3.37", "6831") // output to jaeger, using agent udp

	resource := tracer.NewResource(
		tracer.WithServiceName("your-service-name"),
		tracer.WithEnvironment("dev"),
		tracer.WithServiceVersion("demo"),
	)

	tracer.Init(exporter, resource) // collect all by default
	// tracer.Init(exporter, resource, 0.5) // collect half
}
```

<br>

Create a span in the program with ctx derived from the previous parent span.

```go
	_, span := otel.Tracer(serviceName).Start(
		ctx,
		spanName,
		trace.WithAttributes(attribute.String("foo", "bar")), // customised attributes
	)
	defer span.End()

	// ......
```


<br>

documents https://opentelemetry.io/docs/instrumentation/go/

support OpenTelemetry in other libraries https://opentelemetry.io/registry/?language=go&component=instrumentation
