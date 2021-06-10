package sentry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	sentrygo "github.com/getsentry/sentry-go"
	"github.com/lixinio/kelly"
)

// 参考
// https://github.com/getsentry/sentry-go/blob/v0.10.0/gin/sentrygin.go
// https://github.com/getsentry/sentry-go/blob/master/example/gin/main.go

var (
	ErrSentryConfig error = errors.New("invalid sentry config")
)

const (
	contextDataContext string        = "middleware.sentry.context"
	defaultEnv                       = "test"
	defaultWaitTime    time.Duration = 200
)

type SentryConfig struct {
	DSN              string
	Env              string
	Release          string
	Debug            bool          // 开启Debug模式
	AttachStacktrace bool          // 附带堆栈信息
	Repanic          bool          // 继续panic
	WaitForDelivery  bool          // 是否等待投递结束（避免请求过多拖累sentry服务器）
	WaitTimeout      time.Duration // 投递间隔毫秒

}

// Check for a broken connection, as this is what Gin does already.
func isBrokenPipeError(err interface{}) bool {
	if netErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(sysErr.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(sysErr.Error()), "connection reset by peer") {
				return true
			}
		}
	}
	return false
}

func recoverWithSentry(config *SentryConfig, hub *sentrygo.Hub, r *http.Request) {
	if err := recover(); err != nil {
		if !isBrokenPipeError(err) {
			eventID := hub.RecoverWithContext(
				context.WithValue(r.Context(), sentrygo.RequestContextKey, r),
				err,
			)
			if eventID != nil && config.WaitForDelivery {
				hub.Flush(config.WaitTimeout * time.Millisecond)
			}
		}
		if config.Repanic {
			panic(err)
		}
	}
}

func GetHubFromContext(c *kelly.Context) *sentrygo.Hub {
	if hub := c.Get(contextDataContext); hub != nil {
		if hub, ok := hub.(*sentrygo.Hub); ok {
			return hub
		}
	}
	return nil
}

func Sentry(config *SentryConfig) kelly.HandlerFunc {
	if config == nil {
		panic(fmt.Errorf("sentry config can NOT be empty : %w", ErrSentryConfig))
	}
	if len(config.DSN) < 1 {
		panic(fmt.Errorf("sentry dsn can NOT be empty : %w", ErrSentryConfig))
	}
	if config.Env == "" {
		config.Env = defaultEnv
	}
	if config.Debug {
		config.AttachStacktrace = true
	}
	if config.WaitForDelivery && config.WaitTimeout <= 0 {
		config.WaitTimeout = defaultWaitTime
	}

	err := sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              config.DSN,
		Environment:      config.Env,
		Release:          config.Release,
		Debug:            config.Debug,
		AttachStacktrace: config.AttachStacktrace,
		// BeforeSend: func(event *sentrygo.Event, hint *sentrygo.EventHint) *sentrygo.Event {
		// 	if hint.Context != nil {
		// 		if req, ok := hint.Context.Value(sentrygo.RequestContextKey).(*http.Request); ok {
		// 			// You have access to the original Request
		// 			fmt.Println(req)
		// 		}
		// 	}
		// 	fmt.Println(event)
		// 	return event
		// },
	})
	if err != nil {
		panic(fmt.Errorf("sentry init fail : %w(%s)", ErrSentryConfig, err))
	}

	return func(c *kelly.Context) {
		r := c.Request()
		hub := sentrygo.GetHubFromContext(r.Context())
		if hub == nil {
			hub = sentrygo.CurrentHub().Clone()
		}
		hub.Scope().SetRequest(r)
		c.Set(contextDataContext, hub)

		defer recoverWithSentry(config, hub, r)
		c.InvokeNext()
	}
}
