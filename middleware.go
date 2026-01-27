package rsvp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

// AdaptHandlerFunc provides rsvp as middleware.
//
// It may be used to add rsvp to discrete handlers.
// It is also useful for wrapping rsvp handlers in middleware that requires access to [http.ResponseWriter].
func AdaptHandlerFunc(cfg Config, next HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := WriteResponse(cfg, rw, r, next)
		if err != nil {
			log.Printf("%s", err)
		}
	}
}

// Writes the result of handler to rw according to cfg.
func WriteResponse(cfg Config, rw http.ResponseWriter, r *http.Request, handler HandlerFunc) error {
	var buf bytes.Buffer
	status, err := Write(cfg, &buf, rw.Header(), r, handler)
	if err != nil {
		http.Error(rw, "RSVP failed to write a response", http.StatusInternalServerError)
		return fmt.Errorf("[ERROR] RSVP failed to write a response: %w", err)
	}
	rw.WriteHeader(status)
	_, err = io.Copy(rw, &buf)
	if err != nil {
		return fmt.Errorf("[ERROR] RSVP failed to copy to http.ResponseWriter: %w", err)
	}
	return nil
}

// AdaptHandler provides rsvp as middleware.
//
// It may be used to add rsvp to discrete handlers.
// It is also useful for wrapping rsvp handlers in middleware that requires access to [http.ResponseWriter].
func AdaptHandler(cfg Config, next Handler) http.Handler {
	return AdaptHandlerFunc(cfg, next.ServeHTTP)
}
