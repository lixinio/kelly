package kelly

import "errors"

var (
	// ErrNoContextData 没有Kelly.Context对象
	ErrNoContextData error = errors.New("kelly context data not exist")
	// ErrNoCookie 没有cookie
	ErrNoCookie = errors.New("request cookie not exist")
	// ErrInvalidCookie 错误的cookie
	ErrInvalidCookie = errors.New("request cookie content is invalid")
	// ErrNoHeader 没有请求头
	ErrNoHeader = errors.New("request header not exist")
	// ErrNoQueryVarible 没有请求参数
	ErrNoQueryVarible = errors.New("request query varible not exist")
	// ErrNoFormVarible 没有form表单参数
	ErrNoFormVarible = errors.New("request form varible not exist")
	// ErrNoPathVarible 没有path变量
	ErrNoPathVarible = errors.New("router path varible not exist")
	// ErrNoFileVarible 没有file变量（文件上传）
	ErrNoFileVarible = errors.New("router file varible not exist")
	// ErrInvalidRouterPath 错误的路由路径
	ErrInvalidRouterPath = errors.New("router path is invalid")
	// ErrInvalidHandler 错误的处理句柄
	ErrInvalidHandler = errors.New("handler is invalid")
	// ErrWriteRespFail 写响应失败
	ErrWriteRespFail = errors.New("write response fail")
	// ErrBindFail bind请求参数（到对象）失败
	ErrBindFail = errors.New("bind varible fail")
)
