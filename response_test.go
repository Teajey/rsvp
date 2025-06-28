package rsvp_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Teajey/rsvp"
)

func TestStringBody(t *testing.T) {
	t.Fail()
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	err := res.Write(rec, req, nil)
	if err != nil {
		t.Fatalf("Write failed: %s", err)
	}
	s := rec.Body.String()
	if s != body {
		t.Fatalf("s is not the expected value: %s", s)
	}
}
