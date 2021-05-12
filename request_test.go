package kelly

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHandleQuery(t *testing.T) {
	f := func(c *Context) {
		v, err := c.GetQueryVarible("a")
		if err != nil || v != "b" {
			t.Errorf("query varible is %v|%w", v, err)
			return
		}
		v, err = c.GetQueryVarible("b")
		if err == nil {
			t.Errorf("query varible is %v|%w", v, err)
			return
		}

		v = c.GetDefaultQueryVarible("d", "e")
		if v != "e" {
			t.Errorf("query varible is %v", v)
			return
		}
		v = c.MustGetQueryVarible("c")
		if v != "d" {
			t.Errorf("query varible is %v", v)
			return
		}

		v2, err := c.GetMultiQueryVarible("a")
		if err != nil || !cmp.Equal(
			v2,
			[]string{"f", "b"},
			cmpopts.SortSlices(func(i, j string) bool { return i < j }),
		) {
			t.Errorf("GetMultiQueryVarible varible is %v|%w", v2, err)
			return
		}

		v2, err = c.GetMultiQueryVarible("c")
		if err != nil || !cmp.Equal(v2, []string{"d"}) {
			t.Errorf("GetMultiQueryVarible varible is %v|%w", v2, err)
			return
		}

		// 抛异常
		defer checkError(t, ErrNoQueryVarible)
		c.MustGetQueryVarible("b")

		c.WriteString(http.StatusOK, "abc")
	}

	getFramwork(f, "/?a=b&c=d&a=f", map[string]string{}, map[string]string{}, "")
}

func TestHandleForm(t *testing.T) {
	f := func(c *Context) {
		v, err := c.GetFormVarible("a")
		if err != nil || v != "b" {
			t.Errorf("form varible is %v|%w", v, err)
			return
		}
		v, err = c.GetQueryVarible("b")
		if err == nil {
			t.Errorf("form varible is %v|%w", v, err)
			return
		}

		v = c.GetDefaultFormVarible("d", "e")
		if v != "e" {
			t.Errorf("form varible is %v", v)
			return
		}
		v = c.MustGetFormVarible("c")
		if v != "d" {
			t.Errorf("form varible is %v", v)
			return
		}

		v2, err := c.GetMultiFormVarible("a")
		if err != nil || !cmp.Equal(
			v2,
			[]string{"f", "b"},
			cmpopts.SortSlices(func(i, j string) bool { return i < j }),
		) {
			t.Errorf("GetMultiFormVarible varible is %v|%w", v2, err)
			return
		}

		v2, err = c.GetMultiFormVarible("c")
		if err != nil || !cmp.Equal(v2, []string{"d"}) {
			t.Errorf("GetMultiFormVarible varible is %v|%w", v2, err)
			return
		}

		// 抛异常
		defer checkError(t, ErrNoFormVarible)
		c.MustGetFormVarible("b")

		c.WriteString(http.StatusOK, "abc")
	}

	postFormFramwork(f, "/", [][]string{
		{"a", "b"},
		{"c", "d"},
		{"a", "f"},
		{"x", "发送到 发的撒发"},
	})
}

func TestHandleCookie(t *testing.T) {
	f := func(c *Context) {
		v, err := c.GetCookie("a")
		if err != nil || v != "b" {
			t.Errorf("cookie varible is %v|%w", v, err)
			return
		}
		v, err = c.GetCookie("b")
		if err == nil {
			t.Errorf("cookie varible is %v|%w", v, err)
			return
		}

		v = c.GetDefaultCookie("d", "e")
		if v != "e" {
			t.Errorf("cookie varible is %v", v)
			return
		}
		v = c.MustGetCookie("c")
		if v != "d" {
			t.Errorf("cookie varible is %v", v)
			return
		}

		// 抛异常
		defer checkError(t, ErrNoCookie)
		c.MustGetCookie("b")

		c.WriteString(http.StatusOK, "abc")
	}

	getFramwork(f, "/",
		map[string]string{
			"a": "b",
			"c": "d",
			"x": "发送到 发的撒发",
		},
		map[string]string{},
		"",
	)
}

