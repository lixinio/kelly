package kelly

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type BindObj struct {
	A string `json:"aaa"`
	B string `json:"bbb"`
	C int    `json:"ccc"`
	D bool   `json:"ddd"`
}

type validatorBindObj struct {
}

func (v validatorBindObj) Validate(obj interface{}) error {
	bingObj, ok := obj.(*BindObj)
	if !ok {
		return fmt.Errorf("invalid object type")
	}
	if cmp.Equal(bingObj, &BindObj{
		A: "b",
		B: "d",
		C: 123,
		D: true,
	}) {
		return nil
	}
	return fmt.Errorf("invalid object value")
}

func bindErrorHandler(c *Context, err error) {
	c.WriteString(http.StatusBadRequest, "参数错误: %s", err.Error())
}

func TestBinding(t *testing.T) {
	m := BindMiddleware(
		func() interface{} { return &BindObj{} },
		nil,
		bindErrorHandler,
	)
	f := func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := c.GetBindParameter().(*BindObj)
			if !cmp.Equal(obj, &BindObj{
				A: "b",
				B: "d",
				C: 123,
				D: true,
			}) {
				t.Errorf("GetBindParameter err %v", obj)
			}
			c.ResponseStatusOK()
		}
	}

	resp := kellyFramwork("/", "/?aaa=b&bbb=d&ccc=123&ddd=true", []string{}, m, f)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GetBindParameter error %d", resp.StatusCode)
	}

	groups := []string{
		"/a", "/b", "/c", "/d",
	}
	middlewares := []AnnotationHandlerFunc{
		func(ac *AnnotationContext) HandlerFunc { return func(*Context) {} },
		func(ac *AnnotationContext) HandlerFunc { return func(*Context) {} },
		func(ac *AnnotationContext) HandlerFunc { return func(*Context) {} },
		m,
		f,
	}
	middlewares2 := []AnnotationHandlerFunc{
		m,
		func(ac *AnnotationContext) HandlerFunc { return func(*Context) {} },
		func(ac *AnnotationContext) HandlerFunc { return func(*Context) {} },
		func(ac *AnnotationContext) HandlerFunc { return func(*Context) {} },
	}
	for i := range groups {
		path := strings.Join(groups[:i+1], "")

		for i2 := 0; i2 <= len(middlewares)-2; i2++ {
			resp := kellyFramwork(
				"/", fmt.Sprintf("%s/?aaa=b&bbb=d&ccc=123&ddd=true", path), groups[:i+1], middlewares[i2:]...,
			)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("GetBindParameter error %d", resp.StatusCode)
			}
		}

		for i2 := 1; i2 <= len(middlewares2); i2++ {
			resp := kellyFramwork(
				"/", fmt.Sprintf("%s/?aaa=b&bbb=d&ccc=123&ddd=true", path), groups[:i+1],
				append(middlewares2[:i2], f)...,
			)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("GetBindParameter error %d", resp.StatusCode)
			}
		}
	}
}

func TestQueryBinding(t *testing.T) {
	m := BindMiddleware(
		func() interface{} { return &BindObj{} },
		nil,
		bindErrorHandler,
	)
	f := func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := c.GetBindParameter().(*BindObj)
			if !cmp.Equal(obj, &BindObj{
				A: "b",
				B: "d",
				C: 123,
				D: true,
			}) {
				t.Errorf("GetBindParameter err %v", obj)
			}

			c.ResponseStatusOK()
		}
	}

	resp := kellyFramwork("/", "/?aaa=b&bbb=d&ccc=123&ddd=true", []string{}, m, f)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GetBindParameter error %d", resp.StatusCode)
	}

	resp = kellyFramwork("/", "/?aaa=b&bbb=d&ccc=efg&ddd=true", []string{}, m, f)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("GetBindParameter error %d", resp.StatusCode)
	}
}

