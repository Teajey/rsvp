package rsvp

import (
	"fmt"
	"net/http"
)

// An rsvp.Handler with access to the writer. Only intended as an
// adapter for middleware that takes http.ResponseWriter. Care should
// be taken that the body is not written to, as that must be handled by
// func (*Response) Write
type MiddleHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) Response
}

type MiddleHandlerFunc func(w http.ResponseWriter, r *http.Request) Response

func (f MiddleHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) Response {
	return f(w, r)
}

// Can be used to run a MiddleHandler, It will run before rsvp writes to the
// HTTP response body
func MiddlewareBeforeWriteResponse(cfg *Config, next MiddleHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := next.ServeHTTP(w, r)
		err := response.Write(w, r, cfg)
		if err != nil {
			panic(fmt.Sprintf("Failed to write rsvp.Response: %s", err))
		}
	}
}

func middlewareGetArgs(next Handler) MiddleHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) Response {
		headers := w.Header()
		return next.ServeHTTP(headers, r)
	}
}

// Middleware adapter for individual handlers. Is useful for
// adding rsvp to handlers selectively, or plugging in rsvp
// without the provided ServeMux
func Middleware(cfg *Config, next Handler) http.HandlerFunc {
	return MiddlewareBeforeWriteResponse(cfg, middlewareGetArgs(next))
}
