package rsvp

import (
	stdlog "log"
	"net/http"

	"github.com/Teajey/rsvp/internal/log"
)

// AdaptHandlerFunc provides rsvp as middleware.
//
// It may be used to add rsvp to discrete handlers.
// It is also useful for wrapping rsvp handlers in middleware that requires access to [http.ResponseWriter].
func AdaptHandlerFunc(cfg Config, next HandlerFunc) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		w := response{
			std: wr,
		}
		response := next(&w, r)
		if response.TemplateName == "" && w.defaultTemplateName != "" {
			log.Dev("Using default template name: %v", w.defaultTemplateName)
			response.TemplateName = w.defaultTemplateName
		}
		if err := response.Write(wr, r, cfg); err != nil {
			http.Error(wr, "RSVP failed to write a response", http.StatusInternalServerError)
			stdlog.Printf("[ERROR] RSVP failed to write a response: %s", err)
		}
	}
}

// AdaptHandler provides rsvp as middleware.
//
// It may be used to add rsvp to discrete handlers.
// It is also useful for wrapping rsvp handlers in middleware that requires access to [http.ResponseWriter].
func AdaptHandler(cfg Config, next Handler) http.Handler {
	return AdaptHandlerFunc(cfg, next.ServeHTTP)
}
