package jwt

import (
	"errors"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// NewMapClaims 工具函数：创建一个 MapClaims 对象
func NewMapClaims(audience, issuer, subject string, ttl int) *MapClaims {
	now := time.Now()
	exp := now.Add(time.Duration(ttl) * time.Second).Unix()
	return &MapClaims{
		MapClaims: jwtgo.MapClaims{
			"aud": audience,
			"sub": subject,
			"exp": exp,
			"iat": now.Unix(),
			"iss": issuer,
		},
	}
}

// GenerateHS256Token 签发token
func GenerateHS256Token(claims Claims, secretKey string) (string, error) {
	if realClaims2, ok := claims.(*MapClaims); ok {
		claims = &realClaims2.MapClaims
	} else {
		return "", errors.New("unsupport claims object type")
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// verifyHS256Token 验证token
func verifyHS256Token(tokenStr, secretKey string, claims Claims) (*MapClaims, error) {
	realClaims, ok := claims.(*MapClaims)
	if !ok {
		return nil, errors.New("unsupport claims object type")
	}

	token, err := jwtgo.ParseWithClaims(tokenStr, &realClaims.MapClaims, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return realClaims, nil
}
