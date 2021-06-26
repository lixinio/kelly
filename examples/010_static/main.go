package main

import (
	"net/http"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/static"
)

func main() {
	router := kelly.New(nil)

	router.GET("/", func(c *kelly.Context) {
		c.Redirect(http.StatusFound, "/files")
	})

	router.GET("/files/*path", static.Static(&static.Config{
		Dir:           http.Dir("/home/king"),
		EnableListDir: true,
		Indexfiles:    []string{"index.html"},
		Handler404: func(c *kelly.Context) {
			c.WriteString(http.StatusNotFound, "找不到文件")
		},
	}))

	router.Run(":9999")
}
