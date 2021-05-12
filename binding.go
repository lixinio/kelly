package kelly

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lixinio/kelly/validator"
	"github.com/mitchellh/mapstructure"
)

// binderAdapter 绑定输入的适配接口
type binderAdapter interface {
	Bind(*http.Request, interface{}) error     // 绑定一个对象，根据Content-type自动判断类型
	BindJSON(*http.Request, interface{}) error // 绑定json，从body取数据
	BindXML(*http.Request, interface{}) error  // 绑定xml，从body取数据
	BindForm(*http.Request, interface{}) error // 绑定form，从body/query取数据
}

type binder interface {
	Bind(interface{}) error     // 绑定一个对象，根据Content-type自动判断类型
	BindJSON(interface{}) error // 绑定json，从body取数据
	BindXML(interface{}) error  // 绑定xml，从body取数据
	BindForm(interface{}) error // 绑定form，从body/query取数据
	BindPath(interface{}) error // 绑定path变量

	GetBindParameter() interface{}
	GetBindJSONParameter() interface{}
	GetBindXMLParameter() interface{}
	GetBindFormParameter() interface{}
	GetBindPathParameter() interface{}
}

// BindErrorHandle bind失败的错误处理
type BindErrorHandle func(*Context, error)

func handleBindErr(c *Context, err error) {
	c.WriteJSON(http.StatusBadRequest, H{
		"code":  http.StatusUnprocessableEntity,
		"error": err.Error(),
	})
}

const (
	contextBindKey     = "_binder_key"
	contextBindJSONKey = "_binder_json_key"
	contextBindXMLKey  = "_binder_xml_key"
	contextBindFormKey = "_binder_form_key"
	contextBindPathKey = "_binder_path_key"
)

type binderImp struct {
	c *Context      // http 请求上下文
	b binderAdapter // binder实现
}

func newBinder(c *Context, b binderAdapter) binder {
	return &binderImp{
		c: c,
		b: b,
	}
}

func wrapBindError(message string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("bind error(%s), : %w(%s)", message, ErrBindFail, err)
}

func (b *binderImp) GetBindParameter() interface{} {
	return b.c.MustGet(contextBindKey)
}

func (b *binderImp) GetBindJSONParameter() interface{} {
	return b.c.MustGet(contextBindJSONKey)
}

func (b *binderImp) GetBindXMLParameter() interface{} {
	return b.c.MustGet(contextBindXMLKey)
}

func (b *binderImp) GetBindFormParameter() interface{} {
	return b.c.MustGet(contextBindFormKey)
}

func (b *binderImp) GetBindPathParameter() interface{} {
	return b.c.MustGet(contextBindPathKey)
}

func (b *binderImp) Bind(obj interface{}) error {
	return wrapBindError("bind", b.b.Bind(b.c.Request(), obj))
}

func (b *binderImp) BindJSON(obj interface{}) error {
	return wrapBindError("bind json", b.b.BindJSON(b.c.Request(), obj))
}

func (b *binderImp) BindXML(obj interface{}) error {
	return wrapBindError("bind xml", b.b.BindXML(b.c.Request(), obj))
}

func (b *binderImp) BindForm(obj interface{}) error {
	return wrapBindError("bind form", b.b.BindForm(b.c.Request(), obj))
}

func (b *binderImp) BindPath(obj interface{}) error {
	params := b.c.Get(contextDataKeyPathVarible).(httprouter.Params)
	myData := map[string]string{}
	for _, param := range params {
		myData[param.Key] = param.Value
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           obj,
		TagName:          "json",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return wrapBindError("bind path: mapstructure new decoder", err)
	}

	return wrapBindError("bind path: mapstructure decode", decoder.Decode(myData))
}

// BindMiddleware 绑定query参数中间件
func BindMiddleware(
	objG func() interface{},
	validator validator.Validator,
	errHandler BindErrorHandle,
) AnnotationHandlerFunc {
	if errHandler == nil {
		errHandler = handleBindErr
	}
	return func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := objG()
			err := c.Bind(obj)
			if err == nil {
				if validator != nil {
					err = validator.Validate(obj)
					if err == nil {
						c.Set(contextBindKey, obj)
						c.InvokeNext()
						return
					}
				} else {
					c.Set(contextBindKey, obj)
					c.InvokeNext()
					return
				}
			}

			errHandler(c, err)
		}
	}
}

// BindJSONMiddleware 绑定json参数中间件
func BindJSONMiddleware(
	objG func() interface{},
	validator validator.Validator,
	errHandler BindErrorHandle,
) AnnotationHandlerFunc {
	if errHandler == nil {
		errHandler = handleBindErr
	}

	return func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := objG()
			err := c.BindJSON(obj)
			if err == nil {
				if validator != nil {
					err = validator.Validate(obj)
					if err == nil {
						c.Set(contextBindJSONKey, obj)
						c.InvokeNext()
						return
					}
				} else {
					c.Set(contextBindJSONKey, obj)
					c.InvokeNext()
					return
				}
			}

			errHandler(c, err)
		}
	}
}

// BindXMLMiddleware 绑定xml参数中间件
func BindXMLMiddleware(
	objG func() interface{},
	validator validator.Validator,
	errHandler BindErrorHandle,
) AnnotationHandlerFunc {
	if errHandler == nil {
		errHandler = handleBindErr
	}

	return func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := objG()
			err := c.BindXML(obj)
			if err == nil {
				if validator != nil {
					err = validator.Validate(obj)
					if err == nil {
						c.Set(contextBindXMLKey, obj)
						c.InvokeNext()
						return
					}
				} else {
					c.Set(contextBindXMLKey, obj)
					c.InvokeNext()
					return
				}
			}

			errHandler(c, err)
		}
	}
}

// BindFormMiddleware 绑定form参数中间件
func BindFormMiddleware(
	objG func() interface{},
	validator validator.Validator,
	errHandler BindErrorHandle,
) AnnotationHandlerFunc {
	if errHandler == nil {
		errHandler = handleBindErr
	}

	return func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := objG()
			err := c.BindForm(obj)
			if err == nil {
				if validator != nil {
					err = validator.Validate(obj)
					if err == nil {
						c.Set(contextBindFormKey, obj)
						c.InvokeNext()
						return
					}
				} else {
					c.Set(contextBindFormKey, obj)
					c.InvokeNext()
					return
				}
			}

			errHandler(c, err)
		}
	}
}

// BindPathMiddleware 绑定path参数中间件
func BindPathMiddleware(
	objG func() interface{},
	validator validator.Validator,
	errHandler BindErrorHandle,
) AnnotationHandlerFunc {
	if errHandler == nil {
		errHandler = handleBindErr
	}

	return func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := objG()
			err := c.BindPath(obj)
			if err == nil {
				if validator != nil {
					err = validator.Validate(obj)
					if err == nil {
						c.Set(contextBindPathKey, obj)
						c.InvokeNext()
						return
					}
				} else {
					c.Set(contextBindPathKey, obj)
					c.InvokeNext()
					return
				}
			}

			errHandler(c, err)
		}
	}
}
