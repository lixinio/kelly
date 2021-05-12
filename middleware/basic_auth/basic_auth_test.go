package basic_auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/lixinio/kelly"
	"github.com/lixinio/kelly/test"
)

func TestBasicAuth(t *testing.T) {
	username, password := "username", "password"
	hash := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	f := func(*kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			c.ResponseStatusOK()
		}
	}

	f2 := func(*kelly.AnnotationContext) kelly.HandlerFunc {
		return func(c *kelly.Context) {
			if CurrentUser(c) != username {
				t.Errorf("basic auth fail %v|%v", CurrentUser(c), username)
			}
			c.ResponseStatusOK()
		}
	}

	middleware := BasicAuth(username, password)
	middleware2 := BasicAuthFunc(func(u, p string) bool {
		return u == username && p == password
	})
	for _, m := range []kelly.AnnotationHandlerFunc{middleware, middleware2} {
		expectCode := http.StatusUnauthorized
		resp := test.KellyFramwork("/", "/", map[string]string{}, m, f)
		if resp.StatusCode != expectCode {
			t.Errorf("basic auth fail %d|%d", resp.StatusCode, expectCode)
		}
		resp = test.KellyFramwork("/", "/", map[string]string{"Authorization": hash}, m, f)
		if resp.StatusCode != expectCode {
			t.Errorf("basic auth fail %d|%d", resp.StatusCode, expectCode)
		}

		expectCode = http.StatusOK
		resp = test.KellyFramwork(
			"/", "/", map[string]string{"Authorization": fmt.Sprintf("Basic %s", hash)}, m, f2,
		)
		if resp.StatusCode != expectCode {
			t.Errorf("basic auth fail %d|%d", resp.StatusCode, expectCode)
		}
	}

}
