package jwt

import (
	jwtgo "github.com/dgrijalva/jwt-go"
)

type Claims jwtgo.Claims

type MapClaims struct {
	jwtgo.MapClaims
}

func (mapClaims *MapClaims) Update(key string, value interface{}) *MapClaims {
	mapClaims.MapClaims[key] = value
	return mapClaims
}

func (mapClaims *MapClaims) Get(key string) interface{} {
	value, ok := mapClaims.MapClaims[key]
	if ok {
		return value
	} else {
		return nil
	}
}
