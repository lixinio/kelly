package main

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"github.com/lixinio/kelly"
	"go.opentelemetry.io/otel/sdk/trace"
)

// # http://localhost:16686/
// $ docker run -d --name jaeger \
//   -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
//   -p 5775:5775/udp \
//   -p 6831:6831/udp \
//   -p 6832:6832/udp \
//   -p 5778:5778 \
//   -p 16686:16686 \
//   -p 14268:14268 \
//   -p 14250:14250 \
//   -p 9411:9411 \
//   jaegertracing/all-in-one:1.22

func subHandler(ctx context.Context) {
	ctx, span := otel.Tracer("/bar").Start(ctx, "Run")
	defer span.End()

	// 6. Set status upon error
	span.SetStatus(codes.Unset,"error")

	span.SetAttributes(
		attribute.String("k", "v"),
		attribute.String("k2", "v2"),
	)
}

func Handler(c *kelly.Context) {
	subHandler(c.Context())
	c.ResponseStatusOK()
}

func main() {
	exporter, err:= jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint("http://localhost:14268/api/traces"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("server name"),
			//attribute.String("environment", environment),
			//attribute.Int64("ID", id),
		)),
		trace.WithSampler(trace.TraceIDRatioBased(0.01)),
	)
	otel.SetTracerProvider(provider)

	h := &otelhttp.Handler{}
	r := kelly.New(nil, h)
	r.GET("/t", Handler)

	log.Fatal(http.ListenAndServe(":9999", h))
}
