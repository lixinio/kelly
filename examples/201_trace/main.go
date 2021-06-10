package main

import (
	"context"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/lixinio/kelly"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
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
	ctx, span := trace.StartSpan(ctx, "/bar")
	defer span.End()

	// 6. Set status upon error
	span.SetStatus(trace.Status{
		Code:    trace.StatusCodeUnknown,
		Message: "error",
	})

	span.AddAttributes(
		trace.StringAttribute("k", "v"),
		trace.StringAttribute("k2", "v2"),
	)

	// 7. Annotate our span to capture metadata about our operation
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("bytes to int", 23),
	}, "Invoking doWork")

	ctx, _ = tag.New(
		ctx,
		tag.Upsert(tag.MustNewKey("func"), "sub func"),
	)
}

func Handler(c *kelly.Context) {
	subHandler(c.Context())
	c.ResponseStatusOK()
}

func main() {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: "http://localhost:14268/api/traces",
		Process: jaeger.Process{
			ServiceName: "kelly-demo",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	r := kelly.New(nil)
	r.GET("/t", Handler)
	h := &ochttp.Handler{Handler: r}

	log.Fatal(http.ListenAndServe(":9999", h))
}
