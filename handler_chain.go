package kelly

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handlerChainEntry struct {
	current HandlerFunc
	next    *handlerChainEntry
}

func (handlerChainEntry *handlerChainEntry) ServeHTTP(c *Context) {
	if handlerChainEntry.next != nil {
		c.Set(contextNextHandle, handlerChainEntry.next)
	} else {
		c.Set(contextNextHandle, nil)
	}
	handlerChainEntry.current(c)
}

// HandlerChain 调用链， 依次调用
type HandlerChain struct {
	handlers []HandlerFunc
	head     *handlerChainEntry
	tail     *handlerChainEntry
}

// append 链尾添加回调
func (handlerChain *HandlerChain) append(handler HandlerFunc) *HandlerChain {
	if handler == nil {
		return handlerChain
	}

	handlerChain.handlers = append(handlerChain.handlers, handler)
	entry := &handlerChainEntry{
		current: handler,
		next:    nil,
	}
	if handlerChain.head == nil {
		handlerChain.head = entry
		handlerChain.tail = entry
		return handlerChain
	}
	handlerChain.tail.next = entry
	handlerChain.tail = entry
	return handlerChain
}

func (handlerChain *HandlerChain) appends(handlers []HandlerFunc) *HandlerChain {
	for _, handler := range handlers {
		handlerChain.append(handler)
	}
	return handlerChain
}

func (handlerChain *HandlerChain) ServeHTTP(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if handlerChain.head != nil {
		// 将Httprouter的接口转换成kelly Context
		c := newContext(w, r)
		// 存储path变量
		c.Set(contextDataKeyPathVarible, params)
		handlerChain.head.ServeHTTP(c)
	}
}

func newHandlerChain(handlers ...HandlerFunc) *HandlerChain {
	chain := &HandlerChain{
		handlers: handlers,
		head:     nil,
		tail:     nil,
	}

	for _, handler := range handlers {
		chain.append(handler)
	}
	return chain
}
