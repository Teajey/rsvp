//go:build rsvp_msgpack

package rsvp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
)

func TestRequestMsgpackInteger(t *testing.T) {
	res := rsvp.Body{Data: 2}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/vnd.msgpack")
	rec := httptest.NewRecorder()

	err := write(res, rec, req, rsvp.Config{})
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	assert.Eq(t, "Status code", http.StatusOK, resp.StatusCode)
	assert.Eq(t, "Content type", "application/vnd.msgpack", resp.Header.Get("Content-Type"))
	body := rec.Body.Bytes()
	assert.SlicesEq(t, "body contents", []byte{0x02}, body)
}

func TestRequestMsgpackEmptyMapUsingFileExtension(t *testing.T) {
	res := rsvp.Body{Data: map[string]string{}}
	req := httptest.NewRequest("GET", "/resource.msgpack", nil)
	rec := httptest.NewRecorder()

	err := write(res, rec, req, rsvp.Config{})
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	assert.Eq(t, "Status code", http.StatusOK, resp.StatusCode)
	assert.Eq(t, "Content type", "application/vnd.msgpack", resp.Header.Get("Content-Type"))
	body := rec.Body.Bytes()
	assert.SlicesEq(t, "body contents", []byte{0x80}, body)
}
