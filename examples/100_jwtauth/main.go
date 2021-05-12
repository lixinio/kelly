package main

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/jwt"
)

const secretKey string = "123456789"
const cookieKey string = "JwtToken"
const jwtAud string = "kelly.jwtauth"
const userID string = "123"

type User struct {
	ID   string
	Name string
}

func initAuth(router kelly.Router) {
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

	r.GET("/login/fail", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.Abort(401, "认证失败")
		}
	})

	r.POST("/login", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
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
				claims := jwt.NewMapClaims(
					jwtAud,
					"kelly.jwtauth",
					userID, // sub
					3600,   // 一小时
				)
				claims.Update("username", username)

				// 签发token并写入cookie
				token, _ := jwt.GenerateHS256Token(claims, secretKey)
				c.SetCookie(cookieKey, token, 0, "", "", false, false)
				c.WriteIndentedJSON(http.StatusOK, kelly.H{
					"code": "0",
				})
			} else {
				c.Redirect(http.StatusFound, r.Path()+"/login/fail")
				return
			}
		}
	})
}

func main() {
	router := kelly.New(nil)
	initAuth(router)

	jwtAuth := jwt.JwtAuth(&jwt.JwtAuthConfig{
		TokenGetter: func(c *kelly.Context) (string, error) {
			// 从cookie获取
			token, err := c.GetCookie(cookieKey)
			if err != nil {
				return "", err
			} else if token == "" {
				return "", errors.New("empty token")
			}
			return token, nil
		},
		Authorizator: func(claims jwt.Claims) (interface{}, error) {
			claimsObj, ok := claims.(*jwt.MapClaims)
			if !ok {
				return nil, errors.New("invalid claims object")
			}
			// 写死， 用作测试
			if claimsObj.Get("sub").(string) == userID {
				return &User{
					ID:   userID,
					Name: claimsObj.Get("username").(string),
				}, nil
			} else {
				return nil, errors.New("invalid user id")
			}
		},
		SecretKey: secretKey,
		Audience:  jwtAud,
	})

	router.GET("/", jwtAuth, func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			user, _ := jwt.CurrentUser(c).(*User)
			c.WriteIndentedJSON(http.StatusOK, kelly.H{
				"code": "0",
				"data": user.Name,
			})
		}
	})

	router.Run(":9999")
}
