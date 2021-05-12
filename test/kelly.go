package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lixinio/kelly"
)

func KellyFramwork(
	pattern, path string,
	headers map[string]string,
	handlers ...kelly.AnnotationHandlerFunc,
) *http.Response {
	return kellyFramworkImp(pattern, "", path, headers, map[string]string{}, nil, handlers...)
}

func KellyFormFramwork(
	pattern, path string,
	headers map[string]string, form map[string]string,
	handlers ...kelly.AnnotationHandlerFunc,
) *http.Response {
	return kellyFramworkImp(pattern, "post", path, headers, form, nil, handlers...)
}

func KellyJsonFramwork(
	pattern, path string, body interface{},
	headers map[string]string,
	handlers ...kelly.AnnotationHandlerFunc,
) *http.Response {
	return kellyFramworkImp(pattern, "patch", path, headers, map[string]string{}, body, handlers...)
}

func kellyFramworkImp(
	pattern, method, path string,
	headers map[string]string,
	form map[string]string,
	body interface{},
	handlers ...kelly.AnnotationHandlerFunc,
) *http.Response {
	k := kelly.New(nil)
	if method == "" || strings.ToLower(method) == "get" {
		method = "get"
		k.GET(pattern, handlers...)
	} else if strings.ToLower(method) == "post" {
		method = "post"
		k.POST(pattern, handlers...)
	} else {
		method = "patch"
		k.PATCH(pattern, handlers...)
	}

	var r *http.Request = nil
	if method == "get" {
		r, _ = http.NewRequest(http.MethodGet, path, nil)
	} else if method == "post" {
		var rf http.Request
		rf.ParseForm()
		for k, v := range form {
			rf.Form.Add(k, v)
		}
		reader := strings.NewReader(rf.Form.Encode())
		r, _ = http.NewRequest(http.MethodPost, path, reader)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		if body == nil {
			panic("invalid body")
		}
		jsonData, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		r, _ = http.NewRequest(http.MethodPatch, path, bytes.NewBuffer(jsonData))
		r.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		r.Header.Set(k, v)
	}

	return k.RunTest(r)
}
