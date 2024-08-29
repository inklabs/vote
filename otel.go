package vote

import (
	"log"
	"time"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	sdkResource "go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewResource() *sdkResource.Resource {
	resource, _ := sdkResource.Merge(sdkResource.Default(),
		sdkResource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("voting-demo"),
			semconv.ServiceVersion("0.1.0"),
		))

	return resource
}

func GetTracerProvider(resource *sdkResource.Resource) *sdkTrace.TracerProvider {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatalf("Failed to create Jaeger exporter: %v", err)
	}

	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithResource(resource),
		sdkTrace.WithBatcher(exporter),
	)

	return tracerProvider
}

func GetMeterProvider(resource *sdkResource.Resource) *sdkMetric.MeterProvider {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	meterProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(resource),
		sdkMetric.WithReader(exporter),
	)

	return meterProvider
}

func newStdoutTracerProvider() *sdkTrace.TracerProvider {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			sdkTrace.WithBatchTimeout(time.Second)),
	)
	return tracerProvider
}

func newStdoutMeterProvider() *sdkMetric.MeterProvider {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		log.Fatal(err)
	}

	meterProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(sdkMetric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			sdkMetric.WithInterval(3*time.Second))),
	)

	return meterProvider
}
