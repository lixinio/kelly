package kelly

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHandleResp(t *testing.T) {
	key := "query"
	f := func(value string) func(*Context) {
		return func(c *Context) {
			c.Header().Add(key, value)
			c.SetCookie(key, value, 0, "", "", false, false)
			c.WriteString(http.StatusOK, value)
		}
	}

	qhandler := func(key, value string) {
		resp := getFramwork(f(value), "/", map[string]string{}, map[string]string{}, "")

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Response code is %v|[%s][%s]", resp.StatusCode, key, value)
			return
		}

		result := readBody(resp)
		if result != value {
			t.Errorf("Response body is %v|[%s][%s]", result, key, value)
			return
		}
		if readHeader(resp, key) != value {
			t.Errorf("Response header is %v|[%s][%s]", readHeader(resp, key), key, value)
			return
		}
		if readCookie(resp, key) != value {
			t.Errorf("Response cookie is %v|[%s][%s]", readCookie(resp, key), key, value)
			return
		}
	}

	qhandler(key, "abc")
	qhandler(key, "")
}

func TestHandleAbort(t *testing.T) {
	f := func(code int, value string) func(*Context) {
		return func(c *Context) {
			c.Abort(code, value)
		}
	}

	qhandler := func(code int, value string) {
		resp := getFramwork(f(code, value), "/", map[string]string{}, map[string]string{}, "")

		if resp.StatusCode != code {
			t.Errorf("Response code is %v|[%s]", resp.StatusCode, value)
			return
		}

		if len(value) == 0 {
			value = http.StatusText(code)
		}
		expectResult, err := json.Marshal(H{
			"code":    code,
			"message": value,
		})
		if err != nil {
			t.Errorf("json.Marshal error %w", err)
			return
		}
		result := readBody(resp)
		if result != strings.Trim(string(expectResult), " \n") {
			t.Errorf("Response body is %v|[%s]", result, expectResult)
			return
		}
	}

	qhandler(http.StatusCreated, "")
	qhandler(http.StatusFound, "found")
	qhandler(http.StatusNotFound, "not found")
	qhandler(http.StatusServiceUnavailable, "error")
}

func TestHandleWriteRawJson(t *testing.T) {
	expectResult, _ := json.Marshal(H{
		"code":    http.StatusOK,
		"message": "abcdefg",
	})
	f := func(c *Context) {
		c.WriteRawJSON(http.StatusOK, expectResult)
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteRawJson fail")
		return
	}

	if readContentType(resp) != "application/json" {
		t.Errorf("WriteRawJson content type fail %s", readContentType(resp))
		return
	}

	result := readBody(resp)
	if result != strings.Trim(string(expectResult), " \n") {
		t.Errorf("WriteRawJson is %v|[%s]", result, expectResult)
		return
	}
}

func TestHandleWriteString(t *testing.T) {
	expectResult := "abcdefg"
	f := func(c *Context) {
		c.WriteString(http.StatusOK, expectResult)
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteString fail")
		return
	}

	if readContentType(resp) != "text/plain" {
		t.Errorf("WriteString content type fail %s", readContentType(resp))
		return
	}

	result := readBody(resp)
	if result != strings.Trim(string(expectResult), " \n") {
		t.Errorf("WriteString is %v|[%s]", result, expectResult)
		return
	}
}

func TestHandleWriteJson(t *testing.T) {
	body := H{
		"code":    http.StatusOK,
		"message": "abcdefg",
	}
	expectResult, _ := json.Marshal(body)
	f := func(c *Context) {
		c.WriteJSON(http.StatusOK, body)
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteJson fail")
		return
	}

	if readContentType(resp) != "application/json" {
		t.Errorf("WriteJson content type fail %s", readContentType(resp))
		return
	}

	result := readBody(resp)
	if result != strings.Trim(string(expectResult), " \n") {
		t.Errorf("WriteJson is %v|[%s]", result, expectResult)
		return
	}
}

func TestHandleWriteIndentedJson(t *testing.T) {
	body := H{
		"code":    http.StatusOK,
		"message": "abcdefg",
	}
	expectResult, _ := json.MarshalIndent(body, "", "    ")
	f := func(c *Context) {
		c.WriteIndentedJSON(http.StatusOK, body)
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteIndentedJson fail")
		return
	}

	if readContentType(resp) != "application/json" {
		t.Errorf("WriteIndentedJson content type fail %s", readContentType(resp))
		return
	}

	result := readBody(resp)
	if result != strings.Trim(string(expectResult), " \n") {
		t.Errorf("WriteIndentedJson is %v|[%s]", result, expectResult)
		return
	}
}

