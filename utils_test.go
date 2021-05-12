package kelly

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func wrapKellyHandler(handler HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(newContext(w, r))
	}
}

func getFramwork(
	handler HandlerFunc,
	path string,
	cookies map[string]string,
	headers map[string]string,
	contentType string,
) *http.Response {
	mux := http.NewServeMux()
	mux.HandleFunc("/", wrapKellyHandler(handler))
	r, _ := http.NewRequest(http.MethodGet, path, nil)

	for k, v := range cookies {
		r.AddCookie(&http.Cookie{
			Name:     k,
			Value:    url.QueryEscape(v),
			MaxAge:   0,
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: false,
		})
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	if contentType != "" {
		r.Header.Add("Content-Type", contentType)
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Result()
}

func postFormFramwork(handler HandlerFunc, path string, values [][]string) *http.Response {
	mux := http.NewServeMux()
	mux.HandleFunc("/", wrapKellyHandler(handler))

	var rf http.Request
	rf.ParseForm()
	for _, v := range values {
		rf.Form.Add(v[0], v[1])
	}
	reader := strings.NewReader(rf.Form.Encode())
	r, _ := http.NewRequest(http.MethodPost, path, reader)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Result()
}

func kellyFramwork(pattern, path string, groups []string, handlers ...AnnotationHandlerFunc) *http.Response {
	return kellyFramworkImp(pattern, "", path, groups, map[string]string{}, nil, handlers...)
}

func kellyFormFramwork(pattern, path string, groups []string, form map[string]string, handlers ...AnnotationHandlerFunc) *http.Response {
	return kellyFramworkImp(pattern, "post", path, groups, form, nil, handlers...)
}

func kellyJSONFramwork(pattern, path string, groups []string, body interface{}, handlers ...AnnotationHandlerFunc) *http.Response {
	return kellyFramworkImp(pattern, "patch", path, groups, map[string]string{}, body, handlers...)
}

func kellyFramworkImp(
	pattern, method, path string, groups []string,
	form map[string]string,
	body interface{},
	handlers ...AnnotationHandlerFunc,
) *http.Response {
	k := New(nil)
	if len(groups) == 0 {
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
	} else {
		var g Router = k
		for _, group := range groups {
			g = g.Group(group)
		}

		if method == "" || strings.ToLower(method) == "get" {
			method = "get"
			g.GET(pattern, handlers...)
		} else if strings.ToLower(method) == "post" {
			method = "post"
			g.POST(pattern, handlers...)
		} else {
			method = "patch"
			g.PATCH(pattern, handlers...)
		}
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

	return k.RunTest(r)
}

func readBody(resp *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return strings.Trim(buf.String(), " \n")
}

func readRawBody(resp *http.Response) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.Bytes()
}

func readHeader(resp *http.Response, name string) string {
	return resp.Header.Get(name)
}

func readContentType(resp *http.Response) string {
	return filterFlags(resp.Header.Get("Content-Type"))
}

func readCookie(resp *http.Response, name string) string {
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func newImage() []byte {
	width := 200
	height := 100
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

func checkError(t *testing.T, e error) {
	if r := recover(); r == nil {
		t.Errorf("checkError did not panic")
	} else {
		switch x := r.(type) {
		case error:
			if !errors.Is(x, e) {
				t.Errorf("invalid error type %v", x)
			}
		default:
			t.Errorf("invalid error type %v", x)
		}
	}
}

func TestFilterFlags(t *testing.T) {
	data := [][]string{
		{
			"application/json; charset=utf-8",
			"application/json",
		},
		{
			"application/json charset=utf-8",
			"application/json",
		},
		{
			"text/javascript; charset=utf-8",
			"text/javascript",
		},
		{
			"text/javascript;charset=utf-8",
			"text/javascript",
		},
	}
	for _, item := range data {
		if filterFlags(item[0]) != item[1] {
			t.Errorf("parse [%s]->[%s] fail", item[0], item[1])
			return
		}
	}
}
