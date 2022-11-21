package tracer

import (
	"io"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewConsoleExporter output to console
func NewConsoleExporter() (sdkTrace.SpanExporter, error) {
	return stdout.New(stdout.WithPrettyPrint())
}

// NewFileExporter output to file, note: close the file before ending
func NewFileExporter(filename string) (sdkTrace.SpanExporter, *os.File, error) {
	if filename == "" {
		filename = "traces.json"
	}
	// Write telemetry data to a file.
	f, err := os.Create(filename)
	if err != nil {
		panic("os.Create error: " + err.Error())
	}

	exporter, err := newExporter(f)
	if err != nil {
		panic("newExporter error: " + err.Error())
	}

	return exporter, f, nil
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (sdkTrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
