module github.com/l00p8/tracer

go 1.16

require (
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/exporters/jaeger v1.1.0
	go.opentelemetry.io/otel/sdk v1.1.0
	go.opentelemetry.io/otel/trace v1.1.0
)

require go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.26.1
