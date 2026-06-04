package telemetry

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/labstack/echo/v4"
)

// InitTelemetry initializes standard OpenTelemetry TracerProvider and registers it globally.
// It returns a shutdown cleanup function to flush traces on server termination.
func InitTelemetry(ctx context.Context) (func(), error) {
	// Only initialize and export traces if explicitly enabled or if OTLP endpoint env is set.
	// This mutes background warning connection refused logs in standard local environments.
	if os.Getenv("OTEL_ENABLED") != "true" && os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		log.Println("[OTel Engine] OpenTelemetry is currently disabled. To enable APM tracing, set OTEL_ENABLED=true.")
		return func() {}, nil
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "northstar-cmdb"
	}

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318" // Standard local collector HTTP endpoint port
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String("development"),
		),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, err
	}

	// Create OTLP HTTP trace exporter pointing to endpoint
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}
	// Support unencrypted local connections by default unless specifically overriding
	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") != "false" {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	traceExporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		log.Printf("[OTel Engine] Warning: Failed to initialize OTLP exporter: %v. Running with no-op exporter.", err)
		return func() {}, nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	shutdown := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("[OTel Engine] Error shutting down TracerProvider: %v", err)
		}
	}

	log.Printf("[OTel Engine] OpenTelemetry TracerProvider successfully initialized for service: %v", serviceName)
	return shutdown, nil
}

// Middleware returns an Echo middleware for tracing incoming HTTP requests.
func Middleware() echo.MiddlewareFunc {
	return otelecho.Middleware("northstar-cmdb")
}
