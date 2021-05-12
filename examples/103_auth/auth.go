package main

import (
	"net/http"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/sessions"
)

func InitApiV1(r kelly.Router, mng *sessions.LoginManager) {

	api := r.Group("/api/v1")
	api.GET("/",
		mng.LoginRequired(),
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				// 获取登录用户
				user := mng.GetCurrentUser(c).(*User)
				c.WriteJSON(http.StatusOK, kelly.H{
					"id":      user.id,
					"message": user.name,
				})
			}
		},
	)

	api.GET("/login",
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				// 是否已经登录
				if mng.IsAuthenticated(c) {
					c.Redirect(http.StatusFound, "/api/v1/")
					return
				}

				id := c.GetDefaultQueryVarible("name", "default_user")

				// 登录授权
				mng.Login(c, &User{
					id:   id,
					name: "name",
				})
				c.Redirect(http.StatusFound, "/api/v1/")
			}
		},
	)

	api2 := api.Group("/", mng.LoginRequired())
	api2.GET("/logout",
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				// 注销登录
				mng.Logout(c)
				c.WriteJSON(http.StatusFound, "/logout")
			}
		},
	)
}
