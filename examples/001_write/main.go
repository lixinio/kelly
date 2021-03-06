package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/lixinio/kelly"
)

func main() {
	router := kelly.New(nil)

	router.GET("/", func(c *kelly.Context) {
		c.WriteIndentedJSON(http.StatusOK, kelly.H{
			"message": "ok",
			"code":    "0",
		})
	})

	router.GET("/json", func(c *kelly.Context) {
		c.WriteJSON(http.StatusOK, kelly.H{
			"message": "ok",
			"code":    "0",
		})
	})

	router.GET("/str", func(c *kelly.Context) {
		c.WriteString(http.StatusOK, "你好 %s， 你好 %s", "世界", "中国")
	})

	router.GET("/xml", func(c *kelly.Context) {
		c.WriteXML(http.StatusOK, kelly.H{ // 返回XML
			"code": "/aaa/bbb",
		})
	})

	router.GET("/redirect", func(c *kelly.Context) {
		c.Redirect(http.StatusFound, "http://www.baidu.com")
	})

	router.GET("/html", func(c *kelly.Context) {
		data := `<html>
<body>
	<form action="#" method="post">
	<p>AAA: <input type="text" name="aaa" /></p>
	<p>BBB: <input type="text" name="bbb" /></p>
	<p>CCC: <input type="text" name="ccc" /></p>
	<input type="submit" value="Submit" />
	</form>
</body>
</html>`
		c.WriteHTML(http.StatusOK, data)
	})

	router.GET("/template", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		data := `<form action="#" method="get">
<p>First {{ .First }}: <input type="text" name="fname" value="{{ .First }}"/></p>
<p>Last {{ .Last }}: <input type="text" name="lname" value="{{ .Last }}"/></p>
<input type="submit" value="Submit" />
</form>`

		// 通过闭包预先编译好
		t := template.Must(template.New("t1").Parse(data))
		return func(c *kelly.Context) {
			c.WriteTemplateHTML(http.StatusOK, t, map[string]string{
				"First": "Qiu",
				"Last":  "King",
			})
		}
	})

	router.GET("/image", func(c *kelly.Context) {
		response, err := http.Get("https://mat1.gtimg.com/pingjs/ext2020/qqindex2018/dist/img/qq_logo_2x.png")
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			panic(errors.New("Received non 200 response code"))
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		c.WriteData(http.StatusOK, "image/png", body)
	})

	router.Run(":9999")
}
