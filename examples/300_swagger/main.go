package main

import (
	"net/http"
	"time"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/cors"
	swg "github.com/lixinio/kelly/middleware/swagger"
)

func main() {
	router := kelly.New(nil)
	r := router.Group("/api/v1/xx")

	router.Use(cors.Cors(cors.Config{
		AllowOrigins:     []string{"http://d.lan.lixinio.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Access-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 绑定所有的options请求来支持中间件作跨域处理
	router.OPTIONS("/*path", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.WriteString(http.StatusOK, "ok")
		}
	})

	swagger := swg.NewSwagger(
		r,
		&swg.Config{
			BasePath:     "/api/v1",
			Title:        "Swagger测试工具",
			Description:  "Swagger测试工具描述",
			ApiVersion:   "0.1",
			Debug:        true,
			SwaggerUiUrl: "http://d.lan.lixinio.com/",
			Headers: []swg.SecurityDefinition{
				{
					Name: "access-token",
				},
			},
		},
	)
	// handle /doc and /doc/spec
	swagger.RegisteSwaggerView(router.Group("/doc"))
	// r2 := router.Group("/doc")
	// r2.GET("/", swagger.SwaggerUIView("spec"))
	// r2.GET("/spec", swagger.SwaggerSpecView())

	r.GET("/test",
		swagger.SwaggerFile("swagger.yaml:test"),
		func(c *kelly.Context) {
			c.ResponseStatusOK()
		},
	)

	router.Run(":9999")
}
