// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package response

import (
	"html/template"
	"io"
	"net/http"
)

const htmlContentType = "text/html; charset=utf-8"

func WriteHTML(w http.ResponseWriter, code int, data string) error {
	writeContentType(w, htmlContentType)
	w.WriteHeader(code)

	if _, err := io.WriteString(w, data); err != nil {
		return err
	}
	return nil
}

func WriteTemplateHTML(w http.ResponseWriter, code int, temp *template.Template, data interface{}) error {
	writeContentType(w, htmlContentType)
	w.WriteHeader(code)

	return temp.Execute(w, data)
}
