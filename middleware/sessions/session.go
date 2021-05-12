package sessions

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	gsessions "github.com/gorilla/sessions"
	"github.com/lixinio/kelly"
)

var (
	// 配置为空
	ErrConfigError error = errors.New("session manager config is empty")
	// 配置 name 为空
	ErrConfigNameError error = errors.New("session manager config name is empty")
	// 配置 name 为空
	ErrConfigSecretError error = errors.New("session manager config secret is empty")
	// 配置 store 为空
	ErrConfigStoreError error = errors.New("session manager config store is empty")
	// redis url 错误
	ErrRedisUrlError error = errors.New("parse redis url fail")
)

const (
	// 设置到 kelly.Context的 存储 Key， 适配 CurrentUser
	contextDataKeySession string = "middleware.session"
	// 缺省redis pool size
	defaultRedisPoolSize = 20
	// 缺省的session有效期
	defaultSessionLifetime = 3600 * 2
)

type SessionManagerConfig struct {
	Name             string          // cookie的名称
	Secret           string          // cookie密钥
	RedisURL         string          // redis
	SessionLifetime  int             // session有效期（秒）
	SessionKeyPrefix string          // session前缀
	Store            gsessions.Store // 存储session的容器（cookie/redis/mysql etc.)
}

type SessionManager struct {
	config *SessionManagerConfig
	store  gsessions.Store
}

func parseRedisURL(urlStr string) (int, string, string, error) {
	redisURL, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", "", err
	}

	redisPwd := ""
	if redisURL.User != nil {
		if password, ok := redisURL.User.Password(); ok {
			redisPwd = password
		}
	}

	redisDb := 0
	if len(redisURL.Path) > 1 {
		db := strings.TrimPrefix(redisURL.Path, "/")
		intVar, err := strconv.Atoi(db)
		if err != nil {
			return 0, "", "", err
		}
		redisDb = intVar
	}

	return redisDb, redisURL.Host, redisPwd, nil
}

func NewSessionManager(config *SessionManagerConfig) (*SessionManager, error) {
	if config == nil {
		return nil, ErrConfigError
	}
	if len(config.Name) == 0 {
		return nil, ErrConfigNameError
	}
	if len(config.Secret) == 0 {
		return nil, ErrConfigSecretError
	}
	if len(config.RedisURL) == 0 && config.Store == nil {
		return nil, ErrConfigStoreError
	}
	if config.SessionLifetime <= 0 {
		config.SessionLifetime = defaultSessionLifetime
	}

	var store sessions.Store = nil
	if config.Store == nil {
		redisDB, redisHost, redisPwd, err := parseRedisURL(config.RedisURL)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid redis url %s(%s), error : %w",
				config.RedisURL, err.Error(), ErrRedisUrlError,
			)
		}
		redis_store, err := redistore.NewRediStoreWithDB(
			defaultRedisPoolSize, "tcp", redisHost, redisPwd,
			strconv.Itoa(redisDB), []byte(config.Secret),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"new redis store %s(%s) fail, error : %w",
				config.RedisURL, err.Error(), ErrRedisUrlError,
			)
		}
		if len(config.SessionKeyPrefix) != 0 {
			redis_store.SetKeyPrefix(config.SessionKeyPrefix)
		}
		store = redis_store
	} else {
		store = config.Store
	}

	session := &SessionManager{
		config: config,
		store:  store,
	}
	return session, nil
}

func (session *SessionManager) StartSession(c *kelly.Context) (*sessions.Session, func() error, error) {
	s, err := session.store.Get(c.Request(), session.config.Name)
	if err != nil {
		return nil, nil, err
	} else {
		s.Options.MaxAge = session.config.SessionLifetime
		return s, func() error {
			return s.Save(c.Request(), c.ResponseWriter)
		}, nil
	}
}
