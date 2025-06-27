package rsvp

import (
	"net/http"
	"net/url"
)

type Request struct {
	inner *http.Request
}

func wrapStdRequest(r *http.Request) Request {
	return Request{
		r,
	}
}

func (r *Request) ParseForm() (url.Values, error) {
	err := r.inner.ParseForm()
	return r.inner.Form, err
}

func (r *Request) PathValue(name string) string {
	return r.inner.PathValue(name)
}

func (r *Request) URL() *url.URL {
	return r.inner.URL
}
