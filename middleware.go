package rsvp

import (
	"log"
	"net/http"
)

// NewAdapter returns an rsvp middleware that adapts standard http.Handler to rsvp.Handler, using the provided config.
//
// This is the primary entrypoint to using rsvp.
func NewAdapter(cfg Config) Adapter {
	return Adapter{cfg}
}

type Adapter struct {
	config Config
}

func (a Adapter) AdaptFunc(next func(w ResponseWriter, r *http.Request) Body) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := Write(rw, r, a.config, HandlerFunc(next))
		if err != nil {
			log.Printf("rsvp failed to write a response: %s", err)
			return
		}
	})
}

func (a Adapter) Adapt(next Handler) http.Handler {
	return a.AdaptFunc(next.ServeHTTP)
}
