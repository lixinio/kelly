package telemetry

import (
	"net/http"

	"github.com/lixinio/kelly"
	"go.opencensus.io/plugin/ochttp"
)

type t struct {
	c *kelly.Context
}

func (t *t) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.c.SetRequest(r)
	// 因为要计算执行的时间， 所以必须要在// ochttp.Handler 执行完后续所有的链路
	t.c.InvokeNext()
}

func OChttp(ac *kelly.AnnotationContext) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		h := &ochttp.Handler{Handler: &t{c}}
		h.ServeHTTP(c.ResponseWriter, c.Request())
		// 这里无需再 InvokeNext
	}
}
