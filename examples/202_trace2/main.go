package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/telemetry"
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

func subHandler(ctx context.Context, name string) {
	ctx, span := trace.StartSpan(ctx, name)
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

func Handler(name string) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		subHandler(c.Context(), name)
		time.Sleep(time.Millisecond * 100)

		// 调用http请求
		req, _ := http.NewRequest("GET", "http://127.0.0.1:9998/slave", nil)
		req = req.WithContext(c.Context())
		client := &http.Client{Transport: &ochttp.Transport{}}
		res, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to make the request: %v", err)
		}
		io.Copy(ioutil.Discard, res.Body)
		_ = res.Body.Close()

		c.ResponseStatusOK()
	}
}

func Handler2(name string) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		subHandler(c.Context(), name)
		time.Sleep(time.Millisecond * 500)
		c.ResponseStatusOK()
	}
}

func main() {
	initTrace()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	router1 := kelly.New(nil)
	router1.GET("/main", telemetry.OChttp, Handler("main_request"))
	router2 := kelly.New(nil)
	router2.GET("/slave", telemetry.OChttp, Handler2("sub_request"))

	go watchSignal(cancel)
	go func() {
		defer wg.Done()
		router1.RunContext(ctx, ":9999")
	}()
	go func() {
		defer wg.Done()
		router2.RunContext(ctx, ":9998")
	}()

	wg.Wait()
	fmt.Println("main finished")
}