func TestHandleWriteXml(t *testing.T) {
	// xml本身不区分int
	body := &H{
		"code":    "ok",
		"message": "abcdefg",
	}
	// expectResult, _ := xml.Marshal(body)
	f := func(c *Context) {
		c.WriteXML(http.StatusOK, body)
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteXml fail")
		return
	}

	if readContentType(resp) != "application/xml" {
		t.Errorf("WriteXml content type fail %s", readContentType(resp))
		return
	}

	result := readRawBody(resp)
	expectResult := &H{}
	err := xml.Unmarshal(result, expectResult)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(body, expectResult) {
		t.Errorf("WriteXml cmp fail")
	}
}

func TestHandleWriteHtml(t *testing.T) {
	body := `
	<html>
		<head>
		</head>
		<body>
		</body>
	</html>
	`
	f := func(c *Context) {
		c.WriteHTML(http.StatusOK, strings.Trim(body, " \n"))
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteHtml fail")
		return
	}

	if readContentType(resp) != "text/html" {
		t.Errorf("WriteHtml content type fail %s", readContentType(resp))
		return
	}

	result := readBody(resp)
	if result != strings.Trim(body, " \n") {
		t.Errorf("WriteHtml is %v|[%s]", result, body)
		return
	}
}

func TestHandleWriteTemplateHtml(t *testing.T) {
	body := `
	<html>
		<head>
		</head>
		<body>
		{{ .body }}
		</body>
	</html>
	`

	tmpl := template.Must(template.New("t1").Parse(body))
	data := H{
		"body": "fadsfasdf",
	}
	f := func(c *Context) {
		c.WriteTemplateHTML(http.StatusOK, tmpl, data)
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteTemplateHtml fail")
		return
	}

	if readContentType(resp) != "text/html" {
		t.Errorf("WriteTemplateHtml content type fail %s", readContentType(resp))
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		t.Errorf("WriteTemplateHtml Execute fail %w", err)
		return
	}

	result := readBody(resp)
	if result != strings.Trim(tpl.String(), " \n") {
		t.Errorf("WriteTemplateHtml is %v|[%s]", result, body)
		return
	}
}

func TestHandleWriteData(t *testing.T) {
	data := newImage()

	f := func(c *Context) {
		c.WriteData(http.StatusOK, "image/png", data)
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WriteData fail")
		return
	}

	if readContentType(resp) != "image/png" {
		t.Errorf("WriteData content type fail %s", readContentType(resp))
		return
	}

	result := readRawBody(resp)
	if !cmp.Equal(result, data) {
		t.Errorf("WriteData Equal fail")
		return
	}
}

func TestHandleMisc(t *testing.T) {
	f := func(c *Context) {
		c.ResponseStatusOK()
	}
	resp := getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("ResponseStatusOK fail")
		return
	}

	redirectURL := "https://www.baidu.com"
	f = func(c *Context) {
		c.Redirect(http.StatusFound, redirectURL)
	}
	resp = getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Redirect fail %d|%d", resp.StatusCode, http.StatusOK)
		return
	}
	url, err := resp.Location()
	if err != nil || url.String() != redirectURL {
		t.Errorf("Redirect fail")
		return
	}

	errorMsg := "error"
	qhandler := func(f func(c *Context), code int, identify string) {
		resp = getFramwork(f, "/", map[string]string{}, map[string]string{}, "")
		if resp.StatusCode != code {
			t.Errorf("%s fail (%d|%d)", identify, resp.StatusCode, code)
			return
		}
		expectResult, _ := json.Marshal(H{
			"code":    code,
			"message": errorMsg,
		})
		result := readBody(resp)
		if result != strings.Trim(string(expectResult), " \n") {
			t.Errorf("%s is %v|[%s]", identify, result, expectResult)
			return
		}
	}
	qhandler(func(c *Context) {
		c.ResponseStatusBadRequest(errors.New(errorMsg))
	}, http.StatusBadRequest, "ResponseStatusBadRequest")

	qhandler(func(c *Context) {
		c.ResponseStatusUnauthorized(errors.New(errorMsg))
	}, http.StatusUnauthorized, "ResponseStatusUnauthorized")

	qhandler(func(c *Context) {
		c.ResponseStatusForbidden(errors.New(errorMsg))
	}, http.StatusForbidden, "ResponseStatusForbidden")

	qhandler(func(c *Context) {
		c.ResponseStatusNotFound(errors.New(errorMsg))
	}, http.StatusNotFound, "ResponseStatusNotFound")

	qhandler(func(c *Context) {
		c.ResponseStatusInternalServerError(errors.New(errorMsg))
	}, http.StatusInternalServerError, "ResponseStatusInternalServerError")

}
