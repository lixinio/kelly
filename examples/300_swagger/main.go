package main

import (
	"github.com/lixinio/kelly"
	swg "github.com/lixinio/kelly/middleware/swagger"
)

func main() {
	router := kelly.New(nil)
	r := router.Group("/api/v1/xx")

	swagger := swg.NewSwagger(
		r,
		&swg.Config{
			BasePath:    "/api/v1",
			Title:       "Swagger测试工具",
			Description: "Swagger测试工具描述",
			ApiVersion:  "0.1",
			Debug:       true,
		},
	)
	// handle /doc and /doc/spec
	swagger.RegisteSwaggerView(router.Group("/doc"))
	// r2 := router.Group("/doc")
	// r2.GET("/", swagger.SwaggerUIView("spec"))
	// r2.GET("/spec", swagger.SwaggerSpecView())

	r.GET("/test",
		swagger.SwaggerFile("swagger.yaml:test"),
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				c.ResponseStatusOK()
			}
		},
	)

	router.Run(":9999")
}
