package sessions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lixinio/kelly"
)

var (
	// 配置为空
	ErrLMConfigError error = errors.New("login manager config is empty")
	// 用户获取回调为空
	ErrLMConfigUserGetterError error = errors.New("login manager user getter is empty")
	// session 管理器为空
	ErrLMConfigSMError error = errors.New("login manager config session is empty")
)

const (
	SESSION_USER_ID = "user-id"
)

type User interface {
	GetID() interface{}
}

type UserGetter func(interface{}) interface{}

type LoginManagerConfig struct {
	UnauthorizedHandler kelly.HandlerFunc
	UserGetter          UserGetter
	SessionManager      *SessionManager
}

type LoginManager struct {
	config         *LoginManagerConfig
	sessionManager *SessionManager
}

func NewLoginManager(config *LoginManagerConfig) (*LoginManager, error) {
	if config == nil {
		return nil, ErrLMConfigError
	}
	if config.UserGetter == nil {
		return nil, ErrLMConfigUserGetterError
	}
	if config.SessionManager == nil {
		return nil, ErrLMConfigSMError
	}
	if config.UnauthorizedHandler == nil {
		config.UnauthorizedHandler = func(c *kelly.Context) {
			c.WriteString(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}
	}

	loginManager := &LoginManager{
		config:         config,
		sessionManager: config.SessionManager,
	}
	return loginManager, nil
}

func (loginManager *LoginManager) Login(c *kelly.Context, user User) error {
	s, closer, err := loginManager.sessionManager.StartSession(c)
	if err != nil {
		return err
	}
	s.Values[SESSION_USER_ID] = user.GetID()
	closer()
	return nil
}

func (loginManager *LoginManager) Logout(c *kelly.Context) error {
	s, closer, err := loginManager.sessionManager.StartSession(c)
	if err != nil {
		return err
	}
	s.Options.MaxAge = -1
	delete(s.Values, SESSION_USER_ID)
	closer()
	return nil
}

func (loginManager *LoginManager) GetCurrentUser(c *kelly.Context) interface{} {
	s, _, err := loginManager.sessionManager.StartSession(c)
	if err != nil {
		return nil
	}

	userid, ok := s.Values[SESSION_USER_ID]
	if !ok {
		return nil
	}

	return loginManager.config.UserGetter(userid)
}

func (loginManager *LoginManager) IsAuthenticated(c *kelly.Context) bool {
	return loginManager.GetCurrentUser(c) != nil
}

func (loginManager *LoginManager) MustGetCurrentUser(c *kelly.Context) interface{} {
	user := loginManager.GetCurrentUser(c)
	if user != nil {
		return user
	}
	panic(fmt.Errorf("can NOT get current user"))
}

func (loginManager *LoginManager) LoginRequired() kelly.HandlerFunc {
	return func(c *kelly.Context) {
		if loginManager.IsAuthenticated(c) {
			c.InvokeNext()
			return
		}
		c.Abort(http.StatusUnauthorized, "x")
	}
}
