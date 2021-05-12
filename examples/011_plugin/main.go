package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/kelly"
)

func middleware(title string) kelly.AnnotationHandlerFunc {
	return func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(*kelly.Context) {
			fmt.Println("this is middleware ", title)
		}
	}
}

func initRouter3(r kelly.Router) {
	api := r.Group("/v3", middleware("router3_1"), middleware("router3_2"))
	api.Use(middleware("router3_3"), middleware("router3_4"))

	api.GET("/",
		middleware("view/api/v3_1"),
		middleware("view/api/v3_2"),
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				c.WriteIndentedJSON(http.StatusOK, kelly.H{
					"message": "this is v3",
					"code":    "0",
				})
			}
		},
	)

	api.Use(middleware("router3_5"))
}

func initRouter1(r kelly.Router) {
	api := r.Group("/api/v1", middleware("router1_1"), middleware("router1_2"))
	api.Use(middleware("router1_3"), middleware("router1_4"))

	initRouter3(api)
	api.GET("/",
		middleware("view/api/v1_1"),
		middleware("view/api/v1_2"),
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				c.WriteIndentedJSON(http.StatusOK, kelly.H{
					"message": "this is v1",
					"code":    "0",
				})
			}
		},
	)

	api.Use(middleware("router1_5"))
}

func initRouter2(r kelly.Router) {
	api := r.Group("/api/v2", middleware("router2_1"), middleware("router2_2"))
	api.Use(middleware("router2_3"), middleware("router2_4"))

	api.GET("/",
		middleware("view/api/v2_1"),
		middleware("view/api/v2_2"),
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				c.WriteIndentedJSON(http.StatusOK, kelly.H{
					"message": "this is v2",
					"code":    "0",
				})
			}
		},
	)

	api.Use(middleware("router2_5"))
}

func initRouterFilter(r kelly.Router) {
	r.GET("/user",
		middleware("view/user_1"),
		middleware("view/user_2"),
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				// 从query参数（id）获取
				uid := c.GetDefaultQueryVarible("id", "")
				if uid == "" {
					fmt.Println("param error, break")
					c.WriteIndentedJSON(http.StatusOK, kelly.H{
						"message": "this is user",
						"code":    "404",
					})
				} else {
					fmt.Println("get param ", uid)
					c.InvokeNext()
				}
			}
		},
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				c.WriteIndentedJSON(http.StatusOK, kelly.H{
					"message": "this is user",
					"code":    "0",
				})
			}
		},
	)
}

func main() {
	router := kelly.New(nil, middleware("root"))
	router.Use(middleware("main1"), middleware("main2"))

	initRouter1(router)
	initRouter2(router)
	initRouterFilter(router)
	router.GET("/", middleware("view/"), func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.WriteIndentedJSON(http.StatusOK, kelly.H{
				"message": "this is main",
				"code":    "0",
			})
		}
	})
	router.Use(middleware("main3"))
	router.Run(":9999")
}
