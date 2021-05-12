package swagger

import (
	"strings"

	"github.com/lixinio/kelly"
)

const (
	defaultSwaggerBasePath    = "/api/v1"
	defaultSwaggerTitle       = "Swagger文档"
	defaultSwaggerDescription = "Swagger文档描述"
	defaultSwaggerApiVersion  = "0.0.1"
	defaultSwaggerDocDir      = "./"
	defaultSwaggerUI          = "https://petstore.swagger.io/"
	defaultSwaggerVersion     = "2.0"
)

// Config swagger配置
type Config struct {
	// "api前缀，例如/api/v1"，默认为空
	BasePath string

	// swagger文档标题
	Title string

	// swagger文档描述
	Description string

	// 接口版本
	ApiVersion string

	// swagger ui的地址
	SwaggerUiUrl string

	// swagger文档的地址，用于调试，release直接打包到二进制里面。默认为空
	DocFilePath string

	// 用于支持swagger ui认证头的参数
	Headers []SecurityDefinition

	DocLoader DocLoader
	// 是否调试模式
	Debug bool
}

func defaultConfig(config *Config) *Config {
	if config == nil {
		config = &Config{}
	}
	if len(config.BasePath) == 0 {
		config.BasePath = defaultSwaggerBasePath
	}
	if len(config.Title) == 0 {
		config.Title = defaultSwaggerTitle
	}
	if len(config.Description) == 0 {
		config.Title = defaultSwaggerDescription
	}
	if len(config.ApiVersion) == 0 {
		config.Title = defaultSwaggerApiVersion
	}
	if len(config.SwaggerUiUrl) == 0 {
		config.SwaggerUiUrl = defaultSwaggerUI
	}
	if len(config.DocFilePath) == 0 && config.DocLoader == nil {
		config.DocFilePath = defaultSwaggerDocDir
	}

	return config
}

func NewSwagger(r kelly.Router, config *Config) *Swagger {
	config = defaultConfig(config)
	swagger := newSwagger(config)
	r.Kelly().RegistePreRunHandler(swagger.PreRunHandler)
	return swagger
}

type DocLoader func(key string) ([]byte, error)

type Swagger struct {
	// 是否调试模式
	debugFlag bool
	// url前缀
	baseUrl    string
	cache      *cache
	pathEditor *pathEditor
	config     *Config
	specData   []byte
}

func newSwagger(config *Config) *Swagger {
	opt := &Swagger{
		debugFlag:  config.Debug,
		baseUrl:    config.BasePath,
		cache:      newCache(config.DocFilePath, config.DocLoader),
		pathEditor: newPathEditor(),
		config:     config,
	}

	return opt
}

func (s *Swagger) PreRunHandler(k kelly.Kelly) {
	if s.cache != nil {
		// 运行前一次性生成spec， 然后清理缓存
		s.specData = s.cache.buidAndClear(s.config)
		s.cache = nil
	}
}

// 去掉/api/v1之类的前缀
func (s *Swagger) realPath(r kelly.Router, path string) string {
	if len(s.baseUrl) > 0 && strings.HasPrefix(path, s.baseUrl) {
		path = strings.TrimPrefix(path, s.baseUrl)
	}

	return s.pathEditor.update(path)
}

func (s *Swagger) SwaggerFile(swaggerEntry string) kelly.AnnotationHandlerFunc {
	return func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		path := s.realPath(ac.Router, ac.Path)
		s.cache.getEntry(swaggerEntry, path, ac.Method)
		return nil
	}
}
