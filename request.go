package rsvp

import (
	"log"
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

func (r *Request) ParseForm() url.Values {
	err := r.inner.ParseForm()
	if err != nil {
		log.Printf("Failed to parse form: %s\n", err)
	}
	return r.inner.Form
}

func (r *Request) PathValue(name string) string {
	return r.inner.PathValue(name)
}

func (r *Request) URL() *url.URL {
	return r.inner.URL
}
