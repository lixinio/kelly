package static

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/lixinio/kelly"
)

type Config struct {
	Dir           http.FileSystem   // 托管的目录 eg. http.Dir("/var/www/html")
	EnableListDir bool              // 是否支持枚举文件
	Handler404    kelly.HandlerFunc // 文件/目录不存在的处理
	Indexfiles    []string          // 主页文件
}

// 根据一个目录生成一个HandlerFunc处理文件请求，在绑定Path时，必须使用下面的规则
// r.GET("/static/*path", kelly.Static(http.Dir("/tmp")))
// 在内部依赖于名称为path的路径变量
// 若将*path改成:path，将只能访问根目录的文件，无法嵌套
func Static(config *Config) kelly.AnnotationHandlerFunc {
	staticTemp := `<pre>
{{ range $key, $value := . }}
	<a href="{{ $value.Url }}" style="color: {{ $value.Color }};">{{ $value.Name }}</a>
{{ end }}
</pre>`

	// 初始化模板
	t := template.Must(template.New("staticTemp").Parse(staticTemp))
	if len(config.Indexfiles) > 0 {
		for _, v := range config.Indexfiles {
			if len(v) < 1 {
				panic(fmt.Errorf("invalid index file"))
			} else if strings.ContainsAny(v, "/") {
				panic(fmt.Errorf("invalid index file %s", v))
			}
		}
	}

	// 错误处理
	handler404 := config.Handler404
	if handler404 == nil {
		handler404 = func(c *kelly.Context) {
			c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
	}

	return func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			// 获得Path变量
			file := c.MustGetPathVarible("path")
			fmt.Println("path ", file)
			f, err := config.Dir.Open(file)
			if err != nil {
				handler404(c)
				return
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				handler404(c)
				return
			}

			// 处理文件
			if fi.IsDir() {
				fmt.Println("serve path ", file)
				if len(config.Indexfiles) > 0 {
					if !serverIndex(config, file, c) {
						// 如果找不到index， 又支持枚举
						listDir(config, f, t, c)
						return
					}
				}
				if config.EnableListDir {
					listDir(config, f, t, c)
				}
				return
			}

			http.ServeContent(c, c.Request(), file, fi.ModTime(), f)
		}
	}
}

func serverIndex(config *Config, file string, c *kelly.Context) bool {
	// 自动处理首页的情况 eg. index.html
	var target = ""
	for _, v := range config.Indexfiles {
		newFile := path.Join(file, v)
		f, err := config.Dir.Open(newFile)
		if err != nil {
			continue
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil || fi.IsDir() {
			continue
		}

		target = v
		break
	}

	if target == "" {
		// 如果当前目录匹配不到index
		if config.EnableListDir {
			return false
		} else {
			config.Handler404(c)
		}
	} else {
		c.Redirect(http.StatusFound, target)
	}
	return true
}

// 参考 https://github.com/labstack/echo/blob/master/middleware/static.go
func listDir(config *Config, d http.File, t *template.Template, c *kelly.Context) {
	dirs, err := d.Readdir(-1)
	if err != nil {
		config.Handler404(c)
		return
	}

	data := []map[string]string{}
	for _, d := range dirs {
		name := d.Name()
		color := "#212121"
		if d.IsDir() {
			color = "#e91e63"
			name += "/"
		}

		data = append(data, map[string]string{
			"Name":  name,
			"Color": color,
			"Url":   name,
		})
	}

	c.WriteTemplateHTML(http.StatusOK, t, data)
}
