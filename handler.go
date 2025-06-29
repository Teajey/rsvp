package rsvp

import (
	"fmt"
	"net/http"
)

type ServeMux struct {
	inner  *http.ServeMux
	Config *Config
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		inner:  http.NewServeMux(),
		Config: DefaultConfig(),
	}
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.inner.ServeHTTP(w, r)
}

type Handler interface {
	ServeHTTP(h http.Header, r *http.Request) Response
}

type HandlerFunc func(h http.Header, r *http.Request) Response

func (m *ServeMux) HandleFunc(pattern string, handler HandlerFunc) {
	m.inner.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		response := handler(w.Header(), r)

		err := response.Write(w, r, m.Config)
		if err != nil {
			panic(fmt.Sprintf("Failed to write rsvp.Response: %s", err))
		}
	})
}

func (m *ServeMux) Handle(pattern string, handler Handler) {
	m.HandleFunc(pattern, handler.ServeHTTP)
}
