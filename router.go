package kelly

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router 路由
type Router interface {
	GET(string, ...interface{}) Router
	HEAD(string, ...interface{}) Router
	OPTIONS(string, ...interface{}) Router
	POST(string, ...interface{}) Router
	PUT(string, ...interface{}) Router
	PATCH(string, ...interface{}) Router
	DELETE(string, ...interface{}) Router

	// 新建子路由
	Group(string, ...interface{}) Router

	// 动态插入中间件
	Use(...interface{}) Router

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

func validateHandlers(handlers ...interface{}) []AnnotationHandlerFunc {
	var result []AnnotationHandlerFunc
	for _, item := range handlers {
		if item == nil {
			panic(fmt.Errorf("handler can NOT be empty, : %w", ErrInvalidHandler))
		}

		if f, ok := item.(func(*Context)); ok {
			result = append(result, func(*AnnotationContext) HandlerFunc {
				return f
			})
		} else if f, ok := item.(HandlerFunc); ok {
			result = append(result, func(*AnnotationContext) HandlerFunc {
				return f
			})
		} else if f, ok := item.(func(*AnnotationContext) HandlerFunc); ok {
			result = append(result, f)
		} else if f, ok := item.(AnnotationHandlerFunc); ok {
			result = append(result, f)
		} else {
			panic(fmt.Errorf("handler must be AnnotationHandlerFunc|HandlerFunc , : %w", ErrInvalidHandler))
		}
	}
	return result
}

func (rt *router) validateParam(path string, handlers ...interface{}) []AnnotationHandlerFunc {
	if len(handlers) < 1 {
		panic(fmt.Errorf("must have one handler at least, : %w", ErrInvalidHandler))
	}

	rt.validatePath(path)
	return validateHandlers(handlers...)
}

func (rt *router) methodImp(
	method string,
	path string,
	handlers ...interface{},
) Router {
	annotationHandlers := rt.validateParam(path, handlers...)
	rt.endpoints = append(rt.endpoints, newEndpoint(method, path, annotationHandlers...))
	return rt
}

func (rt *router) GET(path string, handlers ...interface{}) Router {
	return rt.methodImp(GET, path, handlers...)
}

func (rt *router) HEAD(path string, handlers ...interface{}) Router {
	return rt.methodImp(HEAD, path, handlers...)
}

func (rt *router) OPTIONS(path string, handlers ...interface{}) Router {
	return rt.methodImp(OPTIONS, path, handlers...)
}

func (rt *router) POST(path string, handlers ...interface{}) Router {
	return rt.methodImp(POST, path, handlers...)
}

func (rt *router) PUT(path string, handlers ...interface{}) Router {
	return rt.methodImp(PUT, path, handlers...)
}

func (rt *router) PATCH(path string, handlers ...interface{}) Router {
	return rt.methodImp(PATCH, path, handlers...)
}

func (rt *router) DELETE(path string, handlers ...interface{}) Router {
	return rt.methodImp(DELETE, path, handlers...)
}

func (rt *router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rt.rt.ServeHTTP(rw, r)
}

func (rt *router) Group(path string, handlers ...interface{}) Router {
	if path == "/" {
		path = ""
	}

	newRt := newRouterImp(rt.rt, rt.k, rt, path, rt.absolutePath+path, handlers...)
	rt.groups = append(rt.groups, newRt)
	return newRt
}

func (rt *router) Use(handlers ...interface{}) Router {
	annotationHandlers := validateHandlers(handlers...)
	for _, v := range annotationHandlers {
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

func newRouterImp(
	rt *httprouter.Router,
	k Kelly,
	parent *router,
	path, absolutePath string,
	handlers ...interface{},
) *router {
	annotationHandlers := validateHandlers(handlers...)
	return &router{
		rt:           rt,
		path:         path,
		absolutePath: absolutePath,
		middlewares:  annotationHandlers,
		parent:       parent,
		k:            k,
	}
}
