package rsvp

import (
	"net/http"
)

// ServeMux is a wrapper of [http.ServeMux] that consumes [Handler].
type ServeMux struct {
	// Access to the underlying standard http.ServeMux from net/http
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

// Handler defines the basic signature of RSVP's http handlers. It is the same
// as [http.Handler] except that only [http.ResponseWriter.Header] is accessible,
// and a [Response] must be returned.
type Handler interface {
	ServeHTTP(h http.Header, r *http.Request) Response
}

// HandlerFunc is a counterpart to [http.HandlerFunc]
type HandlerFunc func(h http.Header, r *http.Request) Response

func (f HandlerFunc) ServeHTTP(h http.Header, r *http.Request) Response {
	return f(h, r)
}

// Uses the same pattern syntax as [http.ServeMux]
func (m *ServeMux) Handle(pattern string, handler Handler) {
	m.Std.Handle(pattern, AdaptHandler(m.Config, handler))
}

// Uses the same pattern syntax as [http.ServeMux]
func (m *ServeMux) HandleFunc(pattern string, handler func(http.Header, *http.Request) Response) {
	m.Handle(pattern, HandlerFunc(handler))
}
