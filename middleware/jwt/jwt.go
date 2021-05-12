package jwt

import (
	"errors"
	"net/http"
	"strings"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/lixinio/kelly"
)

// 错误码
const (
	ErrGetTokenFail     int = 10000 // 获取token失败
	ErrTokenVerifyFail      = 10001 // 校验token失败
	ErrTokenAuthFail        = 10002 // 认证失败
	ErrAudienceMissing      = 10003 // sub不存在
	ErrAudienceDismatch     = 10004 // sub不匹配
	ErrTokenExpired         = 10005 // token过期
)

const (
	// 设置到 kelly.Context的 存储 Key， 适配 CurrentUser
	contextDataKeyJwtUser string = "middleware.jwt.user"
	// 标准的认证头
	AuthnHeader string = "Authorization"
	AuthnType          = "bearer"
)

// TokenGetter 怎么从请求中获取Token
type TokenGetterFunc func(*kelly.Context) (string, error)

// AuthorizatorFunc 根据解析后的claims验证数据有效性
type AuthorizatorFunc func(Claims) (interface{}, error)

// ErrorHandlerFunc 错误处理函数
type ErrorHandlerFunc func(*kelly.Context, int, error)

type JwtAuthConfig struct {
	TokenGetter  TokenGetterFunc
	Authorizator AuthorizatorFunc
	ErrorHandler ErrorHandlerFunc
	Audience     string
	SecretKey    string
}

func defaultClaimsGetter() Claims {
	return &MapClaims{}
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
	return c.MustGet(contextDataKeyJwtUser)
}

func JwtAuth(config *JwtAuthConfig) kelly.AnnotationHandlerFunc {
	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultErrorHandler
	}

	if config.TokenGetter == nil {
		config.TokenGetter = defaultTokenGetter
	}

	return func(ac *kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			token, err := config.TokenGetter(c)
			if err != nil {
				config.ErrorHandler(c, ErrGetTokenFail, err)
				return
			}

			claims, err := verifyHS256Token(token, config.SecretKey, defaultClaimsGetter())
			if err != nil {
				code := ErrTokenVerifyFail
				validationErr, ok := err.(*jwtgo.ValidationError)
				if ok || validationErr.Errors == jwtgo.ValidationErrorExpired {
					// token超时单独拎出来
					code = ErrTokenExpired
				}

				config.ErrorHandler(c, code, err)
				return
			}

			if len(config.Audience) > 0 {
				if aud, ok := claims.Get("aud").(string); !ok {
					config.ErrorHandler(c, ErrAudienceMissing, errors.New("claims audience missing"))
					return
				} else if aud != config.Audience {
					config.ErrorHandler(c, ErrAudienceDismatch, errors.New("claims audience dismatch"))
					return
				}
			}

			user, err := config.Authorizator(claims)
			if err != nil {
				config.ErrorHandler(c, ErrTokenAuthFail, err)
				return
			}

			c.Set(contextDataKeyJwtUser, user)
			c.InvokeNext()
		}
	}
}
