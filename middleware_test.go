package rsvp_test

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
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
	return rsvp.Ok()
}

func TestMiddlewareBeforeWriteResponse(t *testing.T) {
	cfg := rsvp.DefaultConfig()
	server := httptest.NewServer(rsvp.MiddlewareBeforeWriteResponse(cfg, rsvp.MiddleHandlerFunc(func(w http.ResponseWriter, r *http.Request) rsvp.Response {
		setter := func(cookie *http.Cookie) {
			http.SetCookie(w, cookie)
		}
		return T{t}.handler(setter)
	})))
	defer server.Close()

	res, err := http.Get(server.URL)
	assert.FatalErr(t, "server response", err)

	header, err := httputil.DumpResponse(res, false)
	assert.FatalErr(t, "dump response", err)

	assert.SnapshotText(t, string(header))
}
