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

func (t T) handler(s cookieSetter) rsvp.Response {
	cookies, err := http.ParseCookie("hello=world")
	assert.FatalErr(t.t, "cookie parse", err)
	s(cookies[0])
	return rsvp.Blank()
}

func TestMiddlewareBeforeWriteResponse(t *testing.T) {
	cfg := rsvp.Config{}
	server := httptest.NewServer(rsvp.Adapt(cfg, rsvp.WriterHandlerFunc(func(w http.ResponseWriter, r *http.Request) rsvp.Response {
		setter := func(cookie *http.Cookie) {
			http.SetCookie(w, cookie)
		}
		return T{t}.handler(setter)
	})))
	defer server.Close()

	res, err := http.Get(server.URL)
	assert.FatalErr(t, "server response", err)

	assert.Eq(t, "expected cookie header", "hello=world", res.Header.Get("Set-Cookie"))
}
