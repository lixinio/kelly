package swagger

import (
	"fmt"
	"strings"
)

type SwaggerDocFile map[string]*SwaggerApiEntry

type SwaggerApiEntry struct {
	Description string                 `json:"description,omitempty" yaml:"description"`
	Summary     string                 `json:"summary" yaml:"summary"`
	Tags        []string               `json:"tags" yaml:"tags" binding:"required,dive,required"`
	Parameters  []*SwaggerApiParameter `json:"parameters,omitempty" yaml:"parameters,omitempty" binding:"dive"`
	Produces    []string               `json:"produces,omitempty" yaml:"produces"`
	Responses   map[int]*JsonSchemaObj `json:"responses" yaml:"responses" binding:"required"`
	OperationId string                 `json:"operationId,omitempty"`
}

type SwaggerApiParameter struct {
	Description string         `json:"description,omitempty" yaml:"description"`
	In          string         `json:"in" yaml:"in" binding:"eq=query|eq=path|eq=formData|eq=body|eq=header"`
	Name        string         `json:"name" yaml:"name" binding:"required,max=100,min=1"`
	Required    bool           `json:"required" yaml:"required"`
	Example     interface{}    `json:"example,omitempty" yaml:"example,omitempty"`
	Type        string         `json:"type" yaml:"type" binding:"eq=string|eq=integer|eq=number|eq=boolean|eq=array|eq=object|eq=file"`
	Schema      *JsonSchemaObj `json:"schema,omitempty" yaml:"schema"`
}

type JsonSchemaObj struct {
	Description string                    `json:"description,omitempty" yaml:"description"`
	Type        string                    `json:"type,omitempty" yaml:"type" binding:"required,min=1"`
	Items       *JsonSchemaObj            `json:"items,omitempty" yaml:"items"`
	Properties  map[string]*JsonSchemaObj `json:"properties,omitempty" yaml:"properties"`
	Required    []string                  `json:"required,omitempty" yaml:"required"`
	Example     interface{}               `json:"example,omitempty" yaml:"example,omitempty"`
	Schema      *JsonSchemaObj            `json:"schema,omitempty" yaml:"schema"`
}

type SwaggerPathEntry struct {
	Post   *SwaggerApiEntry `json:"post,omitempty"`
	Get    *SwaggerApiEntry `json:"get,omitempty"`
	Put    *SwaggerApiEntry `json:"put,omitempty"`
	Delete *SwaggerApiEntry `json:"delete,omitempty"`
	Patch  *SwaggerApiEntry `json:"patch,omitempty"`
}

func newSwaggerEntry(method string, entry *SwaggerApiEntry) *SwaggerPathEntry {
	swaggerEntry := &SwaggerPathEntry{}
	swaggerEntry.update(method, entry)
	return swaggerEntry
}

func (swaggerEntry *SwaggerPathEntry) update(method string, entry *SwaggerApiEntry) *SwaggerPathEntry {
	if strings.ToLower(method) == "get" {
		swaggerEntry.Get = entry
	} else if strings.ToLower(method) == "post" {
		swaggerEntry.Post = entry
	} else if strings.ToLower(method) == "put" {
		swaggerEntry.Put = entry
	} else if strings.ToLower(method) == "delete" {
		swaggerEntry.Delete = entry
	} else if strings.ToLower(method) == "patch" {
		swaggerEntry.Patch = entry
	} else {
		panic(fmt.Errorf("invalid method %s, error : %w", method, ErrInvalidHttpMethod))
	}
	return swaggerEntry
}

type SecurityDefinition struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	In          string `json:"in"`
	Name        string `json:"name"`
}
