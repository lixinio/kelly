package sentry

import (
	"errors"
	"testing"
	"time"

	sentrygo "github.com/getsentry/sentry-go"
	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/test"
)

var (
	ErrTest1 error = errors.New("panic error 1")
	ErrTest2 error = errors.New("panic error 2")
	ErrTest3 error = errors.New("panic error 3")
)

func TestSentry(t *testing.T) {
	func() {
		defer test.CheckError(t, ErrSentryConfig)
		Sentry(nil)
	}()
	func() {
		defer test.CheckError(t, ErrSentryConfig)
		Sentry(&SentryConfig{})
	}()

	config := &SentryConfig{
		DSN:              "https://www.baidu.com/530",
		Env:              "test",
		Release:          "1.1.1",
		Debug:            true,
		AttachStacktrace: true,
		Repanic:          true,
		WaitForDelivery:  true,
	}

	func() {
		defer test.CheckError(t, ErrTest1)
		test.KellyFramwork("/", "/", map[string]string{},
			Sentry(config),
			func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
				return func(c *kelly.Context) {
					panic(ErrTest1)
				}
			},
		)
	}()
	func() {
		defer test.CheckError(t, ErrTest2)
		test.KellyFramwork("/", "/", map[string]string{},
			Sentry(config),
			func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
				return func(c *kelly.Context) {
					panic(ErrTest2)
				}
			},
		)
	}()

	config.Repanic = false
	test.KellyFramwork("/", "/", map[string]string{},
		Sentry(config), func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				panic(ErrTest1)
			}
		},
	)
	test.KellyFramwork("/", "/", map[string]string{},
		Sentry(config),
		func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
			return func(c *kelly.Context) {
				panic(ErrTest2)
			}
		},
	)
}

func TestSentryContext(t *testing.T) {
	config := &SentryConfig{
		DSN:              "https://www.baidu.com/530",
		Env:              "test",
		Release:          "1.1.1",
		Debug:            true,
		AttachStacktrace: true,
	}

	for _, f := range []func(*sentrygo.Hub){
		func(hub *sentrygo.Hub) { hub.CaptureMessage("capture message") },
		func(hub *sentrygo.Hub) { hub.CaptureException(ErrTest3) },
	} {
		test.KellyFramwork("/", "/", map[string]string{},
			Sentry(config),
			func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
				return func(c *kelly.Context) {
					hub := GetHubFromContext(c)
					if hub == nil {
						t.Fatal("no hub")
					}

					hub.WithScope(func(scope *sentrygo.Scope) {
						scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
						f(hub)
					})
					c.ResponseStatusOK()
				}
			},
		)
	}
	sentrygo.Flush(time.Second)
}
