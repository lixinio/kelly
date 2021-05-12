package kelly

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lixinio/kelly/binding"
)

const (
	contextNextHandle         = "__next_handler" // 下一个handler
	contextDataKeyPathVarible = "__path_varible" // path变量
)

var (
	gBinder binderAdapter = nil // 全局的binder
)

func init() {
	gBinder = binding.NewBinder()
}

// Context kelly在调用链传递的对象， 包装request/response
type Context struct {
	http.Hijacker
	http.Flusher
	http.ResponseWriter               // 备份http.ResponseWriter
	contextData                       // 支持绑定自定义数据
	r                   *http.Request // 备份http.Request
	response                          // 处理各种输出
	request                           // 读取各种请求参数
	binder                            // 绑定参数到对象
}

// Request 获得http.Request对象
func (c *Context) Request() *http.Request {
	return c.r
}

// SetRequest todo
// 有些中间件可能会重新生成一个新的request
// 比如 r = r.WithContext(context.WithValue(r.Context(), addedTagsKey{}, &tags))
// 需要重新设置request
func (c *Context) SetRequest(r *http.Request) *Context {
	c.r = r
	return c
}

// Context 获得 context.Context
func (c *Context) Context() context.Context {
	return c.r.Context()
}

// InvokeNext 触发调用链的下一个handler
func (c *Context) InvokeNext() {
	n := c.Get(contextNextHandle)
	if n != nil {
		rn := n.(*handlerChainEntry)
		if rn == nil {
			fmt.Println("fuck")
		}
		rn.ServeHTTP(c)
	}
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{
		ResponseWriter: w,
		r:              r,
		contextData:    newContextMapData(),
	}
	c.response = newResponse(c)
	c.request = newRequest(c, r)
	c.binder = newBinder(c, gBinder)
	return c
}

// AnnotationContext，用于记录每个请求的静态信息
type AnnotationContext struct {
	Router Router
	Method string
	Path   string
}
