package openid

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/lixinio/kelly"
)

// 错误码
const (
	ErrGetTokenFail     int = 10000 // 获取token失败
	ErrTokenVerifyFail      = 10001 // 校验token失败
	ErrTokenInvalidType     = 10002 // 认证失败
	ErrTokenAuthFail        = 10003 // 认证失败
	ErrTokenExpired         = 10004 // token过期
)

type MapClaims map[string]interface{}

const (
	// 设置到 kelly.Context的 存储 Key， 适配 CurrentUser
	contextDataKeyOpenIDUser string = "middleware.openid.user"
	// 标准的认证头
	AuthnHeader string = "Authorization"
	AuthnType          = "bearer"
)

// TokenGetter 怎么从请求中获取Token
type TokenGetterFunc func(*kelly.Context) (string, error)

// AuthorizatorFunc 根据解析后的claims验证数据有效性
type AuthorizatorFunc func(*MapClaims) (interface{}, error)

// ErrorHandlerFunc 错误处理函数
type ErrorHandlerFunc func(*kelly.Context, int, error)

type OpenIDAuthConfig struct {
	TokenGetter  TokenGetterFunc
	Authorizator AuthorizatorFunc
	ErrorHandler ErrorHandlerFunc
	Issuer       string
	Audience     string
}

func defaultErrorHandler(c *kelly.Context, code int, err error) {
	c.Abort(http.StatusUnauthorized, err.Error())
}

func defaultTokenGetter(c *kelly.Context) (string, error) {
	authn, err := c.GetHeader(AuthnHeader)
	if err != nil {
		return "", err
	}

	splits := strings.Split(authn, " ")
	if len(splits) != 2 {
		return "", errors.New("invalid bearer token format")
	}

	if strings.ToLower(splits[0]) != AuthnType {
		return "", errors.New("invalid bearer token format")
	}

	return splits[1], nil
}

// CurrentUser 获得当前用户
func CurrentUser(c *kelly.Context) interface{} {
	return c.MustGet(contextDataKeyOpenIDUser)
}

func OpenIDAuth(config *OpenIDAuthConfig) (kelly.AnnotationHandlerFunc, error) {
	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultErrorHandler
	}

	if config.TokenGetter == nil {
		config.TokenGetter = defaultTokenGetter
	}

	provider, err := oidc.NewProvider(context.Background(), config.Issuer)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.Audience,
	})

	return func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			token, err := config.TokenGetter(c)
			if err != nil {
				config.ErrorHandler(c, ErrGetTokenFail, err)
				return
			}

			idtoken, err := verifier.Verify(c.Request().Context(), token)
			if err != nil {
				config.ErrorHandler(c, ErrTokenVerifyFail, err)
				return
			}

			var claims MapClaims
			if err := idtoken.Claims(&claims); err != nil {
				config.ErrorHandler(c, ErrTokenInvalidType, err)
				return
			}

			user, err := config.Authorizator(&claims)
			if err != nil {
				config.ErrorHandler(c, ErrTokenAuthFail, err)
				return
			}

			c.Set(contextDataKeyOpenIDUser, user)
			c.InvokeNext()
		}
	}, nil
}
