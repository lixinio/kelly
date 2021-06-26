package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/lixinio/kelly"
)

type AuthError string

func (err AuthError) Error() string {
	return string(err)
}

// 错误码
var (
	ErrGetTokenFail     error = AuthError("get token fail")               // 获取token失败
	ErrTokenVerifyFail        = AuthError("verify token fail")            // 校验token失败
	ErrTokenAuthFail          = AuthError("auth token fail")              // 认证失败
	ErrAudienceMissing        = AuthError("token audience missing")       // aud不存在
	ErrAudienceDismatch       = AuthError("token audience dismatch fail") // aud不匹配
	ErrTokenExpired           = AuthError("get token fail")               // token过期
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
type ErrorHandlerFunc func(*kelly.Context, error)

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

func defaultErrorHandler(c *kelly.Context, err error) {
	var aerr AuthError
	if errors.As(err, &aerr) {
		c.WriteJSON(http.StatusUnauthorized, kelly.H{
			"code":    http.StatusText(http.StatusUnauthorized),
			"message": aerr.Error(),
			"detail":  err.Error(),
		})
	} else {
		c.WriteJSON(http.StatusUnauthorized, kelly.H{
			"code":    http.StatusText(http.StatusUnauthorized),
			"message": err.Error(),
		})
	}
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

func JwtAuth(config *JwtAuthConfig) kelly.HandlerFunc {
	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultErrorHandler
	}

	if config.TokenGetter == nil {
		config.TokenGetter = defaultTokenGetter
	}

	return func(c *kelly.Context) {
		token, err := config.TokenGetter(c)
		if err != nil {
			config.ErrorHandler(c, fmt.Errorf("get token fail(%v) : %w", err, ErrGetTokenFail))
			return
		}

		claims, err := verifyHS256Token(token, config.SecretKey, defaultClaimsGetter())
		if err != nil {
			aerr := ErrTokenVerifyFail
			validationErr, ok := err.(*jwtgo.ValidationError)
			if ok || validationErr.Errors == jwtgo.ValidationErrorExpired {
				// token超时单独拎出来
				aerr = ErrTokenExpired
			}

			config.ErrorHandler(c, fmt.Errorf("verify fail(%v) : %w", err, aerr))
			return
		}

		if len(config.Audience) > 0 {
			if aud, ok := claims.Get("aud").(string); !ok {
				config.ErrorHandler(c, fmt.Errorf("claims audience missing : %w", ErrAudienceMissing))
				return
			} else if aud != config.Audience {
				config.ErrorHandler(c, fmt.Errorf("claims audience dismatch : %w", ErrAudienceDismatch))
				return
			}
		}

		user, err := config.Authorizator(claims)
		if err != nil {
			config.ErrorHandler(c, fmt.Errorf("auth fail(%v) : %w", err, ErrTokenAuthFail))
			return
		}

		c.Set(contextDataKeyJwtUser, user)
		c.InvokeNext()
	}
}
