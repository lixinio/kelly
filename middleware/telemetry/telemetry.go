package telemetry

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/filters"
	"net/http"

	"github.com/lixinio/kelly"
)

type t struct {
	c *kelly.Context
}

func (t *t) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.c.SetRequest(r).SetResponseWriter(w)
	// 因为要计算执行的时间， 所以必须要在// ochttp.Handler 执行完后续所有的链路
	t.c.InvokeNext()
}

func Otelhttp(c *kelly.Context) {
	h := otelhttp.NewHandler(
		&t{c},
		c.Request().RequestURI,
		otelhttp.WithFilter(
			filters.None(
				filters.Path("/healthz"),
			),
		),
	)

	h.ServeHTTP(c.ResponseWriter, c.Request())
	// 这里无需再 InvokeNext
}
