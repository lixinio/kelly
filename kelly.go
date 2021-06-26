package kelly

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Config 配置参数
type Config struct {
	// https://pkg.go.dev/github.com/julienschmidt/httprouter?utm_source=godoc#Router.RedirectTrailingSlash
	RedirectTrailingSlash bool
	// https://pkg.go.dev/github.com/julienschmidt/httprouter?utm_source=godoc#Router.RedirectFixedPath
	RedirectFixedPath bool
	// https://pkg.go.dev/github.com/julienschmidt/httprouter?utm_source=godoc#Router.HandleMethodNotAllowed
	HandleMethodNotAllowed HandlerFunc
	// https://pkg.go.dev/github.com/julienschmidt/httprouter?utm_source=godoc#Router.NotFound
	HandleNotFound HandlerFunc
	// 调试模式
	Debug bool
}

// Kelly 实例对象
type Kelly interface {
	Router
	http.Handler
	Run(addr string)                          // 同步启动
	RunContext(context.Context, string) error // 异步启动，等待context.Done
	RunTest(r *http.Request) *http.Response   // Debug
	RegistePreRunHandler(PreRunHandler)       // 注册正式运行前运行逻辑
}

type PreRunHandler func(Kelly)

type kellyImp struct {
	hr *httprouter.Router
	*router
	runBeforeHandlers []PreRunHandler // 监听端口前运行的逻辑
	config            *Config         //全局配置
	inited            bool            // 是否已经初始化
}

func (k *kellyImp) RegistePreRunHandler(handler PreRunHandler) {
	if handler == nil {
		panic("invalid PreRunHandler")
	}
	k.runBeforeHandlers = append(k.runBeforeHandlers, handler)
}

func (k *kellyImp) tryInit(addr string) {
	if k.inited {
		return
	}
	k.inited = true

	k.Use(LoggerRouter)
	k.router.doPreRun()

	for _, handler := range k.runBeforeHandlers {
		handler(k)
	}

	k.print(addr)
}

func (k *kellyImp) Run(addr string) {
	k.tryInit(addr)
	log.Fatal(http.ListenAndServe(addr, k.hr))
}

func (k *kellyImp) RunContext(context context.Context, addr string) error {
	k.tryInit(addr)
	srv := &http.Server{Addr: addr, Handler: k.hr}

	go func() {
		select {
		case <-context.Done():
			srv.Shutdown(context)
			break
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
		return err
	}
	return nil
}

func defaultHandleMethodNotAllowed(c *Context) {
	c.WriteString(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}

func defaultHandleNotFound(c *Context) {
	c.WriteString(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func defaultKellyConfig() *Config {
	return &Config{
		RedirectTrailingSlash: true,
		RedirectFixedPath:     false,
	}
}

func newImp(config *Config, handlers ...interface{}) Kelly {
	router := httprouter.New()
	if config == nil {
		config = defaultKellyConfig()
	}
	if config.HandleMethodNotAllowed == nil {
		config.HandleMethodNotAllowed = defaultHandleMethodNotAllowed
	}
	if config.HandleNotFound == nil {
		config.HandleNotFound = defaultHandleNotFound
	}

	router.RedirectTrailingSlash = config.RedirectTrailingSlash
	router.RedirectFixedPath = config.RedirectFixedPath
	router.NotFound = &handlerFuncWrap{config.HandleNotFound}
	router.MethodNotAllowed = &handlerFuncWrap{config.HandleMethodNotAllowed}

	ky := &kellyImp{
		hr:                router,
		config:            config,
		runBeforeHandlers: make([]PreRunHandler, 0),
	}
	ky.router = newRouterImp(router, ky, nil, "", "", handlers...)

	return ky
}

// New 创建一个新实例
func New(config *Config, handlers ...interface{}) Kelly {
	return newImp(config, handlers...)
}

func (k *kellyImp) RunTest(r *http.Request) *http.Response {
	k.tryInit("")

	mux := http.NewServeMux()
	mux.HandleFunc("/", k.router.ServeHTTP)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Result()
}

func (k *kellyImp) ServeHTTP(r http.ResponseWriter, w *http.Request) {
	k.tryInit("")
	k.router.ServeHTTP(r, w)
}

func (k *kellyImp) print(addr string) {
	txt := `
	
****************************************************************************

    @@          @@                       
    @@     /@@.                           @@          @@                       
    @@   #@@              #@@@@*          @@          @@       ###       ##    
    @@ %@@@             @@/    @@(        @@          @@        @@@     @@     
    @@@@  @@.          @@#//////@@        @@          @@         &@@   @@      
    @@     @@@         @@(                @@          @@          @@% @@*      
    @@       @@,        @@@     %         @@          @@           %@@@#       
                           ,%%%@           %%%         %%           @@@        
                                                                   @@&         
                                                                @@@       

****************************************************************************
    `
	fmt.Println(txt)
	fmt.Printf(" * Debug mode: %v\n", k.config.Debug)
	if strings.HasPrefix(addr, ":") {
		fmt.Printf(" * Running on 127.0.0.1%s\n", addr)
	} else {
		fmt.Printf(" * Running on %s\n", addr)
	}
}
