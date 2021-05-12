package main

import (
	"net/http"

	"github.com/lixinio/kelly"
)

func initRouter3(r kelly.Router) {
	api := r.Group("/v3")
	api.GET("/", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.WriteIndentedJSON(http.StatusOK, kelly.H{
				"message": "this is v3",
				"code":    "0",
			})
		}
	})
}

func initRouter1(r kelly.Router) {
	api := r.Group("/api/v1")
	initRouter3(api)
	api.GET("/", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.WriteIndentedJSON(http.StatusOK, kelly.H{
				"message": "this is v1",
				"code":    "0",
			})
		}
	})
}

func initRouter2(r kelly.Router) {
	api := r.Group("/api/v2")
	api.GET("/", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.WriteIndentedJSON(http.StatusOK, kelly.H{
				"message": "this is v2",
				"code":    "0",
			})
		}
	})
}

func main() {
	router := kelly.New(nil)
	initRouter1(router)
	initRouter2(router)
	router.GET("/", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.WriteIndentedJSON(http.StatusOK, kelly.H{
				"message": "this is main",
				"code":    "0",
			})
		}
	})
	router.Run(":9999")
}
