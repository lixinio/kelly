package swagger

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lixinio/kelly"
)

func (s *Swagger) RegisteSwaggerView(r kelly.Router) {
	spec := "spec"
	r.GET("/", s.SwaggerUIView(spec))
	r.GET(fmt.Sprintf("/%s", spec), s.SwaggerSpecView())
}

func (s *Swagger) SwaggerUIView(spec string) kelly.AnnotationHandlerFunc {
	return func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		path := ac.Path
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}

		return func(c *kelly.Context) {
			scheme := "http://"
			if c.Request().TLS != nil {
				scheme = "https://"
			}
			host := fmt.Sprintf("%s%s%s%s", scheme, c.Request().Host, ac.Path, spec)
			host = s.config.SwaggerUiUrl + "?url=" + url.QueryEscape(host)
			c.Redirect(http.StatusFound, host)
		}
	}
}

func (s *Swagger) SwaggerSpecView() kelly.AnnotationHandlerFunc {
	return func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.SetHeader("Access-Control-Allow-Origin", "*")
			c.WriteRawJSON(http.StatusOK, s.specData)
		}
	}
}
