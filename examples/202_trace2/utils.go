package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
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
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: "http://localhost:14268/api/traces",
		Process: jaeger.Process{
			ServiceName: "kelly3-demo",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}
