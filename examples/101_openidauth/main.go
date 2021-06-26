package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/openid"
)

const (
	// issuer           = "http://127.0.0.1:9999"
	aud              = "kelly.openidauth"
	userID           = "123"
	defaultPort uint = 9999
)

type User struct {
	ID   string
	Name string
}

func initAuth(router kelly.Router, handler *Handler, iss string) {
	r := router.Group("/auth")

	r.GET("/login", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		data := `<form action="#" method="post">
<p>{{ .First }}: <input type="text" name="username" value="any" /></p>
<p>{{ .Last }}: <input type="text" name="password" value="456" /></p>
<input type="submit" value="提交" />
</form>`

		// 通过闭包预先编译好
		t := template.Must(template.New("t1").Parse(data))
		return func(c *kelly.Context) {
			c.WriteTemplateHTML(http.StatusOK, t, map[string]string{
				"First": "用户名",
				"Last":  "密码",
			})
		}
	})

	r.GET("/login/fail", func(c *kelly.Context) {
		c.Abort(401, "认证失败")
	})

	r.POST("/login", func(c *kelly.Context) {
		username, err := c.GetFormVarible("username")
		if err != nil {
			c.Redirect(http.StatusFound, r.Path()+"/login/fail")
			return
		}
		password, err := c.GetFormVarible("password")
		if err != nil {
			c.Redirect(http.StatusFound, r.Path()+"/login/fail")
			return
		}

		if len(username) > 0 && password == "456" {
			// 构建jwt claims数据
			claims := makeClaims(aud, userID, username, iss, 1000)
			token, err := handler.issueToken(claims)
			if err != nil {
				c.Abort(500, err.Error())
			}
			c.WriteString(http.StatusOK, token)
		} else {
			c.Redirect(http.StatusFound, r.Path()+"/login/fail")
			return
		}
	})
}

func main() {
	port := defaultPort
	args := os.Args
	if len(args) >= 2 {
		newPort, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil || uint(newPort) == defaultPort {
			panic(err)
		}
		port = uint(newPort)
	}
	iss := fmt.Sprintf("http://127.0.0.1:%d", port)

	router := kelly.New(nil)

	// openssl genrsa -out a.key 2048
	// openssl rsa -in a.key -pubout > a.pub
	handler := initOpenIDServer(router, "/home/king/code/qjw/gosample/build/a.key", iss)
	initAuth(router, handler, iss)

	if port == defaultPort {
		// 运行openid 认证服务器
		router.Run(fmt.Sprintf(":%d", port))
		return
	}

	//
	openidAuth, err := openid.OpenIDAuth(&openid.OpenIDAuthConfig{
		TokenGetter: func(c *kelly.Context) (string, error) {
			// 从cookie获取
			token, err := c.GetQueryVarible("token")
			if err != nil {
				return "", err
			} else if token == "" {
				return "", errors.New("empty token")
			}
			return token, nil
		},
		Authorizator: func(claims *openid.MapClaims) (interface{}, error) {
			// 写死， 用作测试
			if (*claims)["sub"].(string) == userID {
				return &User{
					ID:   userID,
					Name: (*claims)["name"].(string),
				}, nil
			} else {
				return nil, errors.New("invalid user id")
			}
		},
		Audience: aud,
		Issuer:   fmt.Sprintf("http://127.0.0.1:%d", defaultPort),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	router.GET("/", openidAuth, func(c *kelly.Context) {
		user, _ := openid.CurrentUser(c).(*User)
		c.WriteIndentedJSON(http.StatusOK, kelly.H{
			"code": "0",
			"data": user.Name,
		})
	})

	router.Run(fmt.Sprintf(":%d", port))
}
