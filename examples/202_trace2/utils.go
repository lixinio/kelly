package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

func waitSignal(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stop)
}

func watchSignal(cancel context.CancelFunc) {
	stop := make(chan struct{})
	waitSignal(stop)

	<-stop
	fmt.Println("Signal caught; shutting down server")
	cancel()
}

func initTrace() {
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
}
