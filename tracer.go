package tracer

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Attrs          []attribute.KeyValue
	JaegerUrl      string
	Environment    string
}

// NewResource returns a resource describing application.
func newResource(serviceName string, serviceVersion string, env string, attrs ...attribute.KeyValue) (*resource.Resource, error) {
	if attrs == nil {
		attrs = []attribute.KeyValue{}
	}
	attrs = append([]attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
		semconv.DeploymentEnvironmentKey.String(env),
	}, attrs...)
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			attrs...,
		),
	)
}

func newJaegerExporter(url string) (trace.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}

func initProvider(exporter trace.SpanExporter, resource *resource.Resource) *trace.TracerProvider {
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp
}

func InitProvider(cfg *Config) error {
	exp, err := newJaegerExporter(cfg.JaegerUrl)
	if err != nil {
		return err
	}
	res, err := newResource(cfg.ServiceName, cfg.ServiceVersion, cfg.Environment, cfg.Attrs...)
	if err != nil {
		return err
	}
	_ = initProvider(exp, res)
	return nil
}
