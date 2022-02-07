package main

import (
	"net/http"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/validator/obj"
)

type BindObj struct {
	A string   `json:"aaa,omitempty" validate:"required,max=32,min=6" error:"长度6/32"`
	B string   `json:"bbb,omitempty"`
	C int      `json:"ccc,omitempty"`
	D []int    `json:"ddd,omitempty" validate:"dive,min=1,max=10"`
	E []string `json:"eee,omitempty" disable_split:""`
}

func bindErrorHandler(c *kelly.Context, err error) {
	c.WriteString(http.StatusOK, "参数错误: %s", err.Error())
}

func main() {
	router := kelly.New(nil)
	validator := obj.NewValidator()

	// http://127.0.0.1:9999/query?aaa=123456&bbb=2&ccc=3
	router.GET("/query", func(c *kelly.Context) {
		var obj BindObj
		if err := c.Bind(&obj); err == nil {
			c.WriteJSON(http.StatusOK, obj)
		} else {
			c.WriteString(http.StatusOK, "参数错误: %s", err.Error())
		}
	})

	// http://127.0.0.1:9999/query2?aaa=123456&bbb=2&ccc=3&ddd=4,5,6&eee=asd,fgh,ijk
	router.GET("/query2",
		kelly.BindMiddleware(
			func() interface{} { return &BindObj{} },
			validator,
			bindErrorHandler,
		),
		func(c *kelly.Context) {
			bObj := c.GetBindParameter().(*BindObj)
			c.WriteJSON(http.StatusOK, bObj)
		},
	)

	// http://127.0.0.1:9999/path/abcdef/b/1
	router.GET("/path/:aaa/:bbb/:ccc",
		kelly.BindPathMiddleware(
			func() interface{} { return &BindObj{} },
			validator,
			bindErrorHandler,
		),
		func(c *kelly.Context) {
			bObj := c.GetBindPathParameter().(*BindObj)
			c.WriteJSON(http.StatusOK, bObj)
		},
	)

	router.Run(":9999")
}
