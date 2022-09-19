## tracer

在[go.opentelemetry.io/otel](go.opentelemetry.io/otel)基础上封装的链路跟踪库。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/tracer

<br>

## 使用示例

初始化trace，指定exporter和resource。

```go
func initTrace() {
    // exporter := tracer.NewConsoleExporter() // 输出到终端

    // exporter, f, err := tracer.NewFileExporter("trace.json") // 输出到文件

	// exporter, err := tracer.NewJaegerExporter("http://localhost:14268/api/traces") // 输出到jaeger，使用collector http
	exporter, err := tracer.NewJaegerAgentExporter("192.168.3.37", "6831") // 输出到jaeger，使用agent udp

	resource := tracer.NewResource(
		tracer.WithServiceName("your-service-name"),
		tracer.WithEnvironment("dev"),
		tracer.WithServiceVersion("demo"),
	)

	tracer.Init(exporter, resource) // 默认采集全部
	// tracer.Init(exporter, resource, 0.5) // 采集一半
}
```

<br>

在程序中创建一个span，ctx来源于上一个parent span。

```go
	_, span := otel.Tracer(serviceName).Start(
		ctx,
		spanName,
		trace.WithAttributes(attribute.String("foo", "bar")), // 自定义属性
	)
	defer span.End()

	// ......
```


<br>

documents https://opentelemetry.io/docs/instrumentation/go/

support OpenTelemetry in other libraries https://opentelemetry.io/registry/?language=go&component=instrumentation
