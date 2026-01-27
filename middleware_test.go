package rsvp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
)

type cookieSetter func(*http.Cookie)

type T struct {
	t *testing.T
}

func (t T) handler(w rsvp.ResponseWriter, r *http.Request, s cookieSetter) rsvp.Response {
	cookies, err := http.ParseCookie("hello=world")
	assert.FatalErr(t.t, "cookie parse", err)
	s(cookies[0])
	return rsvp.Blank()
}

func middleware(next func(w rsvp.ResponseWriter, r *http.Request, s cookieSetter) rsvp.Response) func(w http.ResponseWriter, r *http.Request) rsvp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) rsvp.HandlerFunc {
		setter := func(cookie *http.Cookie) {
			http.SetCookie(w, cookie)
		}
		return func(w rsvp.ResponseWriter, r *http.Request) rsvp.Response {
			return next(w, r, setter)
		}
	}
}

func TestMiddlewareBeforeWriteResponse(t *testing.T) {
	cfg := rsvp.Config{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := middleware(T{t}.handler)
		b := a(w, r)
		rsvp.AdaptHandlerFunc(cfg, b)(w, r)
	}))
	defer server.Close()

	res, err := http.Get(server.URL)
	assert.FatalErr(t, "server response", err)

	assert.Eq(t, "status code", res.StatusCode, 200)
	assert.Eq(t, "expected cookie header", "hello=world", res.Header.Get("Set-Cookie"))
}
