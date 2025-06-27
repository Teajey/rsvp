package rsvp

import (
	"fmt"
	"html/template"
	"net/http"
)

type ServeMux struct {
	inner        *http.ServeMux
	HtmlTemplate *template.Template
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		http.NewServeMux(),
		nil,
	}
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.inner.ServeHTTP(w, r)
}

type Handler interface {
	ServeHTTP(h http.Header, r *Request) Response
}

type HandlerFunc func(h http.Header, r *Request) Response

func (m *ServeMux) HandleFunc(pattern string, handler HandlerFunc) {
	m.inner.HandleFunc(pattern, func(w http.ResponseWriter, stdReq *http.Request) {
		r := wrapStdRequest(stdReq)

		response := handler(w.Header(), &r)

		err := response.Write(w, stdReq, m.HtmlTemplate)
		if err != nil {
			panic(fmt.Sprintf("Failed to write rsvp.Response: %s", err))
		}
	})
}

func (m *ServeMux) Handle(pattern string, handler Handler) {
	m.HandleFunc(pattern, handler.ServeHTTP)
}
