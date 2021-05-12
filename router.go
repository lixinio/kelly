package kelly

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router 路由
type Router interface {
	GET(string, ...AnnotationHandlerFunc) Router
	HEAD(string, ...AnnotationHandlerFunc) Router
	OPTIONS(string, ...AnnotationHandlerFunc) Router
	POST(string, ...AnnotationHandlerFunc) Router
	PUT(string, ...AnnotationHandlerFunc) Router
	PATCH(string, ...AnnotationHandlerFunc) Router
	DELETE(string, ...AnnotationHandlerFunc) Router

	// 新建子路由
	Group(string, ...AnnotationHandlerFunc) Router

	// 动态插入中间件
	Use(...AnnotationHandlerFunc) Router

	// 返回当前Router的绝对路径
	Path() string
	Kelly() Kelly
}

type router struct {
	rt           *httprouter.Router
	path         string                  // 当前rouer路径
	absolutePath string                  // 绝对路径
	middlewares  []AnnotationHandlerFunc // 中间件
	groups       []*router               // 所有的子Group
	parent       *router                 // 父Group
	endpoints    []*endpoint
	k            Kelly
}

func (rt router) httpRouter() *httprouter.Router {
	return rt.rt
}

func (rt *router) Path() string {
	return rt.absolutePath
}

func (rt *router) Kelly() Kelly {
	return rt.k
}

func (rt *router) validatePath(path string) {
	if len(path) < 1 {
		panic(fmt.Errorf("invalid path (%s), : %w", path, ErrInvalidRouterPath))
	}
	if path == "/" {
		return
	}
	if path[0] != '/' || path[len(path)-1] == '/' {
		panic(fmt.Errorf("invalid path (%s),must beginwith (NOT endwith) /, : %w", path, ErrInvalidRouterPath))
	}
}

func (rt *router) validateParam(path string, handlers ...AnnotationHandlerFunc) {
	if len(handlers) < 1 {
		panic(fmt.Errorf("must have one handle, : %w", ErrInvalidHandler))
	}
	rt.validatePath(path)
}

func (rt *router) methodImp(
	method string,
	path string,
	handlers ...AnnotationHandlerFunc,
) Router {

	rt.validateParam(path, handlers...)
	rt.endpoints = append(rt.endpoints, newEndpoint(method, path, handlers...))
	return rt
}

func (rt *router) GET(path string, handlers ...AnnotationHandlerFunc) Router {
	return rt.methodImp(GET, path, handlers...)
}

func (rt *router) HEAD(path string, handlers ...AnnotationHandlerFunc) Router {
	return rt.methodImp(HEAD, path, handlers...)
}

func (rt *router) OPTIONS(path string, handlers ...AnnotationHandlerFunc) Router {
	return rt.methodImp(OPTIONS, path, handlers...)
}

func (rt *router) POST(path string, handlers ...AnnotationHandlerFunc) Router {
	return rt.methodImp(POST, path, handlers...)
}

func (rt *router) PUT(path string, handlers ...AnnotationHandlerFunc) Router {
	return rt.methodImp(PUT, path, handlers...)
}

func (rt *router) PATCH(path string, handlers ...AnnotationHandlerFunc) Router {
	return rt.methodImp(PATCH, path, handlers...)
}

func (rt *router) DELETE(path string, handlers ...AnnotationHandlerFunc) Router {
	return rt.methodImp(DELETE, path, handlers...)
}

func (rt *router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rt.rt.ServeHTTP(rw, r)
}

func (rt *router) Group(path string, handlers ...AnnotationHandlerFunc) Router {
	rt.validatePath(path)
	if path == "/" {
		path = ""
	}
	newRt := &router{
		rt:           rt.rt,
		path:         path,
		absolutePath: rt.absolutePath + path,
		middlewares:  handlers,
		parent:       rt,
		k:            rt.k,
	}

	rt.groups = append(rt.groups, newRt)
	return newRt
}

func (rt *router) Use(handlers ...AnnotationHandlerFunc) Router {
	for _, v := range handlers {
		rt.middlewares = append(rt.middlewares, v)
	}
	return rt
}

// doPreRun 运行前的一些准备工作
func (rt *router) doPreRun(handlerList ...[]AnnotationHandlerFunc) {
	handlerList = append(handlerList, rt.middlewares)
	// 注入每一层的handler到endpoint
	for _, e := range rt.endpoints {
		e.doPreRun(rt.absolutePath, rt, handlerList...)
	}
	for _, subRouter := range rt.groups {
		subRouter.doPreRun(handlerList...)
	}
}

func newRouterImp(hr *httprouter.Router, handlers ...AnnotationHandlerFunc) *router {
	rt := &router{
		rt:           hr,
		path:         "",
		absolutePath: "",
		middlewares:  handlers,
	}

	return rt
}
