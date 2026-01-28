package rsvp

import (
	"net/http"
)

// ServeMux is a wrapper of [http.ServeMux] that consumes [Handler].
type ServeMux struct {
	// Access to the underlying standard http.ServeMux from net/http.
	//
	// This is useful when you need at least one handler to be http.Handler; particularly when a [Handler] has been wrapped with [http.Handler] middleware.
	Std    *http.ServeMux
	Config Config
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		Std: http.NewServeMux(),
	}
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Std.ServeHTTP(w, r)
}

// Handler defines the basic signature of RSVP's http handlers
type Handler interface {
	ServeHTTP(w ResponseWriter, r *http.Request) Response
}

// HandlerFunc is a counterpart to [http.HandlerFunc]
type HandlerFunc func(w ResponseWriter, r *http.Request) Response

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *http.Request) Response {
	return f(w, r)
}

// Uses the same pattern syntax as [http.ServeMux]
func (m *ServeMux) Handle(pattern string, handler Handler) {
	m.Std.Handle(pattern, AdaptHandler(m.Config, handler))
}

// Uses the same pattern syntax as [http.ServeMux]
func (m *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *http.Request) Response) {
	m.Handle(pattern, HandlerFunc(handler))
}
