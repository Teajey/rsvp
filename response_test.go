package rsvp_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
)

func TestStringBody(t *testing.T) {
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	err := res.Write(rec, req, nil, nil)
	if err != nil {
		t.Fatalf("Write failed: %s", err)
	}
	statusCode := rec.Result().StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	s := rec.Body.String()
	if s != body {
		t.Fatalf("s does not have the expected value: %#v", s)
	}
}
