package binding

import (
	"net/http"
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
)

func Default(method, contentType string) Binding {
	if method == "GET" {
		return Form
	} else {
		switch contentType {
		case MIMEJSON:
			return JSON
		case MIMEXML, MIMEXML2:
			return XML
		default: //case MIMEPOSTForm, MIMEMultipartPOSTForm:
			return Form
		}
	}
}

type Binder struct {
}

func bindWith(r *http.Request, obj interface{}, b Binding) error {
	return b.Bind(r, obj)
}

func (binder *Binder) Bind(r *http.Request, obj interface{}) error {
	header := r.Header.Get("Content-Type")
	bind := Default(r.Method, header)
	return bindWith(r, obj, bind)
}
func (binder *Binder) BindJSON(r *http.Request, obj interface{}) error {
	return bindWith(r, obj, JSON)
}
func (binder *Binder) BindXML(r *http.Request, obj interface{}) error {
	return bindWith(r, obj, XML)
}
func (binder *Binder) BindForm(r *http.Request, obj interface{}) error {
	return bindWith(r, obj, Form)
}

func NewBinder() *Binder {
	return &Binder{}
}
