package main

import (
	"fmt"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/middleware/sessions"
)

type User struct {
	id   string
	name string
}

func (user *User) GetID() interface{} {
	return user.id
}

func main() {
	router := kelly.New(nil)

	sessionManager, err := sessions.NewSessionManager(&sessions.SessionManagerConfig{
		Name:             "session",
		Secret:           "ungeejai2ohH8Ahchohmee9gohchie2E",
		RedisURL:         "redis://127.0.0.1:6379/3",
		SessionKeyPrefix: "session:example:",
		SessionLifetime:  1800,
	})
	if err != nil {
		panic(err)
	}
	loginManager, err := sessions.NewLoginManager(&sessions.LoginManagerConfig{
		SessionManager: sessionManager,
		UserGetter: func(id interface{}) interface{} {
			return &User{
				id:   id.(string),
				name: fmt.Sprintf("%s_name", id.(string)),
			}
		},
	})
	if err != nil {
		panic(err)
	}

	InitApiV1(router, loginManager)
	router.Run(":9999")
}
