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
		status, err := Write(rw, cfg, rw.Header(), r, next)
		if err != nil {
			http.Error(rw, "RSVP failed to write a response", http.StatusInternalServerError)
			log.Printf("rsvp failed to write a response: %s", err)
			return
		}
		rw.WriteHeader(status)
	}
}

// AdaptHandler wraps a [Handler] as an [http.Handler] with the given config.
//
// This is the primary entrypoint to using rsvp.
func AdaptHandler(config Config, next Handler) http.Handler {
	return AdaptHandlerFunc(config, next.ServeHTTP)
}
