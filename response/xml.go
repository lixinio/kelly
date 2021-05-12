// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package response

import (
	"encoding/xml"
	"net/http"
)

var xmlContentType = "application/xml; charset=utf-8"

func WriteXML(w http.ResponseWriter, code int, data interface{}) error {
	writeContentType(w, xmlContentType)
	w.WriteHeader(code)

	return xml.NewEncoder(w).Encode(data)
}
