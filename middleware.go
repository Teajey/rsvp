package rsvp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Teajey/rsvp/internal/dev"
)

// AdaptHandlerFunc provides rsvp as middleware.
//
// NOTE: This function is for advanced lower-level use cases. [ServeMux] should be used by default.
//
// It may be used to add rsvp to discrete handlers.
// It is also useful for wrapping rsvp handlers in middleware that requires access to [http.ResponseWriter].
func AdaptHandlerFunc(cfg Config, next HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := WriteHandler(cfg, rw, r, next)
		if err != nil {
			log.Printf("[RSVP ERROR]: %s", err)
		}
	}
}

// WriteHandler writes the result of handler to rw according to cfg.
//
// NOTE: This function is for advanced lower-level use cases.
func WriteHandler(cfg Config, rw http.ResponseWriter, r *http.Request, handler HandlerFunc) error {
	var buf bytes.Buffer
	status, err := Write(&buf, cfg, rw.Header(), r, handler)
	if err != nil {
		http.Error(rw, "RSVP failed to write a response", http.StatusInternalServerError)
		return fmt.Errorf("writing response: %w", err)
	}
	err = WriteResponse(status, rw, &buf)
	if err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	return nil
}

// WriteResponse calls w.WriteHeader(status) and copies r to w.
//
// NOTE: This function is for advanced lower-level use cases.
//
// This function, alongside [Write], should be used to wrap [Handler] in middleware that requires _write_ access to [http.ResponseWriter]. [AdaptHandler] and [AdaptHandlerFunc] may be used for simpler standard middleware that does not write to [http.ResponseWriter].
//
// See this test for an example: https://github.com/Teajey/rsvp/blob/main/middleware_test.go
func WriteResponse(status int, w http.ResponseWriter, r io.Reader) error {
	dev.Log("Setting status to %d", status)
	w.WriteHeader(status)
	_, err := io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("copying to http.ResponseWriter: %w", err)
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