func TestPathBinding(t *testing.T) {
	m := BindPathMiddleware(
		func() interface{} { return &BindObj{} },
		nil,
		bindErrorHandler,
	)
	f := func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := c.GetBindPathParameter().(*BindObj)
			if !cmp.Equal(obj, &BindObj{
				A: "b",
				B: "d",
				C: 123,
				D: true,
			}) {
				t.Errorf("GetBindPathParameter err %v", obj)
			}

			c.ResponseStatusOK()
		}
	}

	resp := kellyFramwork("/:aaa/:bbb/:ccc/:ddd", "/b/d/123/true", []string{}, m, f)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GetBindPathParameter error %d", resp.StatusCode)
	}

	resp = kellyFramwork("/:aaa/:bbb/:ccc/:ddd", "/b/d/abc/true", []string{}, m, f)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("GetBindPathParameter error %d", resp.StatusCode)
	}
}

func TestFormBinding(t *testing.T) {
	// form 不区分int
	m := BindFormMiddleware(
		func() interface{} { return &BindObj{} },
		nil,
		bindErrorHandler,
	)
	f := func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := c.GetBindFormParameter().(*BindObj)
			if !cmp.Equal(
				obj,
				&BindObj{
					A: "b",
					B: "d",
					C: 123,
					D: true,
				},
				cmpopts.IgnoreFields(BindObj{}, "C", "D"),
			) {
				t.Errorf("TestFormBinding err %v", obj)
			}

			c.ResponseStatusOK()
		}
	}

	resp := kellyFormFramwork("/", "/", []string{}, map[string]string{
		"aaa": "b",
		"bbb": "d",
	}, m, f)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestFormBinding error %d", resp.StatusCode)
	}
}

func TestJsonBinding(t *testing.T) {
	body := &BindObj{
		A: "b",
		B: "d",
		C: 123,
		D: true,
	}
	m := BindJSONMiddleware(
		func() interface{} { return &BindObj{} },
		nil,
		bindErrorHandler,
	)
	f := func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := c.GetBindJSONParameter().(*BindObj)
			if !cmp.Equal(obj, body) {
				t.Errorf("TestJsonBinding err %v", obj)
			}

			c.ResponseStatusOK()
		}
	}

	resp := kellyJSONFramwork("/", "/", []string{}, body, m, f)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestJsonBinding error %d", resp.StatusCode)
	}

	resp = kellyJSONFramwork("/", "/", []string{}, &H{
		"aaa": "b",
		"bbb": "d",
		"ccc": 123,
		"ddd": true,
	}, m, f)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestJsonBinding error %d", resp.StatusCode)
	}

	resp = kellyJSONFramwork("/", "/", []string{}, &H{
		"aaa": "b",
		"bbb": "d",
		"ccc": "e",
		"ddd": true,
	}, m, f)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("TestJsonBinding error %d", resp.StatusCode)
	}
}

func TestBindingValidator(t *testing.T) {
	m := BindMiddleware(
		func() interface{} { return &BindObj{} },
		&validatorBindObj{},
		bindErrorHandler,
	)
	f := func(ac *AnnotationContext) HandlerFunc {
		return func(c *Context) {
			obj := c.GetBindParameter().(*BindObj)
			if !cmp.Equal(obj, &BindObj{
				A: "b",
				B: "d",
				C: 123,
				D: true,
			}) {
				t.Errorf("TestBindingValidator err %v", obj)
			}

			c.ResponseStatusOK()
		}
	}

	resp := kellyFramwork("/", "/?aaa=b&bbb=d&ccc=123&ddd=true", []string{}, m, f)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestBindingValidator error %d", resp.StatusCode)
	}

	resp = kellyFramwork("/", "/?aaa=b&bbb=d&ccc=456&ddd=true", []string{}, m, f)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("TestBindingValidator error %d", resp.StatusCode)
	}

	resp = kellyFramwork("/", "/?aaa=b&bbb=e&ccc=123&ddd=true", []string{}, m, f)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("TestBindingValidator error %d", resp.StatusCode)
	}
}
