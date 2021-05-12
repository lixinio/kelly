package swagger

import "errors"

var (
	// ErrInvalidHttpMethod 不支持的http方法
	ErrInvalidHttpMethod error = errors.New("invalid http method")
	// ErrInvalidDocDirectory 错误的文档目录
	ErrInvalidDocDirectory = errors.New("invalid doc directory")
	// ErrGenerateSwaggerSpecFail 生成swagger spec数据失败
	ErrGenerateSwaggerSpecFail = errors.New("build spec data fail")
	// ErrOpenDocFileFail 打开文档失败
	ErrOpenDocFileFail = errors.New("open document file fail")
	// ErrParseDocFileFail 解析文档失败
	ErrParseDocFileFail = errors.New("parse document file fail")
	// ErrInvalidSwaggerDecorator swagger装饰器非法
	ErrInvalidSwaggerDecorator = errors.New("invalid swagger decorator")
	// ErrSwaggerDecoratorNotExist swagger装饰器不存在
	ErrSwaggerDecoratorNotExist = errors.New("swagger decorator not exists")
)
