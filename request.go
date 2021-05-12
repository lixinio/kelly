package kelly

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type request interface {
	// 根据key获取cookie值
	GetCookie(string) (string, error)
	// 根据key获取cookie值，若不存在，则返回默认值
	GetDefaultCookie(string, string) string
	// 根据key获取cookie值，若不存在，则panic
	MustGetCookie(string) string

	// 根据key获取header值
	GetHeader(string) (string, error)
	// 根据key获取header值，若不存在，则返回默认值
	GetDefaultHeader(string, string) string
	// 根据key获取header值，若不存在，则panic
	MustGetHeader(string) string
	// Content-Type
	ContentType() string

	// 根据key获取PATH变量值
	GetPathVarible(string) (string, error)
	// 根据key获取PATH变量值，若不存在，则panic
	MustGetPathVarible(string) string

	// 根据key获取QUERY变量值，可能包含多个（http://127.0.0.1:9090/path/abc?abc=bbb&abc=aaa）
	GetMultiQueryVarible(string) ([]string, error)
	// 根据key获取QUERY变量值，仅返回第一个
	GetQueryVarible(string) (string, error)
	// 根据key获取QUERY变量值，仅返回第一个,若不存在，则返回默认值
	GetDefaultQueryVarible(string, string) string
	// 根据key获取QUERY变量值，仅返回第一个,若不存在，则panic
	MustGetQueryVarible(string) string

	// 根据key获取FORM变量值，可能get可能包含多个
	GetMultiFormVarible(string) ([]string, error)
	// 根据key获取FORM变量值，仅返回第一个
	GetFormVarible(string) (string, error)
	// 根据key获取FORM变量值，仅返回第一个,若不存在，则返回默认值
	GetDefaultFormVarible(string, string) string
	// 根据key获取FORM变量值，仅返回第一个,若不存在，则panic
	MustGetFormVarible(string) string

	// @ref http.Request.ParseMultipartForm
	ParseMultipartForm() error
	// 获取（上传的）文件信息
	GetFileVarible(string) (multipart.File, *multipart.FileHeader, error)
	MustGetFileVarible(string) (multipart.File, *multipart.FileHeader)
}

type requestImp struct {
	*Context
	*http.Request
}

func (r requestImp) GetCookie(name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("get cookie fail(%s),: %w(%s)", name, ErrNoCookie, err)
	}
	val, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return "", fmt.Errorf("cookie content(%s) invalid fail,: %w(%s)", cookie.Value, ErrInvalidCookie, err)
	}
	return val, nil
}

func (r requestImp) GetDefaultCookie(name, defaultValue string) string {
	if cookie, err := r.GetCookie(name); err == nil {
		return cookie
	}
	return defaultValue
}

func (r requestImp) MustGetCookie(name string) string {
	cookie, err := r.GetCookie(name)
	if err == nil {
		return cookie
	}
	panic(err)
}

func (r requestImp) GetHeader(name string) (string, error) {
	header := r.Header.Get(name)
	if len(header) > 0 {
		return header, nil
	}

	return "", fmt.Errorf("header(%s) not exist: %w", name, ErrNoHeader)
}

func (r requestImp) GetDefaultHeader(name, defaultValue string) string {
	if cookie, err := r.GetHeader(name); err == nil {
		return cookie
	}
	return defaultValue
}

func (r requestImp) MustGetHeader(name string) string {
	cookie, err := r.GetHeader(name)
	if err == nil {
		return cookie
	}
	panic(err)
}

func (r requestImp) ContentType() string {
	if ct, err := r.GetHeader("Content-Type"); err == nil {
		return filterFlags(ct)
	}
	return ""
}

func (r requestImp) GetPathVarible(name string) (string, error) {
	params := r.Context.MustGet(contextDataKeyPathVarible).(httprouter.Params)
	val := params.ByName(name)
	if len(val) > 0 {
		return val, nil
	}
	return val, fmt.Errorf("router path varible(%s) not exist: %w", name, ErrNoPathVarible)
}

func (r requestImp) MustGetPathVarible(name string) string {
	val, err := r.GetPathVarible(name)
	if err == nil {
		return val
	}
	panic(err)
}

// -----------------------------------------------------------------

func (r requestImp) GetMultiQueryVarible(name string) ([]string, error) {
	if values, ok := r.getQueryArray(name); ok {
		return values, nil
	}
	return []string{}, fmt.Errorf("request query(%s) not exist: %w", name, ErrNoQueryVarible)
}

func (r requestImp) GetQueryVarible(name string) (string, error) {
	values, err := r.GetMultiQueryVarible(name)
	if err == nil {
		return values[0], nil
	}
	return "", err
}

func (r requestImp) GetDefaultQueryVarible(name, defaultValue string) string {
	if val, err := r.GetQueryVarible(name); err == nil {
		return val
	}
	return defaultValue
}
func (r requestImp) MustGetQueryVarible(name string) string {
	val, err := r.GetQueryVarible(name)
	if err == nil {
		return val
	}
	panic(err)
}

// -----------------------------------------------------------------

func (r requestImp) GetMultiFormVarible(name string) ([]string, error) {
	if values, ok := r.getPostFormArray(name); ok {
		return values, nil
	}
	return []string{}, fmt.Errorf("request form(%s) not exist: %w", name, ErrNoFormVarible)
}

func (r requestImp) GetFormVarible(name string) (string, error) {
	values, err := r.GetMultiFormVarible(name)
	if err == nil {
		return values[0], nil
	}
	return "", err
}

func (r requestImp) GetDefaultFormVarible(name, defaultValue string) string {
	if val, err := r.GetFormVarible(name); err == nil {
		return val
	}
	return defaultValue
}
func (r requestImp) MustGetFormVarible(name string) string {
	val, err := r.GetFormVarible(name)
	if err == nil {
		return val
	}
	panic(err)
}

func (r requestImp) ParseMultipartForm() error {
	return r.Request.ParseMultipartForm(32 << 20)
}

func (r requestImp) GetFileVarible(name string) (multipart.File, *multipart.FileHeader, error) {
	f, fh, err := r.FormFile(name)
	if err == nil {
		return f, fh, nil
	}
	return f, fh, fmt.Errorf("get file varible(%s) fail, : %w(%s)", name, ErrNoFileVarible, err)
}

func (r requestImp) MustGetFileVarible(name string) (multipart.File, *multipart.FileHeader) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		panic(fmt.Errorf("get file varible(%s) fail, : %w(%s)", name, ErrNoFileVarible, err))
	}
	return file, handler
}

// -----------------------------------------------------------------

// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (r requestImp) getQueryArray(key string) ([]string, bool) {
	req := r.Request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (r requestImp) getPostFormArray(key string) ([]string, bool) {
	req := r.Request
	// request自己会缓存
	req.ParseForm()
	req.ParseMultipartForm(32 << 20) // 32 MB
	if values := req.PostForm[key]; len(values) > 0 {
		return values, true
	}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values, true
		}
	}
	return []string{}, false
}

func newRequest(c *Context, r *http.Request) request {
	return &requestImp{
		Context: c,
		Request: r,
	}
}
