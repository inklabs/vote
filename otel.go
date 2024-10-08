package vote

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	sdkResource "go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewResource() *sdkResource.Resource {
	resource, _ := sdkResource.Merge(
		sdkResource.Default(),
		sdkResource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("voting-demo"),
			semconv.ServiceVersion("0.1.0"),
		))

	return resource
}

func NewOLTPConn() *grpc.ClientConn {
	conn, err := grpc.NewClient("localhost:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC client connection: %v", err)
	}

	return conn
}

func NewOTLPTracerProvider(resource *sdkResource.Resource, conn *grpc.ClientConn) *sdkTrace.TracerProvider {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("Failed to create metric exporter: %v", err)
	}

	batchSpanProcessor := sdkTrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithResource(resource),
		sdkTrace.WithSpanProcessor(batchSpanProcessor),
	)

	return tracerProvider
}

func NewOLTPMeterProvider(resource *sdkResource.Resource, conn *grpc.ClientConn) *sdkMetric.MeterProvider {
	ctx := context.Background()
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("Failed to create metric exporter: %v", err)
	}

	meterProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(resource),
		sdkMetric.WithReader(sdkMetric.NewPeriodicReader(exporter)),
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
