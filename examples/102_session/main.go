package main

import (
	"net/http"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/sessions"
)

func main() {
	router := kelly.New(nil)
	s, err := sessions.NewSessionManager(&sessions.SessionManagerConfig{
		Name:             "test_session",
		Secret:           "eejaijaecoonai1keuTh8iwee0pheiGh",
		RedisURL:         "redis://127.0.0.1:6379/4",
		SessionKeyPrefix: "session:example:",
		SessionLifetime:  1800,
	})
	if err != nil {
		panic(err)
	}

	router.GET("/", func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			session, closer, err := s.StartSession(c)
			if err != nil {
				panic(err)
			}

			session.Values["foo"] = "bar"
			session.Values["bar"] = "fpo"

			closer()
			c.WriteIndentedJSON(http.StatusOK, kelly.H{
				"code": "0",
			})
		}
	})

	router.Run(":9999")
}
