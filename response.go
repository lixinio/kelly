package kelly

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"

	responseBackend "github.com/lixinio/kelly/response"
)

type response interface {
	// 返回紧凑的json
	WriteJSON(int, interface{})
	// 返回xml
	WriteXML(int, interface{})
	// 返回html
	WriteHTML(int, string)
	// 返回模板html
	WriteTemplateHTML(int, *template.Template, interface{})
	// 返回格式化的json
	WriteIndentedJSON(int, interface{})
	// 返回文本
	WriteString(int, string, ...interface{})
	// 返回二进制数据
	WriteData(int, string, []byte)
	// 返回紧凑的json，直接从二进制读数据
	WriteRawJSON(int, []byte)
	// 返回重定向
	Redirect(int, string)
	// 设置header
	SetHeader(string, string)
	// 设置cookie
	SetCookie(string, string, int, string, string, bool, bool)

	Abort(int, string)
	ResponseStatusOK()
	ResponseStatusBadRequest(error)
	ResponseStatusUnauthorized(error)
	ResponseStatusForbidden(error)
	ResponseStatusNotFound(error)
	ResponseStatusInternalServerError(error)
}

type responseImp struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	c *Context
}

func (r *responseImp) SetCookie(
	name string,
	value string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(r, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (r *responseImp) SetHeader(key, value string) {
	if len(value) == 0 {
		r.Header().Del(key)
	} else {
		r.Header().Set(key, value)
	}
}

func (r *responseImp) WriteJSON(code int, obj interface{}) {
	if err := responseBackend.WriteJSON(r, code, obj); err != nil {
		panic(fmt.Errorf("write json fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) WriteRawJSON(code int, content []byte) {
	if err := responseBackend.WriteRawJSON(r, code, content); err != nil {
		panic(fmt.Errorf("write raw json fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) WriteIndentedJSON(code int, obj interface{}) {
	if err := responseBackend.WriteIndentedJSON(r, code, obj); err != nil {
		panic(fmt.Errorf("write indented json fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) WriteHTML(code int, data string) {
	if err := responseBackend.WriteHTML(r, code, data); err != nil {
		panic(fmt.Errorf("write html fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) WriteTemplateHTML(code int, temp *template.Template, data interface{}) {
	if err := responseBackend.WriteTemplateHTML(r, code, temp, data); err != nil {
		panic(fmt.Errorf("write template fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) WriteXML(code int, obj interface{}) {
	if err := responseBackend.WriteXML(r, code, obj); err != nil {
		panic(fmt.Errorf("write xml fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) WriteString(code int, format string, values ...interface{}) {
	if err := responseBackend.WriteString(r, code, format, values); err != nil {
		panic(fmt.Errorf("write string fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) Redirect(code int, location string) {
	if err := responseBackend.Redirect(r, code, r.c.Request(), location); err != nil {
		panic(fmt.Errorf("redirect (%d|%s) fail, : %w(%s)", code, location, ErrWriteRespFail, err))
	}
}

func (r *responseImp) WriteData(code int, contentType string, data []byte) {
	if err := responseBackend.WriteData(r, code, contentType, data); err != nil {
		panic(fmt.Errorf("write data fail, : %w(%s)", ErrWriteRespFail, err))
	}
}

func (r *responseImp) Abort(code int, msg string) {
	if code == http.StatusNoContent {
		r.ResponseWriter.WriteHeader(code)
		return
	}

	if len(msg) == 0 {
		msg = http.StatusText(code)
	}
	r.WriteJSON(code, H{
		"code":    code,
		"message": msg,
	})
}

func (r *responseImp) ResponseStatusOK() {
	r.Abort(http.StatusOK, "")
}
func (r *responseImp) ResponseStatusBadRequest(err error) {
	if err != nil {
		r.Abort(http.StatusBadRequest, err.Error())
	} else {
		r.Abort(http.StatusBadRequest, "")
	}
}
func (r *responseImp) ResponseStatusUnauthorized(err error) {
	if err != nil {
		r.Abort(http.StatusUnauthorized, err.Error())
	} else {
		r.Abort(http.StatusUnauthorized, "")
	}
}
func (r *responseImp) ResponseStatusForbidden(err error) {
	if err != nil {
		r.Abort(http.StatusForbidden, err.Error())
	} else {
		r.Abort(http.StatusForbidden, "")
	}
}
func (r *responseImp) ResponseStatusNotFound(err error) {
	if err != nil {
		r.Abort(http.StatusNotFound, err.Error())
	} else {
		r.Abort(http.StatusNotFound, "")
	}
}
func (r *responseImp) ResponseStatusInternalServerError(err error) {
	if err != nil {
		r.Abort(http.StatusInternalServerError, err.Error())
	} else {
		r.Abort(http.StatusInternalServerError, "")
	}
}

func newResponse(c *Context) *responseImp {
	return &responseImp{
		ResponseWriter: c,
		Hijacker:       c,
		Flusher:        c,
		c:              c,
	}
}

// H 辅助类
// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
type H map[string]interface{}

// MarshalXML allows type H to be used with xml.Marshal
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}

// https://stackoverflow.com/questions/30928770/marshall-map-to-xml-in-go
// UnmarshalXML unmarshals the XML into a map of string to strings,
// creating a key in the map for each tag and setting it's value to the
// tags contents.
//
// The fact this function is on the pointer of Map is important, so that
// if m is nil it can be initialized, which is often the case if m is
// nested in another xml structurel. This is also why the first thing done
// on the first line is initialize it.
type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// UnmarshalXML imp
func (h *H) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// *m = Map{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*h)[e.XMLName.Local] = e.Value
	}
	return nil
}
