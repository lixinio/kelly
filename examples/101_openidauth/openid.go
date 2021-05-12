package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lixinio/kelly"
	"gopkg.in/square/go-jose.v2"
)

const (
	AuthPath       = "/oauth/auth"
	TokenPath      = "/oauth/token"
	RevocationPath = "/oauth/revoke"

	DiscoveryPath = "/.well-known/openid-configuration"
	JWKPath       = "/.well-known/jwks.json"
	UserInfoPath  = "/userinfo"
)

func initOpenIDServer(router kelly.Router, keyPath, issuer string) *Handler {
	handler, err := NewHandler(keyPath, issuer)
	if err != nil {
		panic(err)
	}

	router.GET(DiscoveryPath, handler.handleDiscovery)
	router.GET(JWKPath, handler.handleJWK)
	return handler
}

type Handler struct {
	issuer string
	key    *rsa.PrivateKey
}

func NewHandler(path, issuer string) (*Handler, error) {
	key, err := loadRSAKey(path)
	if err != nil {
		return nil, fmt.Errorf("oauth: load rsa key error: %v", err)
	}
	return &Handler{
		issuer: issuer,
		key:    key,
	}, nil
}

func (h *Handler) handleJWK(ac *kelly.AnnotationContext) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		jwk := jose.JSONWebKey{
			Key:       h.key.Public(),
			Algorithm: "RS256",
			Use:       "sig",
		}

		setKeyID(&jwk)

		jwks := &jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				jwk,
			},
		}

		c.WriteJSON(http.StatusOK, jwks)
	}
}

func setKeyID(jwk *jose.JSONWebKey) {
	kid, _ := jwk.Thumbprint(crypto.SHA1)
	jwk.KeyID = hex.EncodeToString(kid)
}

func loadRSAKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load private key error: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("illegal rsa private key format")
	}

	pkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key error: %v", err)
	}

	return pkey, nil
}

func makeClaims(aud, sub, name, issuer string, ttl int) *jwt.MapClaims {
	now := time.Now()
	exp := now.Add(time.Duration(ttl) * time.Second).Unix()
	rclaims := &jwt.MapClaims{
		"aud":       aud,
		"auth_time": now.Unix(),
		"exp":       exp,
		"iat":       now.Unix(),
		"iss":       issuer,
		"sub":       sub,
		"name":      name,
	}

	return rclaims
}

func (h *Handler) issueToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ts, err := token.SignedString(h.key)
	if err != nil {
		return "", err
	}

	return ts, nil
}
