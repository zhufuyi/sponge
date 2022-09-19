package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var tp *trace.TracerProvider

// Init 初始化链路跟踪，fraction为分数，默认为1.0，值>=1.0表示全部链路都采样, 值<=0表示全部都不采样，0<值<1只采样百分比
func Init(exporter trace.SpanExporter, resource *resource.Resource, fractions ...float64) {
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
		trace.WithResource(resource),
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(fraction))), // 采样率
	)
	// 将TracerProvider注册为全局，这样将来任何导入包go.opentelemetry.io/otel/trace后，就可以默认使用它。
	otel.SetTracerProvider(tp)
	// 跨进程传播context
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

// Close 停止
func Close(ctx context.Context) error {
	if tp == nil {
		return nil
	}
	return tp.Shutdown(ctx)
}
