package rsvp

import (
	stdlog "log"
	"net/http"

	"github.com/Teajey/rsvp/internal/log"
)

// WriterHandler is a [Handler] with access to [http.ResponseWriter].
//
// This interface is primarily intended as an adapter for third-party middleware
// that requires [http.ResponseWriter].
//
// WARNING: Care should be taken that [http.ResponseWriter.WriteHeader] is not used here, as response
// status is managed by [Response.Write]. Calling WriteHeader will cause incorrect
// behavior.
type WriterHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) Response
}

type WriterHandlerFunc func(w http.ResponseWriter, r *http.Request) Response

func (f WriterHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) Response {
	return f(w, r)
}

// Adapt wraps [WriterHandler] and returns the standard [http.HandlerFunc].
func Adapt(cfg Config, handler WriterHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := handler.ServeHTTP(w, r)
		if err := response.Write(w, r, cfg); err != nil {
			http.Error(w, "RSVP failed to write a response", http.StatusInternalServerError)
			stdlog.Printf("[ERROR] RSVP failed to write a response: %s", err)
		}
	}
}

// AdaptFunc is a convenience that wraps [Adapt] for [WriterHandlerFunc]
func AdaptFunc(cfg Config, handler func(w http.ResponseWriter, r *http.Request) Response) http.HandlerFunc {
	return Adapt(cfg, WriterHandlerFunc(handler))
}

// toWriterHandler adapts a [Handler] to [WriterHandler]
func toWriterHandler(h Handler) WriterHandler {
	return WriterHandlerFunc(func(wr http.ResponseWriter, r *http.Request) Response {
		w := response{
			std: wr,
		}
		res := h.ServeHTTP(&w, r)
		if res.TemplateName == "" && w.defaultTemplateName != "" {
			log.Dev("Using default template name: %v", w.defaultTemplateName)
			res.TemplateName = w.defaultTemplateName
		}
		return res
	})
}

// AdaptHandler wraps a [Handler] and returns a standard [http.HandlerFunc].
// Convenience for integrating with net/http.
func AdaptHandler(cfg Config, h Handler) http.HandlerFunc {
	return Adapt(cfg, toWriterHandler(h))
}

// AdaptHandlerFunc wraps a [HandlerFunc] and returns a standard [http.HandlerFunc].
// Convenience for integrating with net/http.
func AdaptHandlerFunc(cfg Config, h func(w ResponseWriter, r *http.Request) Response) http.HandlerFunc {
	return Adapt(cfg, toWriterHandler(HandlerFunc(h)))
}
