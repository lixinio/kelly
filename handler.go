package kelly

import "net/http"

// AnnotationHandlerFunc todo
type AnnotationHandlerFunc func(*AnnotationContext) HandlerFunc

// Handler http请求处理接口
type Handler interface {
	ServeHTTP(*Context)
}

// HandlerFunc http请求处理函数
type HandlerFunc func(*Context)

type handlerFuncWrap struct {
	hf HandlerFunc
}

func (hfw *handlerFuncWrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hfw.hf(newContext(w, r))
}

type handlerWrap struct {
	h Handler
}

func (hw *handlerWrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hw.h.ServeHTTP(newContext(w, r))
}

func wrapHttpHandlerFunc(f http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		f(c.ResponseWriter, c.Request())
	}
}
