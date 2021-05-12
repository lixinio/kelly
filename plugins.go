package kelly

import (
	"fmt"
	"time"
)

func LoggerRouter(ac *AnnotationContext) HandlerFunc {
	fmt.Printf("[Kelly] %v | %4s %s\n",
		time.Now().Format("2006/01/02 15:04:05"),
		ac.Method,
		ac.Path,
	)
	return nil
}

func Logger(ac *AnnotationContext) HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		path := c.Request().URL.Path
		raw := c.Request().URL.RawQuery

		c.InvokeNext()

		latency := time.Now().Sub(start)
		if raw != "" {
			path = path + "?" + raw
		}

		fmt.Printf("[Kelly] %v | %13s | %s %s\n",
			time.Now().Format("2006/01/02 15:04:05"),
			latency,
			c.Request().Method,
			path,
		)
	}
}
