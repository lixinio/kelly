package kelly

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMethod(t *testing.T) {

	f := func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			c.ResponseStatusOK()
		}
	}

	qhandler := func(method string, code int, register func(Router)) {
		k := New(nil).(*kellyImp)
		register(k)
		k.router.doPreRun()
		mux := http.NewServeMux()
		mux.HandleFunc("/", k.router.ServeHTTP)

		r, _ := http.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != code {
			t.Errorf("Method(%s) test error %d|%d", method, resp.StatusCode, code)
		}
	}

	methods := []string{
		http.MethodGet,
		http.MethodDelete,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}
	ops := map[string]func(Router){
		http.MethodGet:     func(r Router) { r.GET("/", f) },
		http.MethodDelete:  func(r Router) { r.DELETE("/", f) },
		http.MethodPost:    func(r Router) { r.POST("/", f) },
		http.MethodPut:     func(r Router) { r.PUT("/", f) },
		http.MethodPatch:   func(r Router) { r.PATCH("/", f) },
		http.MethodHead:    func(r Router) { r.HEAD("/", f) },
		http.MethodOptions: func(r Router) { r.OPTIONS("/", f) },
	}
	for _, m := range methods {
		for _, m2 := range methods {
			if m == m2 || m == http.MethodOptions {
				qhandler(m, http.StatusOK, ops[m2])
			} else if m2 == http.MethodOptions {
				qhandler(m, http.StatusNotFound, ops[m2])
			} else {
				qhandler(m, http.StatusMethodNotAllowed, ops[m2])
			}
		}
	}
}

func TestBreak(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	k1 := New(nil)
	k2 := New(nil)
	go func() {
		select {
		case <-time.After(time.Millisecond * 100):
			cancel()
		}
	}()

	err1 := k1.RunContext(ctx, ":54321")
	err2 := k2.RunContext(ctx, ":54322")
	if err1 != nil {
		t.Errorf("k1 return fail %s", err1)
	}
	if err2 != nil {
		t.Errorf("k2 return fail %s", err2)
	}
}