func TestHandleHeader(t *testing.T) {
	f := func(c *Context) {
		v, err := c.GetHeader("a")
		if err != nil || v != "b" {
			t.Errorf("header is %v|%w", v, err)
			return
		}
		v, err = c.GetCookie("b")
		if err == nil || !errors.Is(err, ErrNoCookie) {
			t.Errorf("header is %v|%w", v, err)
			return
		}

		v = c.GetDefaultHeader("d", "e")
		if v != "e" {
			t.Errorf("header is %v", v)
			return
		}
		v = c.MustGetHeader("c")
		if v != "d" {
			t.Errorf("header is %v", v)
			return
		}

		// 抛异常
		defer checkError(t, ErrNoHeader)
		c.MustGetHeader("b")

		c.WriteString(http.StatusOK, "abc")
	}

	getFramwork(f, "/",
		map[string]string{},
		map[string]string{
			"a": "b",
			"c": "d",
		},
		"",
	)
}

func TestHandleContentType(t *testing.T) {
	f := func(contentType string) func(*Context) {
		return func(c *Context) {
			if c.ContentType() != contentType {
				t.Errorf("contentType is %v|%v", c.ContentType(), contentType)
				return
			}
			c.WriteString(http.StatusOK, "abc")
		}
	}

	getFramwork(f("application/json"), "/",
		map[string]string{},
		map[string]string{},
		"application/json; charset=utf-8",
	)

	getFramwork(f("text/javascript"), "/",
		map[string]string{},
		map[string]string{},
		"text/javascript; ",
	)

	getFramwork(f("text/javascript"), "/",
		map[string]string{},
		map[string]string{},
		"text/javascript ",
	)

	getFramwork(f("text/javascript"), "/",
		map[string]string{},
		map[string]string{},
		"text/javascript",
	)
}

func TestKellyPathVarible(t *testing.T) {
	resp := kellyFramwork("/a/:path", "/a/xxx", []string{}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("path")
			if err != nil || v != "xxx" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}

			v, err = c.GetPathVarible("path2")
			if err == nil || !errors.Is(err, ErrNoPathVarible) {
				t.Errorf("GetPathVarible fail")
				return
			}

			v = c.MustGetPathVarible("path")
			if v != "xxx" {
				t.Errorf("MustGetPathVarible fail, %v", v)
				return
			}

			defer checkError(t, ErrNoPathVarible)
			c.MustGetPathVarible("path2")

			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusOK error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/a/:path", "/a/xxx/yyy", []string{}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusNotFound error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/a/:path/yyy", "/a/xxx/yyy", []string{}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("path")
			if err != nil || v != "xxx" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusOK error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/:path/yyy", "/a/xxx/yyy", []string{"/a"}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("path")
			if err != nil || v != "xxx" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusOK error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/yyy", "/a/xxx/yyy", []string{"/a/:path"}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("path")
			if err != nil || v != "xxx" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusOK error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/a/*path", "/a/xxx", []string{}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("path")
			if err != nil || v != "/xxx" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusOK error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/a/*path", "/a/xxx/yyy", []string{}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("path")
			if err != nil || v != "/xxx/yyy" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusOK error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/*path", "/a/xxx/yyy", []string{"/a"}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("path")
			if err != nil || v != "/xxx/yyy" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusOK error, %v", resp.StatusCode)
		return
	}

}

func TestFileVarible(t *testing.T) {
	key := "file"
	filename := "a.png"
	data := newImage()

	f := func(c *Context) {
		f, header, err := c.GetFileVarible(key)
		if f != nil {
			defer f.Close()
		}

		if err != nil {
			t.Errorf("GetFileVarible fail |%w", err)
			return
		}
		if header.Filename != filename {
			t.Errorf("GetFileVarible filename fail %s|%s", header.Filename, filename)
			return
		}

		result := bytes.NewBuffer(nil)
		if _, err := io.Copy(result, f); err != nil {
			t.Errorf("GetFileVarible Copy fail |%w", err)
			return
		}

		if !cmp.Equal(result.Bytes(), data) {
			t.Errorf("GetFileVarible Equal fail")
			return
		}

		f, header = c.MustGetFileVarible(key)

		defer checkError(t, ErrNoFileVarible)
		c.MustGetFileVarible("file2")

		c.WriteString(http.StatusOK, "abc")
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile(key, filename)
		if err != nil {
			t.Error(err)
			return
		}
		if cnt, err := part.Write(data); err != nil {
			t.Error(err)
			return
		} else if cnt < 1 {
			t.Errorf("part.Write fail %d", cnt)
			return
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", wrapKellyHandler(f))

	r, _ := http.NewRequest(http.MethodPost, "/", pr)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("upload file fail")
		return
	}

}
