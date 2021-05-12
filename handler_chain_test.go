package kelly

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"
)

func wrapHandlerChainHandler(chain *HandlerChain) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chain.ServeHTTP(w, r, httprouter.Params{})
	}
}

func TestHandlerChain(t *testing.T) {
	httpExec := func(chain *HandlerChain) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", wrapHandlerChainHandler(chain))
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
	}

	f := func(p func(c *Context)) func(*Context) {
		return func(c *Context) {
			p(c)
		}
	}

	result := []string{}

	//-------------------------------------------------
	chain := newHandlerChain(
		f(func(c *Context) {
			result = append(result, "1")
			c.InvokeNext()
		}),
		f(func(c *Context) {
			result = append(result, "2")
			c.InvokeNext()
		}),
		f(func(c *Context) {
			result = append(result, "3")
			c.InvokeNext()
		}),
	)
	httpExec(chain)
	if !cmp.Equal(result, []string{"1", "2", "3"}) {
		t.Errorf("HandlerChain exec fail %v", result)
	}
	result = []string{}

	//-------------------------------------------------
	chain = newHandlerChain(
		f(func(c *Context) {
			result = append(result, "1")
			c.InvokeNext()
		}),
		f(func(c *Context) {
			result = append(result, "2")
		}),
		f(func(c *Context) {
			result = append(result, "3")
			c.InvokeNext()
		}),
	)
	httpExec(chain)
	if !cmp.Equal(result, []string{"1", "2"}) {
		t.Errorf("HandlerChain exec fail %v", result)
	}
	result = []string{}

	//-------------------------------------------------
	chain = newHandlerChain()
	chain.append((func(c *Context) {
		result = append(result, "1")
		c.InvokeNext()
	}))
	chain.append((func(c *Context) {
		result = append(result, "2")
		c.InvokeNext()
	}))
	chain.append((func(c *Context) {
		result = append(result, "3")
		c.InvokeNext()
	}))
	httpExec(chain)
	if !cmp.Equal(result, []string{"1", "2", "3"}) {
		t.Errorf("HandlerChain exec fail %v", result)
	}
	result = []string{}

	//-------------------------------------------------
	chain = newHandlerChain(func(c *Context) {
		result = append(result, "1")
		c.InvokeNext()
	})
	chain.append((func(c *Context) {
		result = append(result, "2")
		c.InvokeNext()
	}))
	chain.append((func(c *Context) {
		result = append(result, "3")
		c.InvokeNext()
	}))
	httpExec(chain)
	if !cmp.Equal(result, []string{"1", "2", "3"}) {
		t.Errorf("HandlerChain exec fail %v", result)
	}
	result = []string{}

	//-------------------------------------------------
	chain = newHandlerChain(func(c *Context) {
		result = append(result, "1")
	})
	chain.append((func(c *Context) {
		result = append(result, "2")
		c.InvokeNext()
	}))
	chain.append((func(c *Context) {
		result = append(result, "3")
		c.InvokeNext()
	}))
	httpExec(chain)
	if !cmp.Equal(result, []string{"1"}) {
		t.Errorf("HandlerChain exec fail %v", result)
	}
	result = []string{}

	//-------------------------------------------------
	chain = newHandlerChain(func(c *Context) {
		result = append(result, "1")
		c.InvokeNext()
	})
	chain.append((func(c *Context) {
		result = append(result, "2")
	}))
	chain.append((func(c *Context) {
		result = append(result, "3")
		c.InvokeNext()
	}))
	httpExec(chain)
	if !cmp.Equal(result, []string{"1", "2"}) {
		t.Errorf("HandlerChain exec fail %v", result)
	}
}
