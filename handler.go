package rsvp

import (
	"net/http"
)

// Handler is rsvp's counterpart to [http.Handler].
type Handler interface {
	ServeHTTP(w ResponseWriter, r *http.Request) Body
}

// HandlerFunc is a counterpart to [http.HandlerFunc].
type HandlerFunc func(w ResponseWriter, r *http.Request) Body

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *http.Request) Body {
	return f(w, r)
}
