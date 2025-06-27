package rsvp

import (
	"log"
	"net/http"
)

type ServeMux struct {
	inner *http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		http.NewServeMux(),
	}
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.inner.ServeHTTP(w, r)
}

type Handler interface {
	ServeHTTP(h http.Header, r *Request) Response
}

type HandlerFunc func(h http.Header, r *Request) Response

func (m *ServeMux) HandleFunc(pattern string, handler HandlerFunc) {
	m.inner.HandleFunc(pattern, func(w http.ResponseWriter, stdReq *http.Request) {
		r := wrapStdRequest(stdReq)

		response := handler(w.Header(), &r)

		err := response.Write(w, stdReq)
		if err != nil {
			log.Printf("Failed to write rsvp.Response to bytes: %s\n", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})
}

func (m *ServeMux) Handle(pattern string, handler Handler) {
	m.HandleFunc(pattern, handler.ServeHTTP)
}
