// https://github.com/martini-contrib/auth

package basic_auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/lixinio/kelly"
)

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
func secureCompare(given string, actual string) bool {
	givenSha := sha256.Sum256([]byte(given))
	actualSha := sha256.Sum256([]byte(actual))

	return subtle.ConstantTimeCompare(givenSha[:], actualSha[:]) == 1
}

// BasicRealm is used when setting the WWW-Authenticate response header.
const basicRealm = "Authorization Required"

const (
	// 设置到 kelly.Context的 存储 Key， 适配 CurrentUser
	contextDataKeyBasicUser string = "middleware.basic.user"
)

// Basic returns a Handler that authenticates via Basic Auth. Writes a http.StatusUnauthorized
// if authentication fails.
func BasicAuth(username string, password string) kelly.HandlerFunc {
	var siteAuth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return func(c *kelly.Context) {
		if auth, err := c.GetHeader("Authorization"); err == nil {
			if !secureCompare(auth, "Basic "+siteAuth) {
				unauthorized(c)
				return
			}
			c.Set(contextDataKeyBasicUser, username)
		} else {
			unauthorized(c)
		}
	}
}

// BasicFunc returns a Handler that authenticates via Basic Auth using the provided function.
// The function should return true for a valid username/password combination.
func BasicAuthFunc(authfn func(string, string) bool) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		auth, err := c.GetHeader("Authorization")
		if err != nil {
			unauthorized(c)
			return
		}

		if len(auth) < 6 || auth[:6] != "Basic " {
			unauthorized(c)
			return
		}
		b, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			unauthorized(c)
			return
		}
		tokens := strings.SplitN(string(b), ":", 2)
		if len(tokens) != 2 || !authfn(tokens[0], tokens[1]) {
			unauthorized(c)
			return
		}
		c.Set(contextDataKeyBasicUser, tokens[0])
	}
}

func unauthorized(res http.ResponseWriter) {
	res.Header().Set("WWW-Authenticate", "Basic realm=\""+basicRealm+"\"")
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}

// CurrentUser 获得当前用户
func CurrentUser(c *kelly.Context) interface{} {
	return c.MustGet(contextDataKeyBasicUser)
}
