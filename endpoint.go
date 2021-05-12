package kelly

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	// GET 方法
	GET = "GET"
	// HEAD 方法
	HEAD = "HEAD"
	// OPTIONS 方法
	OPTIONS = "OPTIONS"
	// POST 方法
	POST = "POST"
	// PUT 方法
	PUT = "PUT"
	// PATCH 方法
	PATCH = "PATCH"
	// DELETE 方法
	DELETE = "DELETE"
)

type endpoint struct {
	method   string
	path     string
	handlers []AnnotationHandlerFunc
}

func newEndpoint(method, path string, handlers ...AnnotationHandlerFunc) *endpoint {
	return &endpoint{
		method:   method,
		path:     path,
		handlers: handlers,
	}
}

func (endpoint endpoint) doPreRun(rootPath string, router *router, handlerList ...[]AnnotationHandlerFunc) {
	var f func(path string, handle httprouter.Handle) = nil
	switch endpoint.method {
	case GET:
		f = router.httpRouter().GET
	case HEAD:
		f = router.httpRouter().HEAD
	case OPTIONS:
		f = router.httpRouter().OPTIONS
	case POST:
		f = router.httpRouter().POST
	case PUT:
		f = router.httpRouter().PUT
	case PATCH:
		f = router.httpRouter().PATCH
	case DELETE:
		f = router.httpRouter().DELETE
	default:
		return
	}

	urlPath := rootPath + endpoint.path
	annotationContext := &AnnotationContext{
		Path:   urlPath,
		Method: endpoint.method,
		Router: router,
	}

	// 合并所有router的handler
	chain := newHandlerChain()
	for _, handlers := range handlerList {
		for _, handler := range handlers {
			chain.append(handler(annotationContext))
		}
	}
	//合并自己的handler
	for _, handler := range endpoint.handlers {
		chain.append(handler(annotationContext))
	}

	// 注入到httprouter
	f(urlPath, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		chain.ServeHTTP(w, r, params)
	})
}
