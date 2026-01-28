package rsvp_test

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
)

func echoHandler(w rsvp.ResponseWriter, r *http.Request) rsvp.Body {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("echoHandler error: %s", err)
		return rsvp.Data(err.Error()).StatusBadRequest()
	}
	return rsvp.Data(string(body))
}

func compressionMiddleware(cfg rsvp.Config, next func(w rsvp.ResponseWriter, r *http.Request) rsvp.Body) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		r.Body = http.MaxBytesReader(w, flate.NewReader(r.Body), 10_000_000) // 10MB limit to protect against zip bombs

		fl, err := flate.NewWriter(&buf, flate.BestCompression)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to init deflate writer: %s", err), http.StatusInternalServerError)
			return
		}

		status, err := rsvp.Write(fl, cfg, w.Header(), r, next)
		if err != nil {
			http.Error(w, fmt.Sprintf("RSVP failed to write: %s", err), http.StatusInternalServerError)
			return
		}
		err = fl.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to close deflate writer: %s", err), http.StatusInternalServerError)
			return
		}

		err = rsvp.WriteResponse(status, w, &buf)
		if err != nil {
			log.Printf("Failed to write response: %s", err)
		}
	}
}

func compressionClient(handler http.Handler, method, target, reqBody string) (status int, body string, err error) {
	rec := httptest.NewRecorder()

	var buf bytes.Buffer
	fl, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		err = fmt.Errorf("initing deflate writer: %s", err)
		return
	}

	_, err = fl.Write([]byte(reqBody))
	if err != nil {
		err = fmt.Errorf("writing deflate: %s", err)
		return
	}
	err = fl.Close()
	if err != nil {
		err = fmt.Errorf("closing deflate writer: %s", err)
		return
	}
	req := httptest.NewRequest(method, target, &buf)

	handler.ServeHTTP(rec, req)
	status = rec.Code

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Failed to copy to http.ResponseWriter: %s", err)
		}
	}()

	bodyBytes, err := io.ReadAll(flate.NewReader(res.Body))
	if err != nil {
		err = fmt.Errorf("reading deflated response: %s", err)
		return
	}
	body = string(bodyBytes)

	return
}

func TestWithCompressionMiddleware(t *testing.T) {
	cfg := rsvp.Config{}
	m := rsvp.NewServeMux()
	m.Std.HandleFunc("POST /{$}", compressionMiddleware(cfg, echoHandler))

	reqBody := "Hello, world!"

	status, respBody, err := compressionClient(m, "POST", "/", reqBody)
	assert.FatalErr(t, "client", err)

	assert.Eq(t, "status code", status, 200)
	assert.Eq(t, "expected body", reqBody, respBody)
}
