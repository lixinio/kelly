package kelly

import (
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	resp := kellyFramwork("/a/:path", "/a/xxx/yyy", []string{}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusNotFound error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/a/*path", "/a/xxx/yyy", []string{}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusNotFound error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/a/:path", "/1/2/3/a/xxx", []string{"/1", "/2", "/3"}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusNotFound error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/a/:path", "/1/2/3/a/xxx", []string{"/:1", "/:2", "/:3"}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusNotFound error, %v", resp.StatusCode)
		return
	}

	resp = kellyFramwork("/:a/*path", "/v1/v2/v3/va/xxx/yyy/zzz", []string{"/:1", "/:2", "/:3"}, func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			v, err := c.GetPathVarible("1")
			if err != nil || v != "v1" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			v, err = c.GetPathVarible("2")
			if err != nil || v != "v2" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			v, err = c.GetPathVarible("3")
			if err != nil || v != "v3" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			v, err = c.GetPathVarible("a")
			if err != nil || v != "va" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}
			v, err = c.GetPathVarible("path")
			if err != nil || v != "/xxx/yyy/zzz" {
				t.Errorf("GetPathVarible fail, %v|%w", v, err)
				return
			}

			c.WriteString(http.StatusOK, "ok")
		}
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusNotFound error, %v", resp.StatusCode)
		return
	}
}
