package rsvp_test

import (
	html "html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
)

func TestDefaultTemplateName(t *testing.T) {
	cfg := rsvp.Config{}
	cfg.HtmlTemplate = html.Must(html.New("tm").Parse(`<div>{{if .}}{{.}}{{else}}Nothin!{{end}}</div>`))

	handler := rsvp.AdaptHandlerFunc(cfg, func(w rsvp.ResponseWriter, r *http.Request) rsvp.Body {
		w.DefaultTemplateName("tm")
		return rsvp.Body{Data: "Hello <input> World!"}
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")

	rec := httptest.NewRecorder()
	handler(rec, req)

	resp := rec.Result()

	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Eq(t, "body contents", "<div>Hello &lt;input&gt; World!</div>", rec.Body.String())
}

func TestResponseHeader(t *testing.T) {
	cfg := rsvp.Config{}

	handler := rsvp.AdaptHandlerFunc(cfg, func(w rsvp.ResponseWriter, r *http.Request) rsvp.Body {
		w.Header().Set("hello", "world")
		return rsvp.Blank()
	})

	req := httptest.NewRequest("GET", "/", nil)

	rec := httptest.NewRecorder()
	handler(rec, req)

	resp := rec.Result()

	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "world", resp.Header.Get("hello"))
	assert.Eq(t, "body contents", "", rec.Body.String())
}
