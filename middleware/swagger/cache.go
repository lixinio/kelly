package swagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lixinio/kelly"
	"gopkg.in/yaml.v2"
)

type cache struct {
	fileData    map[string]*SwaggerDocFile   // 文档对象缓存
	swaggerData map[string]*SwaggerPathEntry // // 文件和里面内容的缓存
	docDir      string                       // 文档的目录
	docLoader   DocLoader
}

func newCache(docDir string, docLoader DocLoader) *cache {
	if len(docDir) > 0 {
		if path, err := filepath.Abs(docDir); err != nil {
			panic(fmt.Errorf("invalid doc dir : %w, err %s", ErrInvalidDocDirectory, err.Error()))
		} else {
			if strings.HasSuffix(path, "/") {
				docDir = path
			} else {
				docDir = path + "/"
			}
		}
	}

	return &cache{
		docDir:      docDir,
		docLoader:   docLoader,
		fileData:    make(map[string]*SwaggerDocFile),
		swaggerData: make(map[string]*SwaggerPathEntry),
	}
}

func (fc *cache) build(config *Config) []byte {
	headersDef := make(map[string]SecurityDefinition)
	if len(config.Headers) > 0 {
		for _, v := range config.Headers {
			key := v.Type
			v.In = "header"
			v.Type = "apiKey"
			headersDef[key] = v
		}
	}

	jsonBytes, err := json.MarshalIndent(kelly.H{
		"basePath": config.BasePath,
		"swagger":  defaultSwaggerVersion,
		"info": struct {
			Description string `json:"description"`
			Title       string `json:"title"`
			Version     string `json:"version"`
		}{
			Description: config.Description,
			Title:       config.Title,
			Version:     config.ApiVersion,
		},
		"definition":          struct{}{},
		"paths":               fc.swaggerData,
		"securityDefinitions": headersDef,
	}, "", "    ")
	if err != nil {
		panic(fmt.Errorf("build spec fail : %w, error %s", ErrGenerateSwaggerSpecFail, err.Error()))
	}
	return jsonBytes
}

func (fc *cache) buidAndClear(config *Config) []byte {
	// 生成可以直接使用的二进制数据
	data := fc.build(config)
	// 清除缓存， 释放内存
	fc.fileData = make(map[string]*SwaggerDocFile)
	fc.swaggerData = make(map[string]*SwaggerPathEntry)
	return data
}

func (fc *cache) getFile(filepath string) *SwaggerDocFile {
	// 是否有缓存
	if docFile, ok := fc.fileData[filepath]; ok {
		return docFile
	}

	var yamlFile []byte
	var err error
	if len(fc.docDir) > 0 {
		// 如果指定了目录， 就直接读取
		filepath = fc.docDir + filepath

		// 加载文件
		yamlFile, err = ioutil.ReadFile(filepath)
		if err != nil {
			panic(fmt.Errorf("open document(%s) fail : %w, error %s", filepath, ErrOpenDocFileFail, err))
		}
	} else {
		// 对于一些下载文件逻辑比较特殊， 支持自行实现文件读取， 比如从网络读取文档信息
		yamlFile, err = fc.docLoader(filepath)
		if err != nil {
			panic(fmt.Errorf("open document(%s) fail : %w, error %s", filepath, ErrOpenDocFileFail, err))
		}
	}

	var docFile SwaggerDocFile
	err = yaml.Unmarshal(yamlFile, &docFile)
	if err != nil {
		panic(fmt.Errorf("parse document(%s) fail : %w, error %s", filepath, ErrParseDocFileFail, err))
	}

	// 写入缓存
	fc.fileData[filepath] = &docFile
	return &docFile
}

func (fc *cache) getEntry(swaggerEntry, path, method string) *SwaggerApiEntry {
	// 解析文件路径和内部路径
	filepath, entry, err := parseFileNode(swaggerEntry)
	if err != nil {
		panic(fmt.Errorf("invalid decorator(%s) : %w, error %s", swaggerEntry, ErrInvalidSwaggerDecorator, err.Error()))
	}

	// 找到文件内容
	docFile := fc.getFile(filepath)
	// 找到文件具体的段落
	if data, ok := (*docFile)[entry]; ok {
		var sentry *SwaggerPathEntry
		if v, ok := fc.swaggerData[path]; ok {
			sentry = v
			sentry.update(method, data)
		} else {
			sentry = newSwaggerEntry(method, data)
		}
		fc.swaggerData[path] = sentry
		return data
	} else {
		panic(fmt.Errorf("swagger decorator entry(%s) not exist : %w", entry, ErrSwaggerDecoratorNotExist))
	}
}
