package rsvp

import (
	"log"
	"net/http"
)

// AdaptHandlerFunc wraps a [HandlerFunc] as an [http.HandlerFunc] with the given config.
//
// This is the primary entrypoint to using rsvp.
func AdaptHandlerFunc(cfg Config, next HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := Write(rw, r, cfg, next)
		if err != nil {
			log.Printf("rsvp failed to write a response: %s", err)
			return
		}
	}
}

// AdaptHandler wraps a [Handler] as an [http.Handler] with the given config.
//
// This is the primary entrypoint to using rsvp.
func AdaptHandler(config Config, next Handler) http.Handler {
	return AdaptHandlerFunc(config, next.ServeHTTP)
}
