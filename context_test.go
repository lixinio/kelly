package kelly

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func getKelly(group, append bool, middlewares ...interface{}) (*kellyImp, Router) {
	var k *kellyImp = nil
	var rt Router = nil
	if !group {
		if !append {
			k = New(nil, middlewares...).(*kellyImp)
		} else {
			k = New(nil).(*kellyImp)
			k.Use(middlewares...)
		}
		rt = k
	} else {
		k = New(nil).(*kellyImp)
		if !append {
			rt = k.Group("/", middlewares...)
		} else {
			rt = k.Group("/")
			rt.Use(middlewares...)
		}
	}
	return k, rt
}

func kellyMiddleware(group, append bool, handler AnnotationHandlerFunc, middlewares ...interface{}) *http.Response {
	k, rt := getKelly(group, append, middlewares...)
	rt.GET("/", handler)
	k.router.doPreRun()
	mux := http.NewServeMux()
	mux.HandleFunc("/", k.router.ServeHTTP)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Result()
}

func TestContext(t *testing.T) {
	key := "key"
	value := "value"
	for _, g := range []bool{false, true} {
		for _, a := range []bool{false, true} {
			resp := kellyMiddleware(
				g,
				a,
				func(ac *AnnotationContext) HandlerFunc {
					return func(c *Context) {
						v := c.Get(key)
						if v == nil {
							c.InvokeNext()
							return
						}
						t.Errorf("get value unexpect")
					}
				},
				func(c *Context) {
					v := c.Get(key)
					if v != nil {
						t.Errorf("get value unexpect")
						return
					}

					v = c.GetDefault(key, value)
					if v != value {
						t.Errorf("Get fail %w", v)
						return
					}

					c.InvokeNext()
				},
				func(c *Context) {
					c.Set(key, value)
					c.InvokeNext()
				},
				func(ac *AnnotationContext) HandlerFunc {
					return func(c *Context) {
						v := c.Get(key)
						if v == nil {
							t.Errorf("Get fail %w", v)
							return
						}
						c.MustGet(key)
						c.InvokeNext()
					}
				},
				func(c *Context) {
					v := c.Get(key)
					if v == nil {
						t.Errorf("Get fail %w", v)
						return
					}
					c.MustGet(key)

					defer checkError(t, ErrNoContextData)
					c.MustGet(key + " invalid")
				},
			)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("StatusOK error, %v", resp.StatusCode)
				return
			}
		}
	}
}
