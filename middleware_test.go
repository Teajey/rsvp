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

type compressedResponseWriter struct {
	http.ResponseWriter
	flate *flate.Writer
}

func (w compressedResponseWriter) Write(data []byte) (int, error) {
	return w.flate.Write(data)
}

func compressionMiddleware(cfg rsvp.Config, next func(w rsvp.ResponseWriter, r *http.Request) rsvp.Body) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, flate.NewReader(r.Body), 10_000_000) // 10MB limit to protect against zip bombs

		fl, err := flate.NewWriter(w, flate.BestCompression)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to init deflate writer: %s", err), http.StatusInternalServerError)
			return
		}

		rw := compressedResponseWriter{
			ResponseWriter: w,
			flate:          fl,
		}

		err = rsvp.Write(rw, r, cfg, next)
		if err != nil {
			http.Error(w, fmt.Sprintf("RSVP failed to write: %s", err), http.StatusInternalServerError)
			return
		}
		err = fl.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to close deflate writer: %s", err), http.StatusInternalServerError)
			return
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
	m := http.NewServeMux()
	cfg := rsvp.Config{}
	m.HandleFunc("POST /{$}", compressionMiddleware(cfg, echoHandler))

	reqBody := "Hello, world!"

	status, respBody, err := compressionClient(m, "POST", "/", reqBody)
	assert.FatalErr(t, "client", err)

	assert.Eq(t, "status code", status, 200)
	assert.Eq(t, "expected body", reqBody, respBody)
}
